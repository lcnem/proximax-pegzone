package keeper

// TODO: Define if your module needs Parameters, if not this can be deleted

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lcnem/proximax-pegzone/x/proximax-bridge/internal/types"
)

// GetParams returns the total set of proximax-bridge parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramspace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the proximax-bridge parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramspace.SetParamSet(ctx, &params)
}

func (k Keeper) AddNewCosigner(ctx sdk.Context, address sdk.ValAddress, mainchainPublicKey string) {
	params := k.GetParams(ctx)

	for _, cosigner := range params.Cosigners {
		if cosigner.MainchainPublicKey == mainchainPublicKey {
			return
		}
	}
	params.Cosigners = append(params.Cosigners, types.Cosigner{ValidatorAddress: address.String(), MainchainPublicKey: mainchainPublicKey})
	k.SetParams(ctx, params)
}
