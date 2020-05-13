package keeper

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/peggy/x/oracle"
	"github.com/lcnem/proximax-pegzone/x/proximax-bridge/internal/types"
)

// Keeper of the proximax-bridge store
type Keeper struct {
	storeKey         sdk.StoreKey
	storeKeyForPeg   sdk.StoreKey
	storeKeyForUnpeg sdk.StoreKey
	cdc              *codec.Codec
	paramspace       types.ParamSubspace
	supplyKeeper     types.SupplyKeeper
	slashingKeeper   types.SlashingKeeper
	oracleKeeper     types.OracleKeeper
}

// NewKeeper creates a proximax-bridge keeper
func NewKeeper(cdc *codec.Codec, key, keyForPeg, keyForUnpeg sdk.StoreKey, paramspace types.ParamSubspace, supplyKeeper types.SupplyKeeper, slashingKeeper types.SlashingKeeper, oracleKeeper types.OracleKeeper) Keeper {
	keeper := Keeper{
		storeKey:         key,
		storeKeyForPeg:   keyForPeg,
		storeKeyForUnpeg: keyForUnpeg,
		cdc:              cdc,
		paramspace:       paramspace.WithKeyTable(types.ParamKeyTable()),
		supplyKeeper:     supplyKeeper,
		slashingKeeper:   slashingKeeper,
		oracleKeeper:     oracleKeeper,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) IsUsedHash(ctx sdk.Context, hash string) bool {
	return ctx.KVStore(k.storeKeyForPeg).Has([]byte(hash))
}

func (k Keeper) MarkAsUsedHash(ctx sdk.Context, hash string) {
	ctx.KVStore(k.storeKeyForPeg).Set([]byte(hash), []byte(hash))
}

type UnpegRecord struct {
	Address         sdk.AccAddress `json:"address" yaml:"address"`
	MainchainTxHash string         `json:"mainchain_tx_hash" yaml:"mainchain_tx_hash"`
	Amount          sdk.Coins      `json:"amount" yaml:"amount"`
}

func (k Keeper) SetUnpegRecord(ctx sdk.Context, mainChainTxHash string, accountAddress sdk.AccAddress, amount sdk.Coins) error {
	unpeg := UnpegRecord{Address: accountAddress, MainchainTxHash: mainChainTxHash, Amount: amount}
	unpegBytes, err := json.Marshal(unpeg)
	if err != nil {
		return err
	}
	ctx.KVStore(k.storeKeyForPeg).Set([]byte(mainChainTxHash), unpegBytes)
	return nil
}

func (k Keeper) GetUnpegRecord(ctx sdk.Context, mainChainTxHash string) (UnpegRecord, error) {
	unpeg := UnpegRecord{}
	if !ctx.KVStore(k.storeKeyForPeg).Has([]byte(mainChainTxHash)) {
		return unpeg, errors.New(fmt.Sprintf("Unpeg Record is Not Found: %s", mainChainTxHash))
	}
	unpegBytes := ctx.KVStore(k.storeKeyForPeg).Get([]byte(mainChainTxHash))
	err := json.Unmarshal(unpegBytes, &unpeg)
	return unpeg, err
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

	if err := k.supplyKeeper.MintCoins(
		ctx, types.ModuleName, oracleClaim.Amount,
	); err != nil {
		return err
	}

	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, oracleClaim.ToAddress, oracleClaim.Amount,
	); err != nil {
		panic(err)
	}

	return nil
}

func (k Keeper) ProcessUnpeg(ctx sdk.Context, msg types.MsgUnpeg) error {
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(
		ctx, msg.FromAddress, types.ModuleName, msg.Amount,
	); err != nil {
		return err
	}

	if err := k.supplyKeeper.BurnCoins(
		ctx, types.ModuleName, msg.Amount,
	); err != nil {
		return err
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
	unpegRecord, err := k.GetUnpegRecord(ctx, oracleClaim.TxHash)
	if err != nil {
		return err
	}

	if err := k.supplyKeeper.MintCoins(
		ctx, types.ModuleName, unpegRecord.Amount,
	); err != nil {
		return err
	}

	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, sdk.AccAddress(unpegRecord.Address), unpegRecord.Amount,
	); err != nil {
		panic(err)
	}
	/*
		for _, notCosignedValidator := range oracleClaim.NotCosignedValidators {
			k.slashingKeeper.Slash(ctx, sdk.ConsAddress(notCosignedValidator), sdk.NewDec(0), 0, 0)
		}
	*/

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
