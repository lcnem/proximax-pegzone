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
	chainID string,
	provider string,
	makeClaims bool,
	validatorName string,
	passphrase string,
	validatorAddress sdk.ValAddress,
	cliContext sdkContext.CLIContext,
	rpcURL string,
	custodyAddress string,
	test bool,
	account proximax.Account,
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

	client, err := websocket.NewClient(ctx, conf)
	if err != nil {
		return err
	}
	go client.Listen()

	address, err := proximax.NewAddressFromRaw(custodyAddress)
	if err != nil {
		return err
	}

	err = client.AddPartialAddedHandlers(account.Address, func(tx *proximax.AggregateTransaction) bool {
		partialAddedHandler(account, tx)
		return false
	})

	for {

	}

	cancel()
}

func partialAddedHandler(account proximax.Account, tx *proximax.AggregateTransaction) {
	if len(tx.InnerTransactions) != 1 {
		return
	}
	if tx.InnerTransactions[0].GetAbstractTransaction().Type == proximax.Transfer {
		handleTransferTransaction(account, proximax.TransferTransaction(tx.InnerTransactions[0]))
	}
	if tx.InnerTransactions[0].GetAbstractTransaction().Type == proximax.ModifyMultisig {
		handleModifyMultisigTransaction(account, proximax.TransferTransaction(tx.InnerTransactions[0]))
	}
}

func handleTransferTransaction(account proximax.Account, tx proximax.TransferTransaction) {

}

func handleModifyMultisigTransaction(account proximax.Account, tx proximax.TransferTransaction) {

}
