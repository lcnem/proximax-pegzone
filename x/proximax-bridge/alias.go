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
	NewMsgPegClaim                   = types.NewMsgPegClaim
	NewMsgUnpeg                      = types.NewMsgUnpeg
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
	MsgPegClaim                   = types.MsgPegClaim
	MsgUnpeg                      = types.MsgUnpeg
	MsgUnpegNotCosignedClaim      = types.MsgUnpegNotCosignedClaim
	MsgRequestInvitation          = types.MsgRequestInvitation
	MsgInvitationNotCosignedClaim = types.MsgMsgInvitationNotCosignedClaim
)
