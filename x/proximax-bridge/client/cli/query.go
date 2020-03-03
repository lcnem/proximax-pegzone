package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/lcnem/proximax-pegzone/x/proximax-bridge/internal/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group proximax-bridge queries under a subcommand
	proximaxbridgeQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	proximaxbridgeQueryCmd.AddCommand(
		flags.GetCommands(
			// TODO: Add query Cmds
			GetCmdQueryMainchainMultisigAddress(queryRoute, cdc),
		)...,
	)

	return proximaxbridgeQueryCmd

}

// TODO: Add Query Commands
func GetCmdQueryMainchainMultisigAddress(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "mainchain-multisig-address",
		Short: "Get the mainchain multisig address for collateral",
		Args:  cobra.ExactArgs(0), // Does your request require arguments
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryMainchainMultisigAddress), nil)
			if err != nil {
				fmt.Printf(err.Error())
				return nil
			}

			var out map[string]string
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdQueryCosigners(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "mainchain-multisig-address",
		Short: "Get the mainchain cosigners of multisig address",
		Args:  cobra.ExactArgs(0), // Does your request require arguments
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCosigners), nil)
			if err != nil {
				fmt.Printf(err.Error())
				return nil
			}

			var out map[string]string
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
