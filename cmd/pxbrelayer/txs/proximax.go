package txs

import (
	"context"
	"time"

	msgTypes "github.com/lcnem/proximax-pegzone/x/proximax-bridge"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

func RelayUnpeg(client *sdk.Client, firstCosignatoryPrivateKey string, msg *msgTypes.MsgUnpeg) error {
	multisigAddress, err := sdk.NewAddressFromRaw(msg.MainchainAddress)
	if err != nil {
		return err
	}
	multisigAccountInfo, err := client.Account.GetAccountInfo(context.Background(), multisigAddress)
	if err != nil {
		return err
	}
	multisigAccount, err := client.NewAccountFromPublicKey(multisigAccountInfo.PublicKey)
	if err != nil {
		return err
	}

	firstCosignatory, err := client.NewAccountFromPrivateKey(firstCosignatoryPrivateKey)
	if err != nil {
		return err
	}

	amount := msg.Amount[0].Amount.BigInt().Uint64()
	transferTx, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour*1),
		sdk.NewAddress(msg.MainchainAddress, client.NetworkType()),
		[]*sdk.Mosaic{sdk.Xpx(amount)},
		sdk.NewPlainMessage(""),
	)
	if err != nil {
		return err
	}
	transferTx.ToAggregate(multisigAccount)

	aggregateBoundedTx, err := client.NewBondedAggregateTransaction(
		sdk.NewDeadline(time.Hour*1),
		[]sdk.Transaction{transferTx},
	)
	if err != nil {
		return err
	}

	signedAggregateBoundedTx, err := firstCosignatory.Sign(aggregateBoundedTx)
	if err != nil {
		return err
	}

	lockFundsTx, err := client.NewLockFundsTransaction(
		sdk.NewDeadline(time.Hour*1),
		sdk.XpxRelative(10),
		sdk.Duration(1000),
		signedAggregateBoundedTx,
	)
	if err != nil {
		return err
	}

	signedLockFundsTx, err := firstCosignatory.Sign(lockFundsTx)
	if err != nil {
		return err
	}

	_, err = client.Transaction.Announce(context.Background(), signedLockFundsTx)
	if err != nil {
		return err
	}

	time.Sleep(30 * time.Second)

	_, _ = client.Transaction.AnnounceAggregateBonded(context.Background(), signedAggregateBoundedTx)
	if err != nil {
		return err
	}
	return nil
}

func RelayInvitation() {

}
