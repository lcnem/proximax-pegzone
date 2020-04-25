package relayer

import (
	"context"
	"os"

	sdkContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	proximax "github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-utils/logger"
	tmKv "github.com/tendermint/tendermint/libs/kv"
	tmLog "github.com/tendermint/tendermint/libs/log"
	tmClient "github.com/tendermint/tendermint/rpc/client"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/lcnem/proximax-pegzone/cmd/pxbrelayer/txs"
)

type CosmosSub struct {
	Cdc                *codec.Codec
	TendermintProvider string
	ProximaXProvider   string
	CliCtx             sdkContext.CLIContext
	TxBldr             authtypes.TxBuilder
	ValidatorMonkier   string
	Logger             tmLog.Logger
	ProximaXClient     *proximax.Client
}

func NewCosmosSub(rpcURL string, cdc *codec.Codec, validatorMonkier, chainID, tendermintProvider, proximaXProvicer string, logger tmLog.Logger) CosmosSub {
	cliCtx := sdkContext.NewCLIContext()
	if rpcURL != "" {
		cliCtx = cliCtx.WithNodeURI(rpcURL)
	}
	cliCtx.SkipConfirm = true

	txBldr := authtypes.NewTxBuilderFromCLI(nil).
		WithTxEncoder(utils.GetTxEncoder(cdc)).
		WithChainID(chainID)

	return CosmosSub{
		Cdc:                cdc,
		TendermintProvider: tendermintProvider,
		ProximaXProvider:   proximaXProvicer,
		CliCtx:             cliCtx,
		TxBldr:             txBldr,
		ValidatorMonkier:   validatorMonkier,
		Logger:             logger,
	}
}

func (sub *CosmosSub) Start() {
	conf, err := proximax.NewConfig(context.Background(), []string{sub.ProximaXProvider})
	if err != nil {
		sub.Logger.Error("Failed to initialize ProximaX client", "err", err)
		os.Exit(1)
	}
	sub.ProximaXClient = proximax.NewClient(nil, conf)

	tendermintClient, err := tmClient.NewHTTP(sub.TendermintProvider, "/websocket")
	if err != nil {
		sub.Logger.Error("Failed to start a client", "err", err)
		os.Exit(1)
	}
	tendermintClient.SetLogger(sub.Logger)
	err = tendermintClient.Start()
	if err != nil {
		sub.Logger.Error("Failed to start a client", "err", err)
		os.Exit(1)
	}
	defer tendermintClient.Stop()

	query := "tm.event = 'Tx'"
	out, err := tendermintClient.Subscribe(context.Background(), "test", query, 1000)
	if err != nil {
		sub.Logger.Error("Failed to subscribe to query", "err", err, "query", query)
		os.Exit(1)
	}

	for result := range out {
		tx, ok := result.Data.(tmTypes.EventDataTx)
		if !ok {
			logger.Error("Type casting failed while extracting event data from new tx")
		}

		logger.Info("New transaction witnessed")

		// Iterate over each event inside of the transaction
		for _, event := range tx.Result.Events {
			attributes := event.GetAttributes()
			switch event.Type {
			case "peg_claim":
				sub.handlePegClaim(attributes)
				break
			case "unpeg_not_cosigned_claim":
				sub.handleUnpegNotCosignedClaim(attributes)
				break
			case "invitation_not_cosigned_claim":
				sub.handleUnpegNotCosignedClaim(attributes)
				break
			default:
				break
			}
		}
	}
}

func (sub *CosmosSub) handlePegClaim(attributes []tmKv.Pair) {
	cosmosMsg, err := txs.PegClaimEventToCosmosMsg(attributes)
	if err != nil {
		sub.Logger.Error("Failed to convert PegClaim event to Cosmos Message", "err", err)
		return
	}

	status, err := sub.ProximaXClient.Transaction.GetTransactionStatus(context.Background(), cosmosMsg.MainchainTxHash)
	if err != nil {
		sub.Logger.Error("Transaction.GetTransaction returned error", "err", err)
		return
	}
	if status.Status != "Success" {
		sub.Logger.Error("Transaction status is not Success", "status", status.Status)
		return
	}
	if status.Group != "confirmed" {
		sub.Logger.Error("Transaction is not confirmed", "group", status.Group)
		return
	}
	txs.RelayPeg(sub.CliCtx, sub.TxBldr, cosmosMsg, sub.ValidatorMonkier)
}

func (sub *CosmosSub) handleUnpegNotCosignedClaim(attributes []tmKv.Pair) {
	msg, err := txs.UnpegNotCosignedClaimEventToCosmosMsg(attributes)
	if err != nil {
		sub.Logger.Error("Failed to convert UnpegNotCosignedClaim event to Cosmos Message", "err", err)
		return
	}
	txs.RelayUnpegNotCosigned(sub.CliCtx, sub.TxBldr, msg, sub.ValidatorMonkier)
}

func (sub *CosmosSub) handleInvitationNotCosignedClaim(attributes []tmKv.Pair) {
	msg, err := txs.InvitationNotCosignedClaimEventToCosmosMsg(attributes)
	if err != nil {
		sub.Logger.Error("Failed to convert InvitationNotCosignedClaim event to Cosmos Message", "err", err)
		return
	}
	txs.RelayInvitationNotCosigned(sub.CliCtx, sub.TxBldr, msg, sub.ValidatorMonkier)
}
