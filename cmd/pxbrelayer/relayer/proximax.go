package relayer

import (
	"context"
	"time"

	sdkContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	proximax "github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket"
)

func InitProximaXRelayer(
	cdc *codec.Codec,
	chainID string,
	provider string,
	makeClaims bool,
	validatorName string,
	passphrase string,
	validatorAddress sdk.ValAddress,
	cliContext sdkContext.CLIContext,
	rpcURL string,
	custodyAddress string,
	test bool,
) error {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*30)

	conf, err := proximax.NewConfig(ctx, []string{""})
	if err != nil {
		return err
	}

	if test {
		conf.NetworkType = proximax.PublicTest
	} else {
		conf.NetworkType = proximax.Public
	}

	client, err := websocket.NewClient(ctx, conf)
	if err != nil {
		return err
	}
	go client.Listen()

	address, err := proximax.NewAddressFromRaw(custodyAddress)
	if err != nil {
		return err
	}

	err = client.AddConfirmedAddedHandlers(address, func(tx proximax.Transaction) bool {

		return false
	})

	for {

	}

	cancel()
}
