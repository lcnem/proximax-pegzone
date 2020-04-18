package txs

import (
	sdkContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	types "github.com/lcnem/proximax-pegzone/x/proximax-bridge"
	amino "github.com/tendermint/go-amino"
)

func RelayPeg(
	chainID string,
	cdc *amino.Codec,
	moniker string,
	msg *types.MsgPegClaim,
	rpcURL string,
) error {

	cliCtx := sdkContext.NewCLIContext()
	if rpcURL != "" {
		cliCtx = cliCtx.WithNodeURI(rpcURL)
	}

	cliCtx.SkipConfirm = true

	txBldr := authtypes.NewTxBuilderFromCLI(nil).
		WithTxEncoder(utils.GetTxEncoder(cdc)).
		WithChainID(chainID)

	accountRetriever := authtypes.NewAccountRetriever(cliCtx)

	err := accountRetriever.EnsureExists((sdk.AccAddress(msg.ToAddress)))
	if err != nil {
		return err
	}

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

func RelayUnpegNotCosigned() {

}

func RelayInvitationNotCosigned() {

}
