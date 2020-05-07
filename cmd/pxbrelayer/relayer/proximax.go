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
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket"
	tmLog "github.com/tendermint/tendermint/libs/log"
)

func InitProximaXRelayer(
	cdc *codec.Codec,
	cliContext sdkContext.CLIContext,
	logger tmLog.Logger,
	proximaxNode string,
	proximaxPrivateKey string,
	proximaxMultisigAddress string,
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

	if tx.InnerTransactions[0].GetAbstractTransaction().Type == sdk.Transfer || tx.InnerTransactions[0].GetAbstractTransaction().Type == sdk.ModifyMultisig {
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
