package relayer

import (
	"context"
	"time"

	sdkContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	proximax "github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket"
)

func InitProximaXRelayer(
	cdc *codec.Codec,
	cliContext sdkContext.CLIContext,
	proximaxNode string,
	chainID string,
	rpcURL string,
	validatorName string,
	validatorAddress sdk.ValAddress,
	proximaxPrivateKey string,
	proximaxMultisigAddress string,
	test bool,
) error {

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*30)

	conf, err := proximax.NewConfig(ctx, []string{""})
	if err != nil {
		return err
	}

	if test {
		conf.NetworkType = proximax.PublicTest
	} else {
		conf.NetworkType = proximax.Public
	}

	var generationHash proximax.Hash
	if test {
		generationHash = proximax.Hash{}
	} else {
		generationHash = proximax.Hash{}
	}

	account, err := proximax.NewAccountFromPrivateKey(proximaxPrivateKey, conf.NetworkType, &generationHash)
	if err != nil {
		return err
	}

	address, err := proximax.NewAddressFromRaw(proximaxMultisigAddress)
	if err != nil {
		return err
	}

	client, err := websocket.NewClient(ctx, conf)
	if err != nil {
		return err
	}

	err = client.AddPartialAddedHandlers(address, func(tx *proximax.AggregateTransaction) bool {
		partialAddedHandler(*account, tx)
		return false
	})

	go client.Listen()

	for {

	}

	cancel()

	return nil
}

func partialAddedHandler(account proximax.Account, tx *proximax.AggregateTransaction) {
	if len(tx.InnerTransactions) != 1 {
		return
	}
	if tx.InnerTransactions[0].GetAbstractTransaction().Type == proximax.Transfer || tx.InnerTransactions[0].GetAbstractTransaction().Type == proximax.ModifyMultisig {
		cosignatureTransaction := proximax.NewCosignatureTransactionFromHash(tx.AggregateHash)
		signedTransaction, err := account.SignCosignatureTransaction(cosignatureTransaction)

		if err != nil {
			return
		}
	}
}
