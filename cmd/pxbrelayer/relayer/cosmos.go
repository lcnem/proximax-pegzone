package relayer

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	tmLog "github.com/tendermint/tendermint/libs/log"
	tmClient "github.com/tendermint/tendermint/rpc/client"
	tmTypes "github.com/tendermint/tendermint/types"
)

func InitCosmosRelayer(
	tendermintNode string,
	proximaxNode string,
) error {
	logger := tmLog.NewTMLogger(tmLog.NewSyncWriter(os.Stdout))
	client, err := tmClient.NewHTTP(tendermintNode, "/websocket")

	client.SetLogger(logger)

	err = client.Start()
	if err != nil {
		logger.Error("Failed to start a client", "err", err)
		os.Exit(1)
	}

	defer client.Stop()

	// Subscribe to all tendermint transactions
	query := "tm.event = 'Tx'"

	out, err := client.Subscribe(context.Background(), "test", query, 1000)
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
				switch event.Type {

				}
			}
		case <-quit:
			os.Exit(0)
		}
	}
}
