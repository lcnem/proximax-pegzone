package txs

import (
	"context"
	"encoding/json"
	"math"
	"time"

	msgTypes "github.com/lcnem/proximax-pegzone/x/proximax-bridge"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

func getAccountByAddress(client *sdk.Client, address string) (*sdk.PublicAccount, error) {
	multisigAddress, err := sdk.NewAddressFromRaw(address)
	if err != nil {
		return nil, err
	}
	multisigAccountInfo, err := client.Account.GetAccountInfo(context.Background(), multisigAddress)
	if err != nil {
		return nil, err
	}
	return client.NewAccountFromPublicKey(multisigAccountInfo.PublicKey)
}

func getApprovalDelta(originalNum int32, addedNumber int32) int8 {
	originalMajority := originalNum / 2
	newMajority := (originalNum + addedNumber) / 2
	delta := newMajority - originalMajority
	if delta <= math.MaxInt8 {
		return int8(delta)
	} else {
		return math.MaxInt8
	}
}

func RelayUnpeg(client *sdk.Client, firstCosignatoryPrivateKey, multisigPublicKey string, msg *msgTypes.MsgUnpeg) (string, error) {
	multisigAccount, err := sdk.NewAccountFromPublicKey(multisigPublicKey, client.NetworkType())
	firstCosignatory, err := client.NewAccountFromPrivateKey(firstCosignatoryPrivateKey)
	if err != nil {
		return "", err
	}

	txMsg, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	amount := msg.Amount[0].Amount.BigInt().Uint64()
	transferTx, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour*1),
		sdk.NewAddress(msg.MainchainAddress, client.NetworkType()),
		[]*sdk.Mosaic{sdk.XpxRelative(amount)},
		sdk.NewPlainMessage(string(txMsg)),
	)
	if err != nil {
		return "", err
	}

	transferTx.ToAggregate(multisigAccount)

	aggregateBoundedTx, err := client.NewBondedAggregateTransaction(
		sdk.NewDeadline(time.Hour*1),
		[]sdk.Transaction{transferTx},
	)
	if err != nil {
		return "", err
	}

	signedAggregateBoundedTx, err := firstCosignatory.Sign(aggregateBoundedTx)
	if err != nil {
		return "", err
	}

	lockFundsTx, err := client.NewLockFundsTransaction(
		sdk.NewDeadline(time.Hour*1),
		sdk.XpxRelative(10),
		sdk.Duration(1000),
		signedAggregateBoundedTx,
	)
	if err != nil {
		return "", err
	}

	signedLockFundsTx, err := firstCosignatory.Sign(lockFundsTx)
	if err != nil {
		return "", err
	}

	_, err = client.Transaction.Announce(context.Background(), signedLockFundsTx)
	if err != nil {
		return "", err
	}

	time.Sleep(30 * time.Second)

	_, _ = client.Transaction.AnnounceAggregateBonded(context.Background(), signedAggregateBoundedTx)
	if err != nil {
		return "", err
	}

	return signedAggregateBoundedTx.Hash.String(), nil
}

func RelayInvitation(client *sdk.Client, firstCosignatoryPrivateKey string, msg *msgTypes.MsgRequestInvitation, multisigAccountAddress string) (string, error) {
	multisigAccount, err := getAccountByAddress(client, multisigAccountAddress)
	if err != nil {
		return "", err
	}
	firstCosignatory, err := client.NewAccountFromPrivateKey(firstCosignatoryPrivateKey)
	if err != nil {
		return "", err
	}
	newCosignerAccount, err := client.NewAccountFromPublicKey(msg.NewCosignerPublicKey)
	if err != nil {
		return "", err
	}

	multisigAccountInfo, err := client.Account.GetMultisigAccountInfo(context.Background(), multisigAccount.Address)
	if err != nil {
		return "", err
	}

	// Update majority
	minApprovalDelta := getApprovalDelta(multisigAccountInfo.MinApproval, 1)
	minRemovalDelta := getApprovalDelta(multisigAccountInfo.MinRemoval, 1)

	modifyMultisigTx, err := client.NewModifyMultisigAccountTransaction(
		sdk.NewDeadline(time.Hour*1),
		minApprovalDelta,
		minRemovalDelta,
		[]*sdk.MultisigCosignatoryModification{{sdk.Add, newCosignerAccount}},
	)
	if err != nil {
		return "", err
	}
	modifyMultisigTx.ToAggregate(multisigAccount)

	aggregateBoundedTx, err := client.NewBondedAggregateTransaction(
		sdk.NewDeadline(time.Hour*1),
		[]sdk.Transaction{modifyMultisigTx},
	)
	if err != nil {
		return "", err
	}

	signedAggregateBoundedTx, err := firstCosignatory.Sign(aggregateBoundedTx)
	if err != nil {
		return "", err
	}

	lockFundsTx, err := client.NewLockFundsTransaction(
		sdk.NewDeadline(time.Hour*1),
		sdk.XpxRelative(10),
		sdk.Duration(100),
		signedAggregateBoundedTx,
	)
	if err != nil {
		return "", err
	}

	signedLockFundsTx, err := firstCosignatory.Sign(lockFundsTx)
	if err != nil {
		return "", err
	}

	_, err = client.Transaction.Announce(context.Background(), signedLockFundsTx)
	if err != nil {
		return "", err
	}

	time.Sleep(30 * time.Second)

	_, _ = client.Transaction.AnnounceAggregateBonded(context.Background(), signedAggregateBoundedTx)
	if err != nil {
		return "", err
	}

	return signedAggregateBoundedTx.Hash.String(), nil
}
