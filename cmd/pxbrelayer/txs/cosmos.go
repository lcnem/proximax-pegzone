package txs

import (
	sdkContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	types "github.com/lcnem/proximax-pegzone/x/proximax-bridge"
)

func RelayMsg(
	cliCtx sdkContext.CLIContext,
	txBldr authtypes.TxBuilder,
	validatorMoniker string,
	msg sdk.Msg,
) error {
	// Validate message
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}

	// Prepare tx
	txBldr, err = utils.PrepareTxBuilder(txBldr, cliCtx)
	if err != nil {
		return err
	}

	// Build and sign the transaction
	txBytes, err := txBldr.BuildAndSign(validatorMoniker, keys.DefaultKeyPass, []sdk.Msg{msg})
	if err != nil {
		return err
	}

	// Broadcast to a Tendermint node
	res, err := cliCtx.BroadcastTxSync(txBytes)
	if err != nil {
		return err
	}

	if err = cliCtx.PrintOutput(res); err != nil {
		return err
	}
	return nil
}

func RelayPeg(
	cliCtx sdkContext.CLIContext,
	txBldr authtypes.TxBuilder,
	validatorMoniker string,
	msg types.MsgPegClaim,
) error {
	// Check if destination account exists
	accountRetriever := authtypes.NewAccountRetriever(cliCtx)
	err := accountRetriever.EnsureExists(msg.Address)
	if err != nil {
		return err
	}

	return RelayMsg(cliCtx, txBldr, validatorMoniker, msg)
}

func RelayRecordUnpeg(
	cliCtx sdkContext.CLIContext,
	txBldr authtypes.TxBuilder,
	validatorMoniker string,
	msg types.MsgRecordUnpeg,
) error {
	return RelayMsg(cliCtx, txBldr, validatorMoniker, msg)
}

func RelayNotCosigned(
	cliCtx sdkContext.CLIContext,
	txBldr authtypes.TxBuilder,
	validatorMoniker string,
	msg types.MsgNotCosignedClaim,
) error {
	return RelayMsg(cliCtx, txBldr, validatorMoniker, msg)
}

func RelayNotifyCosigned(
	cliCtx sdkContext.CLIContext,
	txBldr authtypes.TxBuilder,
	validatorMoniker string,
	msg types.MsgNotifyCosigned,
) error {
	return RelayMsg(cliCtx, txBldr, validatorMoniker, msg)
}

func RelayPendingRequestInvitation(
	cliCtx sdkContext.CLIContext,
	txBldr authtypes.TxBuilder,
	validatorMoniker string,
	msg types.MsgPendingRequestInvitation,
) error {
	// Check if destination account exists
	accountRetriever := authtypes.NewAccountRetriever(cliCtx)
	err := accountRetriever.EnsureExists(sdk.AccAddress(msg.Address))
	if err != nil {
		return err
	}

	return RelayMsg(cliCtx, txBldr, validatorMoniker, msg)
}

func RelayConfirmedInvitation(
	cliCtx sdkContext.CLIContext,
	txBldr authtypes.TxBuilder,
	validatorMoniker string,
	msg types.MsgConfirmedInvitation,
) error {
	return RelayMsg(cliCtx, txBldr, validatorMoniker, msg)
}
