package txs

import (
	"strings"

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

func UnpegNotCosignedClaimEventToCosmosMsg(attributes []tmKv.Pair) (*msgTypes.MsgUnpegNotCosignedClaim, error) {
	var address sdk.ValAddress
	var txHash string
	var notCosignedValidators []sdk.ValAddress
	var err error

	for _, attribute := range attributes {
		key := string(attribute.GetKey())
		val := string(attribute.GetValue())
		switch key {
		case "cosmos_sender":
			address, err = sdk.ValAddressFromBech32(val)
			break
		case "tx_hash":
			txHash = val
			break
		case "not_cosigned_validators":
			for _, addr := range strings.Split(val, ",") {
				valAddress, err := sdk.ValAddressFromBech32(addr)
				if err != nil {
					break
				}
				notCosignedValidators = append(notCosignedValidators, valAddress)
			}
			break
		}
	}
	if err != nil {
		return nil, err
	}
	cosmosMsg := msgTypes.NewMsgUnpegNotCosignedClaim(address, txHash, notCosignedValidators)
	return &cosmosMsg, nil
}

func InvitationNotCosignedClaimEventToCosmosMsg(attributes []tmKv.Pair) (*msgTypes.MsgInvitationNotCosignedClaim, error) {
	var address sdk.ValAddress
	var txHash string
	var notCosignedValidators []sdk.ValAddress
	var err error

	for _, attribute := range attributes {
		key := string(attribute.GetKey())
		val := string(attribute.GetValue())
		switch key {
		case "cosmos_sender":
			address, err = sdk.ValAddressFromBech32(val)
			break
		case "tx_hash":
			txHash = val
			break
		case "not_cosigned_validators":
			for _, addr := range strings.Split(val, ",") {
				valAddress, err := sdk.ValAddressFromBech32(addr)
				if err != nil {
					break
				}
				notCosignedValidators = append(notCosignedValidators, valAddress)
			}
			break
		}
	}
	if err != nil {
		return nil, err
	}
	cosmosMsg := msgTypes.NewMsgInvitationNotCosignedClaim(address, txHash, notCosignedValidators)
	return &cosmosMsg, nil
}

func UnpegEventToCosmosMsg(attributes []tmKv.Pair) (*msgTypes.MsgUnpeg, error) {
	var address sdk.AccAddress
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
				break
			}
			break
		case "mainchain_address":
			mainchainAddress = val
			break
		case "amount":
			amount, err = sdk.ParseCoins(val)
			if err != nil {
				break
			}
			break
		case "first_cosigner_address":
			firstCosignerAddress, err = sdk.ValAddressFromBech32(val)
			if err != nil {
				break
			}
		}
	}
	if err != nil {
		return nil, err
	}
	cosmosMsg := msgTypes.NewMsgUnpeg(address, mainchainAddress, amount, firstCosignerAddress)
	return &cosmosMsg, nil
}

func RequestInvitationEventToCosmosMsg(attributes []tmKv.Pair) (*msgTypes.MsgRequestInvitation, error) {
	var address sdk.ValAddress
	var mainchainAddress string
	var firstCosignerAddress sdk.ValAddress
	var err error

	for _, attribute := range attributes {
		key := string(attribute.GetKey())
		val := string(attribute.GetValue())
		switch key {
		case "address":
			address, err = sdk.ValAddressFromBech32(val)
			if err != nil {
				break
			}
			break
		case "mainchain_address":
			mainchainAddress = val
			break
		case "first_cosigner_address":
			firstCosignerAddress, err = sdk.ValAddressFromBech32(val)
			if err != nil {
				break
			}
		}
	}
	if err != nil {
		return nil, err
	}
	cosmosMsg := msgTypes.NewMsgRequestInvitation(address, mainchainAddress, firstCosignerAddress)
	return &cosmosMsg, nil
}
