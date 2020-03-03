package keeper

import (
	"github.com/lcnem/proximax-pegzone/x/proximax-bridge/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewQuerier creates a new querier for proximax-bridge clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryMainchainMultisigAddress:
			return queryMainchainMultisigAddress(ctx, k)
		case types.QueryCosigners:
			return queryCosigners(ctx, k)
		//case types.QueryParams:
		//	return queryParams(ctx, k)
		// TODO: Put the modules query routes
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown proximax-bridge query endpoint")
		}
	}
}

// TODO: Add the modules query functions
// They will be similar to the above one: queryParams()

func queryMainchainMultisigAddress(ctx sdk.Context, k Keeper) ([]byte, error) {
	address, err := k.GetMainchainMultisigAddress(ctx)
	if err != nil {
		address = ""
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, types.QueryResMainchainMultisigAddress{MainchainMultisigAddress: address})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryCosigners(ctx sdk.Context, k Keeper) ([]byte, error) {
	cosigners, err := k.GetCosigners(ctx)
	if err != nil {
		cosigners = []types.Cosigner{}
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, cosigners)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
