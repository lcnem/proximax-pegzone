package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Query endpoints supported by the proximax-bridge querier
const (
	//TODO: Describe query parameters, update <action> with your query
	// Query<Action>    = "<action>"
	QueryMainchainMultisigAddress = "mainchain_multisig_address"
	QueryCosigners                = "cosigners"
)

/*
Below you will be able how to set your own queries:


// QueryResList Queries Result Payload for a query
type QueryResList []string

// implement fmt.Stringer
func (n QueryResList) String() string {
	return strings.Join(n[:], "\n")
}

*/

type QueryResMainchainMultisigAddress struct {
	MainchainMultisigAddress string `json:"mainchain_multisig_address"`
}

type Cosigner struct {
	ValidatorAddress sdk.ValAddress `json:"validator_address"`
	MainchainAddress string         `json:"mainchain_address"`
}
