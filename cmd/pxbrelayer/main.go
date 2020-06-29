package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/lcnem/proximax-pegzone/app"
	"github.com/lcnem/proximax-pegzone/cmd/pxbrelayer/relayer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	amino "github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/client/rpc"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/cli"
	tmLog "github.com/tendermint/tendermint/libs/log"
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
	rootCmd.PersistentFlags().String(flags.FlagChainID, "PXB", "Chain ID of tendermint node")
	rootCmd.PersistentFlags().String(FlagRPCURL, "PXB", "RPC URL of tendermint node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		relayerCmd(),
	)

	executor := cli.PrepareMainCmd(rootCmd, "PXB", DefaultCLIHome)
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

func relayerCmd() *cobra.Command {
	relayerCmd := &cobra.Command{
		Use:     "start [tendermint_node] [proximax_node] [validator_from_name] [proximax_private_key] [proximax_multisig_publickey] --chain-id [chain-id]",
		Short:   "Initializes a web socket which streams live events from the ProximaX network and relays them to the Cosmos network",
		Args:    cobra.ExactArgs(5),
		Example: "pxbrelayer start validator http://localhost:7475 http://bctestnet1.brimstone.xpxsirius.io:3000 3A700AE4105431BB0F24440AAA0CD08E1FD0D87D60EB805278F7DB3EA7C7D62D VBK6ZOASJSJOFUOX7XUHHZVCBO4Q11GCF726AKHG --chain-id=testing",
		RunE:    RunRelayerCmd,
	}

	return relayerCmd
}

func RunRelayerCmd(cmd *cobra.Command, args []string) error {
	chainID := viper.GetString(flags.FlagChainID)
	if strings.TrimSpace(chainID) == "" {
		return errors.New("Must specify a 'chain-id'")
	}

	tendermintNode := args[0]
	if len(strings.Trim(tendermintNode, "")) == 0 {
		return errors.New(fmt.Sprintf("invalid [tendermint_node]: %s", tendermintNode))
	}

	proximaXNode := args[1]
	if len(strings.Trim(proximaXNode, "")) == 0 {
		return errors.New(fmt.Sprintf("invalid [proximax_node]: %s", proximaXNode))
	}

	validatorMoniker := args[2]
	if len(strings.Trim(validatorMoniker, "")) == 0 {
		return errors.New(fmt.Sprintf("invalid [validator_from_name]: %s", validatorMoniker))
	}

	cosignerPrivateKey := args[3]
	if len(strings.Trim(cosignerPrivateKey, "")) == 0 {
		return errors.New(fmt.Sprintf("invalid [proximax_private_key]: %s", cosignerPrivateKey))
	}

	multisigPublicKey := args[4]
	if len(strings.Trim(multisigPublicKey, "")) == 0 {
		return errors.New(fmt.Sprintf("invalid [proximax_multisig_address]: %s", multisigPublicKey))
	}

	inBuf := bufio.NewReader(cmd.InOrStdin())
	logger := tmLog.NewTMLogger(tmLog.NewSyncWriter(os.Stdout))

	validatorAddress, validatorName, err := relayer.LoadValidatorCredentials(validatorMoniker, inBuf)
	if err != nil {
		return err
	}

	cliCtx := relayer.LoadTendermintCLIContext(appCodec, validatorAddress, validatorName, tendermintNode, chainID)
	txBldr := authtypes.NewTxBuilderFromCLI(nil).
		WithTxEncoder(utils.GetTxEncoder(appCodec)).
		WithChainID(chainID)

	proximaXSub, err := relayer.NewProximaxSub(appCodec, cliCtx, txBldr, logger, chainID, validatorMoniker, validatorAddress, proximaXNode, cosignerPrivateKey, multisigPublicKey)
	if err != nil {
		return err
	}
	cosmosSub, err := relayer.NewCosmosSub(appCodec, cliCtx, txBldr, logger, tendermintNode, chainID, validatorMoniker, validatorAddress, proximaXNode, cosignerPrivateKey, multisigPublicKey)
	if err != nil {
		return err
	}

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)

	go proximaXSub.Start(exitSignal)
	go cosmosSub.Start(exitSignal)
	<-exitSignal
	return nil
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
