package types

// GenesisState - all proximax-bridge state that must be provided at genesis
type GenesisState struct {
	// TODO: Fill out what is needed by the module for genesis
	MainchainMultisigAddress string     `json:"mainchain_multisig_address"`
	Cosigners                []Cosigner `json:"cosigners"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	/* TODO: Fill out with what is needed for genesis state*/
	mainchainMultisigAddress string,
	cosigners []Cosigner,
) GenesisState {

	return GenesisState{
		// TODO: Fill out according to your genesis state
		MainchainMultisigAddress: mainchainMultisigAddress,
		Cosigners:                cosigners,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		// TODO: Fill out according to your genesis state, these values will be initialized but empty
		MainchainMultisigAddress: "",
		Cosigners:                []Cosigner{},
	}
}

// ValidateGenesis validates the proximax-bridge genesis parameters
func ValidateGenesis(data GenesisState) error {
	// TODO: Create a sanity check to make sure the state conforms to the modules needs
	return nil
}
