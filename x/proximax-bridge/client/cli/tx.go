package cli

import (
	"bufio"
	"errors"
	"fmt"
	"strings"

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
		GetCmdPeg(cdc),
		GetCmdUnpeg(cdc),
		GetCmdRequestInvitation(cdc),
	)...)

	return proximaxbridgeTxCmd
}

func GetCmdPeg(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "peg [key_or_address] [mainchain_tx_hash] [amount]",
		Short: "Peg",
		Args:  cobra.ExactArgs(3), // Does your request require arguments
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			cosmosSender, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			mainchainTxHash := args[1]
			if len(strings.Trim(mainchainTxHash, "")) == 0 {
				return errors.New(fmt.Sprintf("invalid [mainchain_tx_hash]: %s", mainchainTxHash))
			}

			coins, err := sdk.ParseCoins(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgPeg(cosmosSender, mainchainTxHash, coins)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdUnpeg(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "unpeg [key_or_address] [mainchain_address] [amount] [first_cosigner_address]",
		Short: "Unpeg",
		Args:  cobra.ExactArgs(4), // Does your request require arguments
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			cosmosSender, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			mainChainAddress := args[1]
			if len(strings.Trim(mainChainAddress, "")) == 0 {
				return errors.New(fmt.Sprintf("invalid [mainchain_address]: %s", mainChainAddress))
			}

			amount, err := sdk.ParseCoins(args[2])
			if err != nil {
				return err
			}

			firstCosignerAddress, err := sdk.ValAddressFromBech32(args[3])
			if err != nil {
				return err
			}

			msg := types.NewMsgUnpeg(cosmosSender, mainChainAddress, amount, firstCosignerAddress)
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
		Use:   "request-invitation [from_key_or_address] [new_cosigner_public_key] [first_cosigner_address]",
		Short: "Request invitation for multisig cosigner",
		Args:  cobra.ExactArgs(3), // Does your request require arguments
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			newCosignerPublicKey := args[1]
			if len(strings.Trim(newCosignerPublicKey, "")) == 0 {
				return errors.New(fmt.Sprintf("invalid [new_cosigner_public_key]: %s", newCosignerPublicKey))
			}

			firstCosignerAddress, err := sdk.ValAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgRequestInvitation(sdk.ValAddress(cliCtx.FromAddress), newCosignerPublicKey, firstCosignerAddress)
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
