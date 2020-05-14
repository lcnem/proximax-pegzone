package relayer

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	sdkContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	cosmosSdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lcnem/proximax-pegzone/cmd/pxbrelayer/txs"
	msgTypes "github.com/lcnem/proximax-pegzone/x/proximax-bridge"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket"
	tmLog "github.com/tendermint/tendermint/libs/log"
)

func InitProximaXRelayer(
	cdc *codec.Codec,
	cli sdkContext.CLIContext,
	logger tmLog.Logger,
	tendermintNode string,
	chainID string,
	proximaxNode string,
	proximaxPrivateKey string,
	proximaxMultisigAddress string,
	validatorAddress cosmosSdk.ValAddress,
	validatorMoniker string,
) error {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*30)
	defer cancel()

	conf, err := sdk.NewConfig(ctx, []string{proximaxNode})
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

	err = wsClient.AddPartialAddedHandlers(account.Address, func(tx *sdk.AggregateTransaction) bool {
		partialAddedHandler(client, logger, account, tx)
		return false
	})

	err = wsClient.AddStatusHandlers(account.Address, func(info *sdk.StatusInfo) bool {
		hash := info.Hash.String()

		msg := msgTypes.NewMsgUnpegNotCosignedClaim(validatorAddress, hash)
		err := txs.RelayUnpegNotCosigned(cdc, cli, tendermintNode, chainID, msg, validatorMoniker)
		if err != nil {
			logger.Error("Failed to Relay UnpegNotCosigned", "err", err)
		}
		return false
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
