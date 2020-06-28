package relayer

import (
	"context"
	"fmt"
	"os"

	sdkContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/lcnem/proximax-pegzone/cmd/pxbrelayer/txs"
	msgTypes "github.com/lcnem/proximax-pegzone/x/proximax-bridge"
	"github.com/lcnem/proximax-pegzone/x/proximax-bridge/types"
	proximax "github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-utils/logger"
	tmKv "github.com/tendermint/tendermint/libs/kv"
	tmLog "github.com/tendermint/tendermint/libs/log"
	tmClient "github.com/tendermint/tendermint/rpc/client"
	tmTypes "github.com/tendermint/tendermint/types"
)

type CosmosSub struct {
	Cdc    *codec.Codec
	Logger tmLog.Logger
	CliCtx sdkContext.CLIContext
	TxBldr authtypes.TxBuilder

	ChainId string

	ValidatorMoniker         string
	ValidatorAddress         sdk.ValAddress
	ProximaxPrivateKey       string
	ProximxMultisigPublicKey string

	TendermintClient *tmClient.HTTP
	ProximaXClient   *proximax.Client
}

func NewCosmosSub(cdc *codec.Codec, cliCtx sdkContext.CLIContext, txBldr authtypes.TxBuilder, logger tmLog.Logger, tendermintNode, chainID, validatorMoniker string, validatorAddress sdk.ValAddress, proximaXNode, proximaXPrivateKey, proximaXMultisibPublicKey string) (CosmosSub, error) {
	conf, err := proximax.NewConfig(context.Background(), []string{proximaXNode})
	if err != nil {
		return CosmosSub{}, err
	}

	tendermintClient, err := tmClient.NewHTTP(tendermintNode, "/websocket")
	if err != nil {
		return CosmosSub{}, err
	}
	tendermintClient.SetLogger(logger)

	return CosmosSub{
		Cdc:                      cdc,
		Logger:                   logger,
		CliCtx:                   cliCtx,
		TxBldr:                   txBldr,
		ChainId:                  chainID,
		ValidatorMoniker:         validatorMoniker,
		ValidatorAddress:         validatorAddress,
		ProximaxPrivateKey:       proximaXPrivateKey,
		ProximxMultisigPublicKey: proximaXMultisibPublicKey,
		TendermintClient:         tendermintClient,
		ProximaXClient:           proximax.NewClient(nil, conf),
	}, nil
}

func (sub *CosmosSub) Start(exitSignal chan os.Signal) {

	err := sub.TendermintClient.Start()
	if err != nil {
		sub.Logger.Error("Failed to start a client", "err", err)
		os.Exit(1)
	}
	defer sub.TendermintClient.Stop()

	query := "tm.event = 'Tx'"
	out, err := sub.TendermintClient.Subscribe(context.Background(), "test", query, 1000)
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
	cosmosMsg, consumed, err := txs.PegEventToCosmosMsg(attributes)
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

	transferTx, ok := tx.(*proximax.TransferTransaction)
	if !ok {
		sub.Logger.Error(fmt.Sprintf("Faild to get transfer transaction: %+v", transferTx))
		return
	}

	unit := int64(proximax.XpxRelative(1).Amount)
	var amount int64 = 0
	for _, mosaic := range transferTx.Mosaics {
		amount += int64(mosaic.Amount) / unit
	}

	var request int64 = 0
	for _, coin := range cosmosMsg.Amount {
		request += coin.Amount.Int64()
	}

	if (consumed + request) > amount {
		sub.Logger.Error(fmt.Sprintf("Request amount exceeds remainning %s, request=%d, remainning=%d", cosmosMsg.MainchainTxHash, request, amount-consumed))
		return
	}

	remaiining := amount - consumed - request
	msg := types.NewMsgPegClaim(cosmosMsg.Address, cosmosMsg.MainchainTxHash, cosmosMsg.Amount, remaiining, sub.ValidatorAddress)
	err = txs.RelayPeg(sub.CliCtx, sub.TxBldr, sub.ValidatorMoniker, msg)
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
	recordMsg := msgTypes.NewMsgRecordUnpeg(msg.Address, txHash, msg.Amount, firstCosignatory.PublicKey, sub.ValidatorAddress)
	err = txs.RelayRecordUnpeg(sub.CliCtx, sub.TxBldr, sub.ValidatorMoniker, recordMsg)
	if err != nil {
		sub.Logger.Error(fmt.Sprintf("Faild while broadcast transaction: %+v", err))
	}
}

func (sub *CosmosSub) handleRequestInvitationEvent(attributes []tmKv.Pair) {
	msg, multisigAddress, err := txs.RequestInvitationEventToCosmosMsg(attributes)
	if err != nil {
		sub.Logger.Error("Failed to convert RequestInvitation event to Cosmos Message", "err", err)
		return
	}
	if msg.FirstCosignerAddress.String() != sub.ValidatorAddress.String() {
		return
	}
	txHash, err := txs.RelayInvitation(sub.ProximaXClient, sub.ProximaxPrivateKey, msg, multisigAddress)
	if err != nil {
		sub.Logger.Error("Failed to broadcase ProximaX transaction to add new cosigner", "err", err)
		return
	}

	account, err := sub.ProximaXClient.NewAccountFromPrivateKey(sub.ProximaxPrivateKey)
	if err != nil {
		sub.Logger.Error("Failed to Get Account", "err", err)
		return
	}
	pubKey := account.PublicAccount.PublicKey

	pendingMsg := msgTypes.NewMsgPendingRequestInvitation(msg.Address, msg.NewCosignerPublicKey, msg.FirstCosignerAddress, pubKey, txHash)
	err = txs.RelayPendingRequestInvitation(sub.CliCtx, sub.TxBldr, sub.ValidatorMoniker, pendingMsg)
	if err != nil {
		sub.Logger.Error("Failed to broadcase Cosmos transaction to notify pending request", "err", err)
		return
	}
}
