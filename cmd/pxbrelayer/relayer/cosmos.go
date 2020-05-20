package relayer

import (
	"context"
	"fmt"
	"os"

	sdkContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	msgTypes "github.com/lcnem/proximax-pegzone/x/proximax-bridge"
	proximax "github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-utils/logger"
	tmKv "github.com/tendermint/tendermint/libs/kv"
	tmLog "github.com/tendermint/tendermint/libs/log"
	tmClient "github.com/tendermint/tendermint/rpc/client"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/lcnem/proximax-pegzone/cmd/pxbrelayer/txs"
)

type CosmosSub struct {
	Cdc                      *codec.Codec
	RpcUrl                   string
	ChainId                  string
	TendermintProvider       string
	ProximaXProvider         string
	CliCtx                   sdkContext.CLIContext
	TxBldr                   authtypes.TxBuilder
	ValidatorMonkier         string
	ValidatorAddress         sdk.ValAddress
	ProximaxPrivateKey       string
	ProximxMultisigPublicKey string
	Logger                   tmLog.Logger
	ProximaXClient           *proximax.Client
}

func NewCosmosSub(rpcURL string, cdc *codec.Codec, validatorMonkier string, validatorAddress sdk.ValAddress, chainID, tendermintProvider, proximaXProvicer, proximaxPrivateKey, proximxMultisigPublicKey string, logger tmLog.Logger) CosmosSub {
	return CosmosSub{
		Cdc:                      cdc,
		RpcUrl:                   rpcURL,
		ChainId:                  chainID,
		TendermintProvider:       tendermintProvider,
		ProximaXProvider:         proximaXProvicer,
		ValidatorMonkier:         validatorMonkier,
		ValidatorAddress:         validatorAddress,
		ProximaxPrivateKey:       proximaxPrivateKey,
		ProximxMultisigPublicKey: proximxMultisigPublicKey,
		Logger:                   logger,
	}
}

func (sub *CosmosSub) Start(exitSignal chan os.Signal) {
	conf, err := proximax.NewConfig(context.Background(), []string{sub.ProximaXProvider})
	if err != nil {
		sub.Logger.Error("Failed to initialize ProximaX client", "err", err)
		os.Exit(1)
	}
	sub.ProximaXClient = proximax.NewClient(nil, conf)

	tendermintClient, err := tmClient.NewHTTP(sub.TendermintProvider, "/websocket")
	if err != nil {
		sub.Logger.Error("Failed to start a client", "err", err)
		os.Exit(1)
	}
	tendermintClient.SetLogger(sub.Logger)
	err = tendermintClient.Start()
	if err != nil {
		sub.Logger.Error("Failed to start a client", "err", err)
		os.Exit(1)
	}
	defer tendermintClient.Stop()

	query := "tm.event = 'Tx'"
	out, err := tendermintClient.Subscribe(context.Background(), "test", query, 1000)
	if err != nil {
		sub.Logger.Error("Failed to subscribe to query", "err", err, "query", query)
		os.Exit(1)
	}

	for {
		select {
		case result := <-out:
			tx, ok := result.Data.(tmTypes.EventDataTx)
			if !ok {
				logger.Error("Type casting failed while extracting event data from new tx")
			}

			logger.Info("New transaction witnessed")

			// Iterate over each event inside of the transaction
			for _, event := range tx.Result.Events {
				attributes := event.GetAttributes()
				switch event.Type {
				case "peg":
					sub.handlePegEvent(attributes)
					break
				case "unpeg":
					sub.handleUnpegEvent(attributes)
					break
				case "request_invitation":
					sub.handleRequestInvitationEvent(attributes)
					break
				default:
					break
				}
			}
		case <-exitSignal:
			return
		}
	}
}

func (sub *CosmosSub) handlePegEvent(attributes []tmKv.Pair) {
	cosmosMsg, err := txs.PegEventToCosmosMsg(attributes)
	if err != nil {
		sub.Logger.Error("Failed to convert PegClaim event to Cosmos Message", "err", err)
		return
	}

	tx, err := sub.ProximaXClient.Transaction.GetTransaction(context.Background(), cosmosMsg.MainchainTxHash)
	if err != nil {
		sub.Logger.Error("Transaction is not found", "err", err)
		return
	}
	if typ := tx.GetAbstractTransaction().Type; typ != proximax.Transfer {
		sub.Logger.Info("Transaction type is not transfer", "hash", cosmosMsg.MainchainTxHash, "type", typ)
		return
	}

	status, err := sub.ProximaXClient.Transaction.GetTransactionStatus(context.Background(), cosmosMsg.MainchainTxHash)
	if err != nil {
		sub.Logger.Error("Transaction.GetTransaction returned error", "err", err)
		return
	}
	if status.Status != "Success" {
		sub.Logger.Error("Transaction status is not Success", "status", status.Status)
		return
	}
	if status.Group != "confirmed" {
		sub.Logger.Error("Transaction is not confirmed", "group", status.Group)
		return
	}
	err = txs.RelayPeg(sub.Cdc, sub.RpcUrl, sub.ChainId, cosmosMsg, sub.ValidatorMonkier, sub.ValidatorAddress)
	if err != nil {
		sub.Logger.Error(fmt.Sprintf("Faild while broadcast transaction: %+v", err))
	}
}

func (sub *CosmosSub) handleUnpegEvent(attributes []tmKv.Pair) {
	msg, err := txs.UnpegEventToCosmosMsg(attributes)
	if err != nil {
		sub.Logger.Error("Failed to convert Unpeg event to Cosmos Message", "err", err)
		return
	}
	if msg.FirstCosignerAddress.String() != sub.ValidatorAddress.String() {
		return
	}
	txHash, err := txs.RelayUnpeg(sub.ProximaXClient, sub.ProximaxPrivateKey, sub.ProximxMultisigPublicKey, msg)
	if err != nil {
		sub.Logger.Error("Failed to Relay Transaction to ProximaX", "err", err)
		return
	}

	firstCosignatory, err := sub.ProximaXClient.NewAccountFromPrivateKey(sub.ProximaxPrivateKey)
	publicKey := firstCosignatory.PublicAccount.PublicKey

	recordMsg := msgTypes.NewMsgRecordUnpeg(msg.Address, msg.FromAddress, txHash, msg.Amount, publicKey)
	err = txs.RelayRecordUnpeg(sub.Cdc, sub.RpcUrl, sub.ChainId, &recordMsg, sub.ValidatorMonkier, msg.FirstCosignerAddress)
	if err != nil {
		sub.Logger.Error(fmt.Sprintf("Faild while broadcast transaction: %+v", err))
	}
}

func (sub *CosmosSub) handleRequestInvitationEvent(attributes []tmKv.Pair) {
	msg, err := txs.RequestInvitationEventToCosmosMsg(attributes)
	if err != nil {
		sub.Logger.Error("Failed to convert RequestInvitation event to Cosmos Message", "err", err)
		return
	}
	if msg.FirstCosignerAddress.String() != sub.ValidatorAddress.String() {
		return
	}
	txs.RelayInvitation(sub.ProximaXClient, sub.ProximaxPrivateKey, msg)
}
