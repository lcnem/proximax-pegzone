package proximax_bridge

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lcnem/proximax-pegzone/x/proximax-bridge/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) []abci.ValidatorUpdate {
	// TODO: Define logic for when you would like to initalize a new genesis
	k.SetParams(ctx, types.NewParams(data.MainchainMultisigAddress, data.Cosigners))

	return []abci.ValidatorUpdate{}
}

func GetGenesisStateFromAppState(cdc *codec.Codec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}
	return genesisState
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, k Keeper) (data types.GenesisState) {
	params := k.GetParams(ctx)

	// TODO: Define logic for exporting state
	return types.NewGenesisState(
		params.MainchainMultisigAddress, params.Cosigners,
	)
}
