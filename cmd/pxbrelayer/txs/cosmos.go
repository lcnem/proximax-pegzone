package txs

import (
	sdkContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	types "github.com/lcnem/proximax-pegzone/x/proximax-bridge"
)

func RelayPeg(
	cliCtx sdkContext.CLIContext,
	txBldr authtypes.TxBuilder,
	msg *types.MsgPegClaim,
	moniker string,
) error {

	// Check if destination account exists
	accountRetriever := authtypes.NewAccountRetriever(cliCtx)
	err := accountRetriever.EnsureExists((sdk.AccAddress(msg.ToAddress)))
	if err != nil {
		return err
	}

	// Validate message
	err = msg.ValidateBasic()
	if err != nil {
		return err
	}

	// Prepare tx
	txBldr, err = utils.PrepareTxBuilder(txBldr, cliCtx)
	if err != nil {
		return err
	}

	// Build and sign the transaction
	txBytes, err := txBldr.BuildAndSign(moniker, keys.DefaultKeyPass, []sdk.Msg{msg})
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

func RelayUnpegNotCosigned(
	cliCtx sdkContext.CLIContext,
	txBldr authtypes.TxBuilder,
	msg *types.MsgUnpegNotCosignedClaim,
	moniker string,
) error {

	err := msg.ValidateBasic()
	if err != nil {
		return err
	}

	txBldr, err = utils.PrepareTxBuilder(txBldr, cliCtx)
	if err != nil {
		return err
	}

	txBytes, err := txBldr.BuildAndSign(moniker, keys.DefaultKeyPass, []sdk.Msg{msg})
	if err != nil {
		return err
	}

	res, err := cliCtx.BroadcastTxSync(txBytes)
	if err != nil {
		return err
	}

	if err = cliCtx.PrintOutput(res); err != nil {
		return err
	}
	return nil
}

func RelayInvitationNotCosigned(
	cliCtx sdkContext.CLIContext,
	txBldr authtypes.TxBuilder,
	msg *types.MsgInvitationNotCosignedClaim,
	moniker string,
) error {

	err := msg.ValidateBasic()
	if err != nil {
		return err
	}

	txBldr, err = utils.PrepareTxBuilder(txBldr, cliCtx)
	if err != nil {
		return err
	}

	txBytes, err := txBldr.BuildAndSign(moniker, keys.DefaultKeyPass, []sdk.Msg{msg})
	if err != nil {
		return err
	}

	res, err := cliCtx.BroadcastTxSync(txBytes)
	if err != nil {
		return err
	}

	if err = cliCtx.PrintOutput(res); err != nil {
		return err
	}
	return nil
}
