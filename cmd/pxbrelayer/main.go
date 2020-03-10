package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/lcnem/proximax-pegzone/app"
	"github.com/lcnem/proximax-pegzone/cmd/pxbrelayer/relayer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	sdkContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	amino "github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/client/rpc"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/cli"
)

var appCodec *amino.Codec

const FlagRPCURL = "rpc-url"

func init() {

	// Read in the configuration file for the sdk
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	appCodec = app.MakeCodec()

	DefaultCLIHome := os.ExpandEnv("$HOME/.pxbcli")

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(flags.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentFlags().String(FlagRPCURL, "", "RPC URL of tendermint node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// Construct Initialization Commands
	initCmd.AddCommand(
		proximaxRelayerCmd(),
		flags.LineBreak,
		cosmosRelayerCmd(),
	)

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		initCmd,
	)

	executor := cli.PrepareMainCmd(rootCmd, "PX", DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		log.Fatal("failed executing CLI command", err)
	}
}

var rootCmd = &cobra.Command{
	Use:          "pxbrelayer",
	Short:        "Relayer service which listens for and relays ProximaX network events",
	SilenceUsage: true,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialization subcommands",
}

func proximaxRelayerCmd() *cobra.Command {
	ethereumRelayerCmd := &cobra.Command{
		Use:     "proximax [proximax_node] [validator_from_name] --chain-id [chain-id]",
		Short:   "Initializes a web socket which streams live events from the ProximaX network and relays them to the Cosmos network",
		Args:    cobra.ExactArgs(3),
		Example: "pxbrelayer init proximax http://localhost:7545 validator --chain-id=testing",
		RunE:    RunProximaxRelayerCmd,
	}

	return ethereumRelayerCmd
}

func cosmosRelayerCmd() *cobra.Command {
	cosmosRelayerCmd := &cobra.Command{
		Use:     "cosmos [tendermint_node] [proximax_node]",
		Short:   "Initializes a web socket which streams live events from the Cosmos network and relays them to the ProximaX network",
		Args:    cobra.ExactArgs(2),
		Example: "pxbrelayer init cosmos tcp://localhost:26657 http://localhost:7545",
		RunE:    RunCosmosRelayerCmd,
	}

	return cosmosRelayerCmd
}

// RunProximaxRelayerCmd executes the initProximaxRelayerCmd with the provided parameters
func RunProximaxRelayerCmd(cmd *cobra.Command, args []string) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())

	// Parse chain's ID
	chainID := viper.GetString(flags.FlagChainID)
	if strings.TrimSpace(chainID) == "" {
		return errors.New("Must specify a 'chain-id'")
	}
	rpcURL := viper.GetString(FlagRPCURL)

	// Get the validator's name and account address using their moniker
	validatorAccAddress, validatorName, err := sdkContext.GetFromFields(inBuf, args[1], false)
	if err != nil {
		return err
	}

	// Convert the validator's account address into type ValAddress
	validatorAddress := sdk.ValAddress(validatorAccAddress)

	// Set up our CLIContext
	cliCtx := sdkContext.NewCLIContext().
		WithCodec(appCodec).
		WithFromAddress(sdk.AccAddress(validatorAddress)).
		WithFromName(validatorName)

	return relayer.InitProximaXRelayer(appCodec, cliCtx, args[0], chainID, rpcURL, validatorName, validatorAddress, "", false)
}

// RunCosmosRelayerCmd executes the initCosmosRelayerCmd with the provided parameters
func RunCosmosRelayerCmd(cmd *cobra.Command, args []string) error {
	return relayer.InitCosmosRelayer(args[0], args[1])
}

func initConfig(cmd *cobra.Command) error {
	return viper.BindPFlag(flags.FlagChainID, cmd.PersistentFlags().Lookup(flags.FlagChainID))
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
