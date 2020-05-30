package proximax_bridge

import (
	"github.com/lcnem/proximax-pegzone/x/proximax-bridge/keeper"
	"github.com/lcnem/proximax-pegzone/x/proximax-bridge/types"
)

const (
	// TODO: define constants that you would like exposed from the internal package

	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	StoreKeyForPeg    = types.StoreKeyForPeg
	StoreKeyForUnpeg  = types.StoreKeyForUnpeg
	StoreKeyForCosign = types.StoreKeyForCosign
	StoreKeyForInvite = types.StoreKeyForInvite
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
	NewMsgPeg                      = types.NewMsgPeg
	NewMsgPegClaim                 = types.NewMsgPegClaim
	NewMsgUnpeg                    = types.NewMsgUnpeg
	NewMsgRecordUnpeg              = types.NewMsgRecordUnpeg
	NewMsgNotifyCosigned           = types.NewMsgNotifyCosigned
	NewMsgNotCosignedClaim         = types.NewMsgNotCosignedClaim
	NewMsgRequestInvitation        = types.NewMsgRequestInvitation
	NewMsgPendingRequestInvitation = types.NewMsgPendingRequestInvitation
	NewMsgConfirmedInvitation      = types.NewMsgConfirmedInvitation

	// variable aliases
	ModuleCdc = types.ModuleCdc
	// TODO: Fill out variable aliases
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	Params       = types.Params

	// TODO: Fill out module types
	MsgPeg                      = types.MsgPeg
	MsgPegClaim                 = types.MsgPegClaim
	MsgUnpeg                    = types.MsgUnpeg
	MsgRecordUnpeg              = types.MsgRecordUnpeg
	MsgNotifyCosigned           = types.MsgNotifyCosigned
	MsgNotCosignedClaim         = types.MsgNotCosignedClaim
	MsgRequestInvitation        = types.MsgRequestInvitation
	MsgPendingRequestInvitation = types.MsgPendingRequestInvitation
	MsgConfirmedInvitation      = types.MsgConfirmedInvitation

	Cosigner = types.Cosigner
)
