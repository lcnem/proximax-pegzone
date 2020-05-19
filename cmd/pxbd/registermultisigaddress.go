package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil"

	bridge "github.com/lcnem/proximax-pegzone/x/proximax-bridge"
)

// AddGenesisAccountCmd returns add-genesis-account cobra Command.
func RegisterMultisigAddressCmd(
	ctx *server.Context, cdc *codec.Codec, defaultNodeHome, defaultClientHome string,
) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "register-multisig [multisig_account_address]",
		Short: "Register Multisig Account Address to genesis.json",
		Long:  `Register Multisig Account Address to genesis.json.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			multisigAddress := args[0]
			if len(strings.Trim(multisigAddress, "")) == 0 {
				return errors.New(fmt.Sprintf("invalid [multisig_account_address]: %s", multisigAddress))
			}

			genFile := config.GenesisFile()
			appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			bridgeState := bridge.GetGenesisStateFromAppState(cdc, appState)
			bridgeState.MainchainMultisigAddress = multisigAddress

			bridgeStateBz, err := cdc.MarshalJSON(bridgeState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}

			appState[bridge.ModuleName] = bridgeStateBz

			appStateJSON, err := cdc.MarshalJSON(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	cmd.Flags().String(flagClientHome, defaultClientHome, "client's home directory")

	return cmd
}
