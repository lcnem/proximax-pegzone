package proximax_bridge

import (
	"github.com/lcnem/proximax-pegzone/x/proximax-bridge/internal/keeper"
	"github.com/lcnem/proximax-pegzone/x/proximax-bridge/internal/types"
)

const (
	// TODO: define constants that you would like exposed from the internal package

	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	StoreKeyForPeg    = types.StoreKeyForPeg
	StoreKeyForUnpeg  = types.StoreKeyForUnpeg
	DefaultParamspace = types.DefaultParamspace
	QuerierRoute      = types.QuerierRoute
)

var (
	// functions aliases
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	RegisterCodec       = types.RegisterCodec
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	// TODO: Fill out function aliases
	NewMsgPeg                        = types.NewMsgPeg
	NewMsgPegClaim                   = types.NewMsgPegClaim
	NewMsgUnpeg                      = types.NewMsgUnpeg
	NewMsgRecordUnpeg                = types.NewMsgRecordUnpeg
	NewMsgUnpegNotCosignedClaim      = types.NewMsgUnpegNotCosignedClaim
	NewMsgRequestInvitation          = types.NewMsgRequestInvitation
	NewMsgInvitationNotCosignedClaim = types.NewMsgMsgInvitationNotCosignedClaim

	// variable aliases
	ModuleCdc = types.ModuleCdc
	// TODO: Fill out variable aliases
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	Params       = types.Params

	// TODO: Fill out module types
	MsgPeg                        = types.MsgPeg
	MsgPegClaim                   = types.MsgPegClaim
	MsgUnpeg                      = types.MsgUnpeg
	MsgRecordUnpeg                = types.MsgRecordUnpeg
	MsgUnpegNotCosignedClaim      = types.MsgUnpegNotCosignedClaim
	MsgRequestInvitation          = types.MsgRequestInvitation
	MsgInvitationNotCosignedClaim = types.MsgInvitationNotCosignedClaim
)
