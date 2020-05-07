package txs

import (
	sdkContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	types "github.com/lcnem/proximax-pegzone/x/proximax-bridge"
)

func RelayPeg(
	cdc *codec.Codec,
	rpcURL string,
	chainID string,
	claim *types.MsgPeg,
	moniker string,
) error {
	msg := types.NewMsgPegClaim(sdk.ValAddress(claim.Address), claim.MainchainTxHash, claim.ToAddress, claim.Amount)

	cliCtx := sdkContext.NewCLIContext().
		WithCodec(cdc).
		WithFromAddress(sdk.AccAddress(msg.Address))
	if rpcURL != "" {
		cliCtx = cliCtx.WithNodeURI(rpcURL)
	}
	cliCtx.SkipConfirm = true

	txBldr := authtypes.NewTxBuilderFromCLI(nil).
		WithTxEncoder(utils.GetTxEncoder(cdc)).
		WithChainID(chainID)

	// Check if destination account exists
	accountRetriever := authtypes.NewAccountRetriever(cliCtx)
	err := accountRetriever.EnsureExists(msg.ToAddress)
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
	cdc *codec.Codec,
	rpcURL string,
	chainID string,
	msg *types.MsgUnpegNotCosignedClaim,
	moniker string,
) error {
	cliCtx := sdkContext.NewCLIContext()
	if rpcURL != "" {
		cliCtx = cliCtx.WithNodeURI(rpcURL)
	}
	cliCtx.SkipConfirm = true

	txBldr := authtypes.NewTxBuilderFromCLI(nil).
		WithTxEncoder(utils.GetTxEncoder(cdc)).
		WithChainID(chainID)

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
	cdc *codec.Codec,
	rpcURL string,
	chainID string,
	msg *types.MsgInvitationNotCosignedClaim,
	moniker string,
) error {
	cliCtx := sdkContext.NewCLIContext()
	if rpcURL != "" {
		cliCtx = cliCtx.WithNodeURI(rpcURL)
	}
	cliCtx.SkipConfirm = true

	txBldr := authtypes.NewTxBuilderFromCLI(nil).
		WithTxEncoder(utils.GetTxEncoder(cdc)).
		WithChainID(chainID)

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
