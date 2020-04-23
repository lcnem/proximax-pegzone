package txs

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmKv "github.com/tendermint/tendermint/libs/kv"

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
		case "cosmos_sender":
			cosmosSender, err = sdk.ValAddressFromBech32(val)
			break
		case "mainchain_tx_hash":
			mainchainTxHash = val
			break
		case "to_address":
			toAddress, err = sdk.AccAddressFromBech32(val)
			break
		case "amount":
			amount, err = sdk.ParseCoins(val)
			break
		}
	}
	if err != nil {
		return nil, err
	}
	cosmosMsg := msgTypes.NewMsgPegClaim(cosmosSender, mainchainTxHash, toAddress, amount)
	return &cosmosMsg, nil
}

