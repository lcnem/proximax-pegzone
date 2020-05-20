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
	storeKey          sdk.StoreKey
	storeKeyForPeg    sdk.StoreKey
	storeKeyForCosign sdk.StoreKey
	cdc               *codec.Codec
	paramspace        types.ParamSubspace
	supplyKeeper      types.SupplyKeeper
	slashingKeeper    types.SlashingKeeper
	oracleKeeper      types.OracleKeeper
}

// NewKeeper creates a proximax-bridge keeper
func NewKeeper(cdc *codec.Codec, key, keyForPeg, keyForCosign sdk.StoreKey, paramspace types.ParamSubspace, supplyKeeper types.SupplyKeeper, slashingKeeper types.SlashingKeeper, oracleKeeper types.OracleKeeper) Keeper {
	keeper := Keeper{
		storeKey:          key,
		storeKeyForPeg:    keyForPeg,
		storeKeyForCosign: keyForCosign,
		cdc:               cdc,
		paramspace:        paramspace.WithKeyTable(types.ParamKeyTable()),
		supplyKeeper:      supplyKeeper,
		slashingKeeper:    slashingKeeper,
		oracleKeeper:      oracleKeeper,
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

type CosignersRecord struct {
	MainchainTxHadh    string   `json:"mainchain_tx_hash" yaml:"mainchain_tx_hash"`
	CosignerPublicKeys []string `json:"cosigner_public_keys" yaml:"cosigner_public_keys"`
}

func (k Keeper) SetCosigners(ctx sdk.Context, mainChainTxHash string, cosignerPublicKey string) error {
	cosignerRecord, err := k.GetCosignersRecord(ctx, mainChainTxHash)
	if err != nil {
		cosignerRecord = CosignersRecord{MainchainTxHadh: mainChainTxHash, CosignerPublicKeys: []string{}}
	}
	for _, key := range cosignerRecord.CosignerPublicKeys {
		if key == cosignerPublicKey {
			return nil
		}
	}
	cosignerRecord.CosignerPublicKeys = append(cosignerRecord.CosignerPublicKeys, cosignerPublicKey)
	cosignerRecordBytes, err := json.Marshal(cosignerRecord)
	if err != nil {
		return err
	}
	ctx.KVStore(k.storeKeyForCosign).Set([]byte(mainChainTxHash), cosignerRecordBytes)
	return nil
}

func (k Keeper) GetCosignersRecord(ctx sdk.Context, mainChainTxHash string) (CosignersRecord, error) {
	cosignersRecord := CosignersRecord{}
	if !ctx.KVStore(k.storeKeyForCosign).Has([]byte(mainChainTxHash)) {
		return cosignersRecord, errors.New(fmt.Sprintf("CosignersRecord Record is Not Found: %s", mainChainTxHash))
	}
	unpegBytes := ctx.KVStore(k.storeKeyForCosign).Get([]byte(mainChainTxHash))
	err := json.Unmarshal(unpegBytes, &cosignersRecord)
	return cosignersRecord, err
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
		ctx, types.ModuleName, oracleClaim.Address, oracleClaim.Amount,
	); err != nil {
		panic(err)
	}

	return nil
}

func (k Keeper) ProcessUnpeg(ctx sdk.Context, msg types.MsgUnpeg) error {
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(
		ctx, msg.Address, types.ModuleName, msg.Amount,
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

func searchStringFromArray(values []string, key string) bool {
	for _, value := range values {
		if value == key {
			return true
		}
	}
	return false
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

	cosignerRecord, err := k.GetCosignersRecord(ctx, oracleClaim.TxHash)
	if err != nil {
		return err
	}

	param := k.GetParams(ctx)

	notCosignedValidatorAddrs := []sdk.ValAddress{}
	for _, cosigner := range param.Cosigners {
		if !searchStringFromArray(cosignerRecord.CosignerPublicKeys, cosigner.MainchainPublicKey) {
			valAddress, err := sdk.ValAddressFromBech32(cosigner.ValidatorAddress)
			if err == nil {
				notCosignedValidatorAddrs = append(notCosignedValidatorAddrs, valAddress)
			}
		}
	}

	for _, notCosignedValidator := range notCosignedValidatorAddrs {
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
