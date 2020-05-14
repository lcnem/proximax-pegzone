package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	// TODO: Register the modules msgs
	cdc.RegisterConcrete(MsgPeg{}, "proximaxbridge/MsgPeg", nil)
	cdc.RegisterConcrete(MsgPegClaim{}, "proximaxbridge/MsgPegClaim", nil)
	cdc.RegisterConcrete(MsgUnpeg{}, "proximaxbridge/MsgUnpeg", nil)
	cdc.RegisterConcrete(MsgRecordUnpeg{}, "proximaxbridge/MsgRecordUnpeg", nil)
	cdc.RegisterConcrete(MsgRequestInvitation{}, "proximaxbridge/MsgRequestInvitation", nil)
	cdc.RegisterConcrete(MsgUnpegNotCosignedClaim{}, "proximaxbridge/MsgUnpegNotCosignedClaim", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
