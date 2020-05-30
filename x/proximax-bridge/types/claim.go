package types

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/peggy/x/oracle"
)

func CreateOracleClaimFromMsgPegClaim(cdc *codec.Codec, msg MsgPegClaim) (oracle.Claim, error) {
	oracleID := fmt.Sprintf("%s,%s,%d", msg.MainchainTxHash, msg.Amount.String(), msg.Remainning)
	claimBytes, err := json.Marshal(msg)
	if err != nil {
		return oracle.Claim{}, err
	}
	claimString := string(claimBytes)
	claim := oracle.NewClaim(oracleID, msg.ValidatorAddress, claimString)
	return claim, nil
}

func CreateOracleClaimFromMsgNotCosignedClaim(cdc *codec.Codec, msg MsgNotCosignedClaim) (oracle.Claim, error) {
	oracleID := msg.TxHash
	claimBytes, err := json.Marshal(msg)
	if err != nil {
		return oracle.Claim{}, err
	}
	claimString := string(claimBytes)
	claim := oracle.NewClaim(oracleID, msg.Address, claimString)
	return claim, nil
}

// CreateOracleClaimFromOracleString converts a JSON string into an OracleClaimContent struct used by this module.
// In general, it is expected that the oracle module will store claims in this JSON format
// and so this should be used to convert oracle claims.
func CreateMsgPegClaimFromOracleString(oracleClaimString string) (MsgPegClaim, error) {
	var oracleClaim MsgPegClaim

	bz := []byte(oracleClaimString)
	if err := json.Unmarshal(bz, &oracleClaim); err != nil {
		return MsgPegClaim{}, sdkerrors.Wrap(ErrJSONMarshalling, fmt.Sprintf("failed to parse claim: %s", err.Error()))
	}

	return oracleClaim, nil
}

// CreateOracleClaimFromOracleString converts a JSON string into an OracleClaimContent struct used by this module.
// In general, it is expected that the oracle module will store claims in this JSON format
// and so this should be used to convert oracle claims.
func CreateMsgNotCosignedClaimFromOracleString(oracleClaimString string) (MsgNotCosignedClaim, error) {
	var oracleClaim MsgNotCosignedClaim

	bz := []byte(oracleClaimString)
	if err := json.Unmarshal(bz, &oracleClaim); err != nil {
		return MsgNotCosignedClaim{}, sdkerrors.Wrap(ErrJSONMarshalling, fmt.Sprintf("failed to parse claim: %s", err.Error()))
	}

	return oracleClaim, nil
}
