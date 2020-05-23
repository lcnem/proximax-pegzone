package relayer

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/lcnem/proximax-pegzone/cmd/pxbrelayer/txs"
	msgTypes "github.com/lcnem/proximax-pegzone/x/proximax-bridge"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket"
	tmLog "github.com/tendermint/tendermint/libs/log"
)

func InitProximaXRelayer(
	inBuf io.Reader,
	cdc *codec.Codec,
	logger tmLog.Logger,

	tendermintNode string,
	chainID string,
	validatorMoniker string,

	proximaxNode string,
	proximaxPrivateKey string,
	proximaxMultisigAddress string,
) error {
	validatorAddress, validatorName, err := LoadValidatorCredentials(validatorMoniker, inBuf)
	if err != nil {
		return err
	}

	cliCtx := LoadTendermintCLIContext(cdc, validatorAddress, validatorName, tendermintNode, chainID)
	txBldr := authtypes.NewTxBuilderFromCLI(nil).
		WithTxEncoder(utils.GetTxEncoder(cdc)).
		WithChainID(chainID)

	conf, err := sdk.NewConfig(context.Background(), []string{proximaxNode})
	if err != nil {
		return err
	}

	client := sdk.NewClient(nil, conf)
	wsClient, err := websocket.NewClient(context.Background(), conf)
	if err != nil {
		return err
	}

	account, err := client.NewAccountFromPrivateKey(proximaxPrivateKey)
	if err != nil {
		return err
	}
	multisigAccount, err := client.NewAccountFromPublicKey(proximaxMultisigAddress)
	if err != nil {
		return err
	}

	err = wsClient.AddPartialAddedHandlers(account.Address, func(tx *sdk.AggregateTransaction) bool {
		partialAddedHandler(client, logger, account, tx)
		return false
	})

	err = wsClient.AddStatusHandlers(account.Address, func(info *sdk.StatusInfo) bool {
		hash := info.Hash.String()

		msg := msgTypes.NewMsgNotCosignedClaim(validatorAddress, hash)
		err := txs.RelayNotCosigned(cliCtx, txBldr, validatorMoniker, msg)
		if err != nil {
			logger.Error("Failed to Relay NotCosigned", "err", err)
		}
		return false
	})

	err = wsClient.AddCosignatureHandlers(multisigAccount.Address, func(info *sdk.SignerInfo) bool {
		txHash := info.ParentHash.String()
		signerPublicKey := info.Signer

		msg := msgTypes.NewMsgNotifyCosigned(validatorAddress, txHash, signerPublicKey)
		err := txs.RelayNotifyCosigned(cliCtx, txBldr, validatorMoniker, msg)
		if err != nil {
			logger.Error("Failed to Relay NotifyCosigned", "err", err)
		}

		return true
	})

	err = wsClient.AddConfirmedAddedHandlers(multisigAccount.Address, func(info sdk.Transaction) bool {
		aggregateTx, ok := info.(*sdk.AggregateTransaction)
		if ok {
			txHash := aggregateTx.TransactionHash.String()

			for _, tx := range aggregateTx.InnerTransactions {
				_, ok := tx.(*sdk.ModifyMultisigAccountTransaction)
				if ok {
					msg := msgTypes.NewMsgConfirmedInvitation(validatorAddress, txHash)
					err := txs.RelayConfirmedInvitation(cliCtx, txBldr, validatorMoniker, msg)
					if err != nil {
						logger.Error("Failed to Relay ConfirmedInvitation", "err", err)
					}
				}
			}
		}

		return true
	})

	go wsClient.Listen()

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
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
