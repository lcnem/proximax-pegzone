package cli

import (
	"bufio"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/lcnem/proximax-pegzone/x/proximax-bridge/internal/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	proximaxbridgeTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%S transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	proximaxbridgeTxCmd.AddCommand(flags.PostCommands(
		// TODO: Add tx based commands
		// GetCmd<Action>(cdc)
		GetCmdUnpeg(cdc),
		GetCmdRequestInvitation(cdc),
	)...)

	return proximaxbridgeTxCmd
}

func GetCmdUnpeg(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "unpeg [key_or_address] [amount] [mainchain_address] [first_cosigner_address]",
		Short: "Unpeg",
		Args:  cobra.ExactArgs(4), // Does your request require arguments
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			coins, err := sdk.ParseCoins(args[2])
			if err != nil {
				return err
			}

			firstCosignerAddress, err := sdk.ValAddressFromBech32(args[3])
			if err != nil {
				return err
			}

			msg := types.NewMsgUnpeg(cliCtx.GetFromAddress(), args[1], coins, firstCosignerAddress)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdRequestInvitation(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "request-invitation [from_key_or_address] [mainchain_address] [first_cosigner_address]",
		Short: "Request invitation for multisig cosigner",
		Args:  cobra.ExactArgs(3), // Does your request require arguments
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			firstCosignerAddress, err := sdk.ValAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgRequestInvitation(sdk.ValAddress(cliCtx.FromAddress), args[1], firstCosignerAddress)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// Example:
//
// GetCmd<Action> is the CLI command for doing <Action>
// func GetCmd<Action>(cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "/* Describe your action cmd */",
// 		Short: "/* Provide a short description on the cmd */",
// 		Args:  cobra.ExactArgs(2), // Does your request require arguments
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			cliCtx := context.NewCLIContext().WithCodec(cdc)
// 			inBuf := bufio.NewReader(cmd.InOrStdin())
// 			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

// 			msg := types.NewMsg<Action>(/* Action params */)
// 			err = msg.ValidateBasic()
// 			if err != nil {
// 				return err
// 			}

// 			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
// 		},
// 	}
// }
