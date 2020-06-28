package relayer

import (
	"context"
	"fmt"
	"os"

	sdkContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	cosmosSdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/lcnem/proximax-pegzone/cmd/pxbrelayer/txs"
	msgTypes "github.com/lcnem/proximax-pegzone/x/proximax-bridge"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket"
	tmLog "github.com/tendermint/tendermint/libs/log"
	tmClient "github.com/tendermint/tendermint/rpc/client"
)

type ProximaXSub struct {
	Cdc    *codec.Codec
	CliCtx sdkContext.CLIContext
	TxBldr authtypes.TxBuilder
	Logger tmLog.Logger

	ChainId string

	ValidatorMoniker string
	ValidatorAddress cosmosSdk.ValAddress
	SignerAccount    *sdk.Account
	MultisigAccount  *sdk.PublicAccount

	TendermintClient *tmClient.HTTP
	ProximaXClient   *sdk.Client
	ProximaXWsClient websocket.CatapultClient
}

func NewProximaxSub(cdc *codec.Codec, cliCtx sdkContext.CLIContext, txBldr authtypes.TxBuilder, logger tmLog.Logger, chainID, validatorMoniker string, validatorAddress cosmosSdk.ValAddress, proximaXNode, proximaXPrivateKey, proximaXMultisibPublicKey string) (ProximaXSub, error) {
	conf, err := sdk.NewConfig(context.Background(), []string{proximaXNode})
	if err != nil {
		return ProximaXSub{}, err
	}

	client := sdk.NewClient(nil, conf)
	wsClient, err := websocket.NewClient(context.Background(), conf)
	if err != nil {
		return ProximaXSub{}, err
	}

	account, err := client.NewAccountFromPrivateKey(proximaXPrivateKey)
	if err != nil {
		return ProximaXSub{}, err
	}
	multisigAccount, err := client.NewAccountFromPublicKey(proximaXMultisibPublicKey)
	if err != nil {
		return ProximaXSub{}, err
	}

	return ProximaXSub{
		Cdc:              cdc,
		CliCtx:           cliCtx,
		TxBldr:           txBldr,
		Logger:           logger,
		ChainId:          chainID,
		ValidatorMoniker: validatorMoniker,
		ValidatorAddress: validatorAddress,
		SignerAccount:    account,
		MultisigAccount:  multisigAccount,
		ProximaXClient:   client,
		ProximaXWsClient: wsClient,
	}, nil
}

func (sub *ProximaXSub) Start(exitSignal chan os.Signal) error {
	err := sub.ProximaXWsClient.AddPartialAddedHandlers(sub.SignerAccount.Address, func(tx *sdk.AggregateTransaction) bool {
		partialAddedHandler(sub.ProximaXClient, sub.Logger, sub.SignerAccount, tx)
		return false
	})
	if err != nil {
		return err
	}

	err = sub.ProximaXWsClient.AddStatusHandlers(sub.SignerAccount.Address, func(info *sdk.StatusInfo) bool {
		hash := info.Hash.String()

		msg := msgTypes.NewMsgNotCosignedClaim(sub.ValidatorAddress, hash)
		err := txs.RelayNotCosigned(sub.CliCtx, sub.TxBldr, sub.ValidatorMoniker, msg)
		if err != nil {
			sub.Logger.Error("Failed to Relay NotCosigned", "err", err)
		}
		return false
	})
	if err != nil {
		return err
	}

	err = sub.ProximaXWsClient.AddCosignatureHandlers(sub.MultisigAccount.Address, func(info *sdk.SignerInfo) bool {
		txHash := info.ParentHash.String()
		signerPublicKey := info.Signer

		msg := msgTypes.NewMsgNotifyCosigned(sub.ValidatorAddress, txHash, signerPublicKey)
		err := txs.RelayNotifyCosigned(sub.CliCtx, sub.TxBldr, sub.ValidatorMoniker, msg)
		if err != nil {
			sub.Logger.Error("Failed to Relay NotifyCosigned", "err", err)
		}

		return true
	})
	if err != nil {
		return err
	}

	err = sub.ProximaXWsClient.AddConfirmedAddedHandlers(sub.MultisigAccount.Address, func(info sdk.Transaction) bool {
		aggregateTx, ok := info.(*sdk.AggregateTransaction)
		if ok {
			txHash := aggregateTx.TransactionHash.String()

			for _, tx := range aggregateTx.InnerTransactions {
				_, ok := tx.(*sdk.ModifyMultisigAccountTransaction)
				if ok {
					msg := msgTypes.NewMsgConfirmedInvitation(sub.ValidatorAddress, txHash)
					err := txs.RelayConfirmedInvitation(sub.CliCtx, sub.TxBldr, sub.ValidatorMoniker, msg)
					if err != nil {
						sub.Logger.Error("Failed to Relay ConfirmedInvitation", "err", err)
					}
				}
			}
		}

		return true
	})
	if err != nil {
		return err
	}

	go sub.ProximaXWsClient.Listen()

	<-exitSignal

	return nil
}

func partialAddedHandler(client *sdk.Client, logger tmLog.Logger, account *sdk.Account, tx *sdk.AggregateTransaction) {
	if len(tx.InnerTransactions) != 1 {
		return
	}

	if txType := tx.InnerTransactions[0].GetAbstractTransaction().Type; txType == sdk.Transfer || txType == sdk.ModifyMultisig {
		if tx.Signer.PublicKey == account.PublicAccount.PublicKey {
			return
		}
		for _, cos := range tx.Cosignatures {
			logger.Info(fmt.Sprintf("Singed Cosigner: %s %s", cos.Signer.PublicKey, account.PublicAccount.PublicKey))
			if cos.Signer.PublicKey == account.PublicAccount.PublicKey {
				return
			}
		}

		cosignatureTransaction := sdk.NewCosignatureTransactionFromHash(tx.TransactionHash)
		signedFirstAccountCosignatureTransaction, err := account.SignCosignatureTransaction(cosignatureTransaction)
		if err != nil {
			logger.Error("Failed to subscribe to query", "err", err)
			return
		}
		_, err = client.Transaction.AnnounceAggregateBondedCosignature(context.Background(), signedFirstAccountCosignatureTransaction)
		if err != nil {
			logger.Error("Failed to subscribe to query", "err", err)
			return
		}
		logger.Info(fmt.Sprintf("Signed Transaction: %s", tx.TransactionHash))
	}
}
