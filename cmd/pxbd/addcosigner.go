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
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"

	bridge "github.com/lcnem/proximax-pegzone/x/proximax-bridge"
)

// AddGenesisAccountCmd returns add-genesis-account cobra Command.
func AddCosignerCmd(
	ctx *server.Context, cdc *codec.Codec, defaultNodeHome, defaultClientHome string,
) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "add-cosigner [address_or_key_name] [cosigner_public_key]",
		Short: "Add a cosigner account to genesis.json",
		Long:  `Add a cosigner account to genesis.json.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			validator, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			mainchainPublicKey := args[1]
			if len(strings.Trim(mainchainPublicKey, "")) == 0 {
				return errors.New(fmt.Sprintf("invalid [cosigner_public_key]: %s", mainchainPublicKey))
			}

			genFile := config.GenesisFile()
			appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			bridgeState := bridge.GetGenesisStateFromAppState(cdc, appState)
			for _, cosigner := range bridgeState.Cosigners {
				if cosigner.MainchainPublicKey == mainchainPublicKey {
					return errors.New(fmt.Sprintf("Cosigner has already been added: %s", mainchainPublicKey))
				}
			}

			cosigner := bridge.Cosigner{ValidatorAddress: validator.String(), MainchainPublicKey: mainchainPublicKey}
			bridgeState.Cosigners = append(bridgeState.Cosigners, cosigner)

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
