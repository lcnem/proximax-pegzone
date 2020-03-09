package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/peggy/x/oracle"
	"github.com/lcnem/proximax-pegzone/x/proximax-bridge/internal/types"
)

// Keeper of the proximax-bridge store
type Keeper struct {
	storeKey       sdk.StoreKey
	cdc            *codec.Codec
	paramspace     types.ParamSubspace
	supplyKeeper   types.SupplyKeeper
	slashingKeeper types.SlashingKeeper
	oracleKeeper   types.OracleKeeper
}

// NewKeeper creates a proximax-bridge keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramspace types.ParamSubspace, supplyKeeper types.SupplyKeeper, slashingKeeper types.SlashingKeeper, oracleKeeper types.OracleKeeper) Keeper {
	keeper := Keeper{
		storeKey:       key,
		cdc:            cdc,
		paramspace:     paramspace.WithKeyTable(types.ParamKeyTable()),
		supplyKeeper:   supplyKeeper,
		slashingKeeper: slashingKeeper,
		oracleKeeper:   oracleKeeper,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// ProcessClaim processes a new claim coming in from a validator
func (k Keeper) ProcessPegClaim(ctx sdk.Context, claim types.MsgPegClaim) (oracle.Status, error) {
	oracleClaim, err := types.CreateOracleClaimFromMsgPegClaim(k.cdc, claim)
	if err != nil {
		return oracle.Status{}, err
	}

	return k.oracleKeeper.ProcessClaim(ctx, oracleClaim)
}

// ProcessSuccessfulClaim processes a claim that has just completed successfully with consensus
func (k Keeper) ProcessSuccessfulPegClaim(ctx sdk.Context, claim string) error {
	oracleClaim, err := types.CreateMsgPegClaimFromOracleString(claim)
	if err != nil {
		return err
	}

	err = k.supplyKeeper.MintCoins(ctx, types.ModuleName, oracleClaim.Amount)

	if err != nil {
		return err
	}

	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, oracleClaim.ToAddress, oracleClaim.Amount,
	); err != nil {
		panic(err)
	}

	return nil
}

// ProcessClaim processes a new claim coming in from a validator
func (k Keeper) ProcessUnpegNotCosignedClaim(ctx sdk.Context, claim types.MsgUnpegNotCosignedClaim) (oracle.Status, error) {
	oracleClaim, err := types.CreateOracleClaimFromMsgUnpegNotCosignedClaim(k.cdc, claim)
	if err != nil {
		return oracle.Status{}, err
	}

	return k.oracleKeeper.ProcessClaim(ctx, oracleClaim)
}

// ProcessSuccessfulClaim processes a claim that has just completed successfully with consensus
func (k Keeper) ProcessSuccessfulUnpegNotCosignedClaim(ctx sdk.Context, claim string) error {
	oracleClaim, err := types.CreateMsgUnpegNotCosignedClaimFromOracleString(claim)
	if err != nil {
		return err
	}

	for _, notCosignedValidator := range oracleClaim.NotCosignedValidators {
		k.slashingKeeper.Slash(ctx, sdk.ConsAddress(notCosignedValidator), sdk.NewDec(0), 0, 0)
	}

	return nil
}

// ProcessClaim processes a new claim coming in from a validator
func (k Keeper) ProcessInvitationNotCosignedClaim(ctx sdk.Context, claim types.MsgInvitationNotCosignedClaim) (oracle.Status, error) {
	oracleClaim, err := types.CreateOracleClaimFromMsgInvitationNotCosignedClaim(k.cdc, claim)
	if err != nil {
		return oracle.Status{}, err
	}

	return k.oracleKeeper.ProcessClaim(ctx, oracleClaim)
}

// ProcessSuccessfulClaim processes a claim that has just completed successfully with consensus
func (k Keeper) ProcessSuccessfulInvitationNotCosignedClaim(ctx sdk.Context, claim string) error {
	oracleClaim, err := types.CreateMsgInvitationNotCosignedClaimFromOracleString(claim)
	if err != nil {
		return err
	}

	for _, notCosignedValidator := range oracleClaim.NotCosignedValidators {
		k.slashingKeeper.Slash(ctx, sdk.ConsAddress(notCosignedValidator), sdk.NewDec(0), 0, 0)
	}

	return nil
}
