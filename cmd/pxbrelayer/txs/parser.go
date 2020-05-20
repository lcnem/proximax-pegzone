package txs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmKv "github.com/tendermint/tendermint/libs/kv"

	msgTypes "github.com/lcnem/proximax-pegzone/x/proximax-bridge"
)

func PegEventToCosmosMsg(attributes []tmKv.Pair) (*msgTypes.MsgPeg, error) {
	var cosmosSender sdk.AccAddress
	var mainchainTxHash string
	var toAddress sdk.AccAddress
	var amount sdk.Coins
	var err error

	for _, attribute := range attributes {
		key := string(attribute.GetKey())
		val := string(attribute.GetValue())
		switch key {
		case "cosmos_sender":
			cosmosSender, err = sdk.AccAddressFromBech32(val)
			break
		case "mainchain_tx_hash":
			mainchainTxHash = val
			break
		case "cosmos_receiver":
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
	cosmosMsg := msgTypes.NewMsgPeg(cosmosSender, mainchainTxHash, toAddress, amount)
	return &cosmosMsg, nil
}

func UnpegEventToCosmosMsg(attributes []tmKv.Pair) (*msgTypes.MsgUnpeg, error) {
	var address sdk.AccAddress
	var cosmosAccount sdk.AccAddress
	var mainchainAddress string
	var amount sdk.Coins
	var firstCosignerAddress sdk.ValAddress
	var err error

	for _, attribute := range attributes {
		key := string(attribute.GetKey())
		val := string(attribute.GetValue())
		switch key {
		case "cosmos_sender":
			address, err = sdk.AccAddressFromBech32(val)
			if err != nil {
				return nil, err
			}
			break
		case "cosmos_account":
			cosmosAccount, err = sdk.AccAddressFromBech32(val)
			if err != nil {
				return nil, err
			}
			break
		case "mainchain_address":
			mainchainAddress = val
			break
		case "amount":
			amount, err = sdk.ParseCoins(val)
			if err != nil {
				return nil, err
			}
			break
		case "first_cosigner_address":
			firstCosignerAddress, err = sdk.ValAddressFromBech32(val)
			if err != nil {
				return nil, err
			}
		}
	}
	cosmosMsg := msgTypes.NewMsgUnpeg(address, cosmosAccount, mainchainAddress, amount, firstCosignerAddress)
	return &cosmosMsg, nil
}

func RequestInvitationEventToCosmosMsg(attributes []tmKv.Pair) (*msgTypes.MsgRequestInvitation, string, error) {
	var address sdk.ValAddress
	var multisigAccountPublicKey string
	var newCosignerPublicKey string
	var firstCosignerAddress sdk.ValAddress
	var err error

	for _, attribute := range attributes {
		key := string(attribute.GetKey())
		val := string(attribute.GetValue())
		switch key {
		case "cosmos_account":
			address, err = sdk.ValAddressFromBech32(val)
			if err != nil {
				break
			}
			break
		case "multisig_address":
			multisigAccountPublicKey = val
			break
		case "new_cosigner_public_key":
			newCosignerPublicKey = val
			break
		case "first_cosigner_address":
			firstCosignerAddress, err = sdk.ValAddressFromBech32(val)
			if err != nil {
				break
			}
		}
	}
	if err != nil {
		return nil, "", err
	}
	cosmosMsg := msgTypes.NewMsgRequestInvitation(address, newCosignerPublicKey, firstCosignerAddress)
	return &cosmosMsg, multisigAccountPublicKey, nil
}
