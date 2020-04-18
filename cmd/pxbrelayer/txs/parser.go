package txs

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmKv "github.com/tendermint/tendermint/libs/kv"

	"github.com/lcnem/proximax-pegzone/cmd/pxbrelayer/types"
	msgTypes "github.com/lcnem/proximax-pegzone/x/proximax-bridge"
)

func PegClaimEventToCosmosMsg(attributes []tmKv.Pair) (*msgTypes.MsgPegClaim, error) {
	var cosmosSender sdk.ValAddress
	var mainchainTxHash string
	var toAddress sdk.AccAddress
	var amount sdk.Coins
	var err error

	for _, attribute := range attributes {
		key := string(attribute.GetKey())
		val := string(attribute.GetValue())

		switch key {
		case types.CosmosSender.String():
			cosmosSender, err = sdk.ValAddressFromBech32(val)
			if err != nil {
				msg := fmt.Sprintf("Invalid Validator Address: %s", val)
				return nil, errors.New(msg)
			}
			break
		case types.MainchainTxHash.String():
			mainchainTxHash = val
			break
		case types.ToAddress.String():
			toAddress, err = sdk.AccAddressFromBech32(val)
			if err != nil {
				msg := fmt.Sprintf("Invalid Account Address: %s", val)
				return nil, errors.New(msg)
			}
			break
		case types.Coin.String():
			amount, err = sdk.ParseCoins(val)
			if err != nil {
				msg := fmt.Sprintf("Invalid Coins: %s", val)
				return nil, errors.New(msg)
			}
			break
		}
	}
	cosmosMsg := msgTypes.NewMsgPegClaim(cosmosSender, mainchainTxHash, toAddress, amount)
	return &cosmosMsg, nil
}
