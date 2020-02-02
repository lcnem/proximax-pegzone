package main

import (
	"log"
	"os"

	"github.com/lcnem/proximax-pegzone/app"
	"github.com/lcnem/proximax-pegzone/cmd/pxbrelayer/relayer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/flags"
	amino "github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/client/rpc"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/cli"
)

var appCodec *amino.Codec

const FlagRPCURL = "rpc-url"

// FlagMakeClaims : optional flag for the proximax relayer to automatically make OracleClaims upon every ProphecyClaim
const FlagMakeClaims = "make-claims"

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

	// Add --make-claims to init cmd as optional flag
	initCmd.PersistentFlags().String(FlagMakeClaims, "", "Make oracle claims everytime a prophecy claim is witnessed")

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
		Use:     "proximax [proximax_node] [validator_from_name] --make-claims [make-claims] --chain-id [chain-id]",
		Short:   "Initializes a web socket which streams live events from the ProximaX network and relays them to the Cosmos network",
		Args:    cobra.ExactArgs(3),
		Example: "pxbrelayer init proximax http://localhost:7545 validator --make-claims=false --chain-id=testing",
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
	return relayer.InitProximaXRelayer()
}

// RunCosmosRelayerCmd executes the initCosmosRelayerCmd with the provided parameters
func RunCosmosRelayerCmd(cmd *cobra.Command, args []string) error {
	return relayer.InitCosmosRelayer()
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
