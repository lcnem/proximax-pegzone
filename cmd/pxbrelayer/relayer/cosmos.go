package relayer

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	proximax "github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	amino "github.com/tendermint/go-amino"
	tmKv "github.com/tendermint/tendermint/libs/kv"
	tmLog "github.com/tendermint/tendermint/libs/log"
	tmClient "github.com/tendermint/tendermint/rpc/client"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/lcnem/proximax-pegzone/cmd/pxbrelayer/txs"
	"github.com/lcnem/proximax-pegzone/cmd/pxbrelayer/types"
)

var proximaxClient *proximax.Client
var tendermintClient *tmClient.HTTP
var logger tmLog.Logger

func InitCosmosRelayer(
	tendermintNode string,
	proximaxNode string,
	cdc *amino.Codec,
	moniker string,
	chainId string,
	rpcUrl string,
) error {
	logger = tmLog.NewTMLogger(tmLog.NewSyncWriter(os.Stdout))

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*30)
	conf, err := proximax.NewConfig(ctx, []string{""})
	if err != nil {
		logger.Error("Failed to initialize ProximaX client", "err", err)
		os.Exit(1)
	}
	proximaxClient = proximax.NewClient(nil, conf)

	tendermintClient, err := tmClient.NewHTTP(tendermintNode, "/websocket")
	tendermintClient.SetLogger(logger)

	err = tendermintClient.Start()
	if err != nil {
		logger.Error("Failed to start a client", "err", err)
		os.Exit(1)
	}

	fmt.Printf("Started Cosmos Relayer %+v\n", tendermintClient)

	defer tendermintClient.Stop()

	// Subscribe to all tendermint transactions
	query := "tm.event = 'Tx'"

	out, err := tendermintClient.Subscribe(context.Background(), "test", query, 1000)
	if err != nil {
		logger.Error("Failed to subscribe to query", "err", err, "query", query)
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case result := <-out:
			tx, ok := result.Data.(tmTypes.EventDataTx)
			if !ok {
				logger.Error("Type casting failed while extracting event data from new tx")
			}

			logger.Info("New transaction witnessed")

			// Iterate over each event inside of the transaction
			for _, event := range tx.Result.Events {
				eventType := getEventType(event.Type)
				switch eventType {
				case types.MsgPegClaim:
					handlePegClaim(ctx, event.GetAttributes(), cdc, moniker, chainId, rpcUrl)
					break
				case types.MsgUnpegNotCosignedClaim:
					handleUnpegNotCosignedClaim()
					break
				case types.MsgInvitationNotCosignedClaim:
					handleUnpegNotCosignedClaim()
					break
				default:
					break
				}
			}
		case <-quit:
			cancel()
			os.Exit(0)
		}
	}
}

func getEventType(eventType string) types.Event {
	switch eventType {
	case types.MsgPegClaim.String():
		return types.MsgPegClaim
	case types.MsgUnpegNotCosignedClaim.String():
		return types.MsgUnpegNotCosignedClaim
	case types.MsgInvitationNotCosignedClaim.String():
		return types.MsgInvitationNotCosignedClaim
	default:
		return types.Unsupported
	}
}

func handlePegClaim(ctx context.Context, attributes []tmKv.Pair, cdc *amino.Codec, moniker, chainId, rpcUrl string) {
	// Parse attributes
	cosmosMsg, err := txs.PegClaimEventToCosmosMsg(attributes)
	if err != nil {
		logger.Error("Failed to convert PegClaim event to Cosmos Message", "err", err)
		return
	}

	// Validate
	status, err := proximaxClient.Transaction.GetTransactionStatus(ctx, cosmosMsg.MainchainTxHash)
	if err != nil {
		logger.Error("Transaction.GetTransaction returned error", "err", err)
		return
	}
	if status.Status != "Success" {
		logger.Error("Transaction status is not Success", "status", status.Status)
		return
	}
	if status.Group != "confirmed" {
		logger.Error("Transaction is not confirmed", "group", status.Group)
		return
	}

	// Broadcast Transaction
	txs.RelayPeg(chainId, cdc, moniker, cosmosMsg, rpcUrl)
	return
}

func handleUnpegNotCosignedClaim() {

}

func handleInvitationNotCosignedClaim() {

}
