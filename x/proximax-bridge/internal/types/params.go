package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/x/params"
)

// Default parameter namespace
const (
	DefaultParamspace = ModuleName
	// TODO: Define your default parameters
)

// Parameter store keys
var (
	// TODO: Define your keys for the parameter store
	// KeyParamName          = []byte("ParamName")
	KeyMainchainMultisigAddress = []byte("MainchainMultisigAddress")
	KeyCosigners                = []byte("Cosigners")
)

// ParamKeyTable for proximax-bridge module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Params - used for initializing default parameter for proximax-bridge at genesis
type Params struct {
	// TODO: Add your Paramaters to the Paramter struct
	// KeyParamName string `json:"key_param_name"`
	MainchainMultisigAddress string     `json:"mainchain_address"`
	Cosigners                []Cosigner `json:"cosigners"`
}

type Cosigner struct {
	ValidatorAddress   string `json:"validator_address"`
	MainchainPublicKey string `json:"mainchain_public_key"`
}

// NewParams creates a new Params object
func NewParams(mainchainMultisigAddress string, cosigners []Cosigner) Params {
	return Params{
		// TODO: Create your Params Type
		MainchainMultisigAddress: mainchainMultisigAddress,
		Cosigners:                cosigners,
	}
}

// String implements the stringer interface for Params
func (p Params) String() string {
	value, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(value)
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		// TODO: Pair your key with the param
		// params.NewParamSetPair(KeyParamName, &p.ParamName),
		params.NewParamSetPair(KeyMainchainMultisigAddress, &p.MainchainMultisigAddress, func(value interface{}) error { return nil }),
		params.NewParamSetPair(KeyCosigners, &p.Cosigners, func(value interface{}) error { return nil }),
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams("", []Cosigner{})
}
