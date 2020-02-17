package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// TODO: Describe your actions, these will implment the interface of `sdk.Msg`
// verify interface at compile time
var _ sdk.Msg = &MsgPegClaim{}

// MsgPegClaim - struct for unjailing jailed validator
type MsgPegClaim struct {
	Address         sdk.AccAddress `json:"address" yaml:"address"`
	MainchainTxHash string         `json:"mainchain_tx_hash" yaml:"mainchain_tx_hash"`
	Amount          []sdk.Coin     `json:"amount" yaml:"amount"`
}

// NewMsgPegClaim creates a new MsgPegClaim instance
func NewMsgPegClaim(validatorAddr sdk.ValAddress) MsgPegClaim {
	return MsgPegClaim{
		ValidatorAddr: validatorAddr,
	}
}

const pegClaimConst = "pegClaim"

// nolint
func (msg MsgPegClaim) Route() string { return RouterKey }
func (msg MsgPegClaim) Type() string  { return pegClaimConst }
func (msg MsgPegClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddr)}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgPegClaim) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgPegClaim) ValidateBasic() error {
	if msg.ValidatorAddr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing validator address")
	}
	return nil
}

// TODO: Describe your actions, these will implment the interface of `sdk.Msg`
// verify interface at compile time
var _ sdk.Msg = &MsgUnpeg{}

// MsgUnpeg - struct for unjailing jailed validator
type MsgUnpeg struct {
	Address              sdk.AccAddress `json:"address" yaml:"address"`
	MainchainAddress     string         `json:"mainchain_address" yaml:"mainchain_address"`
	Amount               []sdk.Coin     `json:"amount" yaml:"amount"`
	FirstCosignerAddress sdk.ValAddress `json:"first_cosigner_address" yaml:"first_cosigner_address"`
}

// NewMsgUnpeg creates a new MsgUnpeg instance
func NewMsgUnpeg(validatorAddr sdk.ValAddress) MsgUnpeg {
	return MsgUnpeg{
		ValidatorAddr: validatorAddr,
	}
}

const unpegConst = "unpeg"

// nolint
func (msg MsgUnpeg) Route() string { return RouterKey }
func (msg MsgUnpeg) Type() string  { return unpegConst }
func (msg MsgUnpeg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddr)}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgUnpeg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgUnpeg) ValidateBasic() error {
	if msg.ValidatorAddr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing validator address")
	}
	return nil
}

// TODO: Describe your actions, these will implment the interface of `sdk.Msg`
// verify interface at compile time
var _ sdk.Msg = &MsgUnpegNotCosignedClaim{}

// MsgUnpegNotCosignedClaim - struct for unjailing jailed validator
type MsgUnpegNotCosignedClaim struct {
	ValidatorAddress     sdk.ValAddress `json:"validator_address" yaml: "validator_address"`
	TxHash               string         `json:"tx_hash" yaml:"tx_hash"`
	FirstCosignerAddress sdk.ValAddress `json:"first_cosigner_address" yaml:"first_cosigner_address"`
}

// NewMsgUnpegNotCosignedClaim creates a new MsgUnpegNotCosignedClaim instance
func NewMsgUnpegNotCosignedClaim(validatorAddr sdk.ValAddress) MsgUnpegNotCosignedClaim {
	return MsgUnpegNotCosignedClaim{
		ValidatorAddr: validatorAddr,
	}
}

const unpegNotCosignedClaimConst = "unpegNotCosignedClaim"

// nolint
func (msg MsgUnpegNotCosignedClaim) Route() string { return RouterKey }
func (msg MsgUnpegNotCosignedClaim) Type() string  { return unpegNotCosignedClaimConst }
func (msg MsgUnpegNotCosignedClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddr)}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgUnpegNotCosignedClaim) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgUnpegNotCosignedClaim) ValidateBasic() error {
	if msg.ValidatorAddr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing validator address")
	}
	return nil
}

// TODO: Describe your actions, these will implment the interface of `sdk.Msg`

// verify interface at compile time
var _ sdk.Msg = &MsgRequestInvitation{}

// MsgRequestInvitation - struct for unjailing jailed validator
type MsgRequestInvitation struct {
	ValidatorAddress     sdk.ValAddress `json:"validator_address" yaml: "validator_address"`
	MainchainAddress     string         `json:"mainchain_address" yaml:"mainchain_address"`
	FirstCosignerAddress sdk.ValAddress `json:"first_cosigner_address" yaml:"first_cosigner_address"`
}

// NewMsgRequestInvitation creates a new MsgRequestInvitation instance
func NewMsgRequestInvitation(validatorAddr sdk.ValAddress) MsgRequestInvitation {
	return MsgRequestInvitation{
		ValidatorAddr: validatorAddr,
	}
}

const requestInvitationConst = "requestInvitation"

// nolint
func (msg MsgRequestInvitation) Route() string { return RouterKey }
func (msg MsgRequestInvitation) Type() string  { return requestInvitationConst }
func (msg MsgRequestInvitation) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddr)}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgRequestInvitation) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgRequestInvitation) ValidateBasic() error {
	if msg.ValidatorAddr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing validator address")
	}
	return nil
}

// TODO: Describe your actions, these will implment the interface of `sdk.Msg`

// verify interface at compile time
var _ sdk.Msg = &MsgMsgInvitationNotCosignedClaim{}

// MsgMsgInvitationNotCosignedClaim - struct for unjailing jailed validator
type MsgMsgInvitationNotCosignedClaim struct {
	ValidatorAddress     sdk.ValAddress `json:"validator_address" yaml: "validator_address"`
	MainchainAddress     string         `json:"mainchain_address" yaml:"mainchain_address"`
	FirstCosignerAddress sdk.ValAddress `json:"first_cosigner_address" yaml:"first_cosigner_address"`
}

// NewMsgMsgInvitationNotCosignedClaim creates a new MsgMsgInvitationNotCosignedClaim instance
func NewMsgMsgInvitationNotCosignedClaim(validatorAddr sdk.ValAddress) MsgMsgInvitationNotCosignedClaim {
	return MsgMsgInvitationNotCosignedClaim{
		ValidatorAddr: validatorAddr,
	}
}

const invitationNotCosignedConst = "invitationNotCosigned"

// nolint
func (msg MsgMsgInvitationNotCosignedClaim) Route() string { return RouterKey }
func (msg MsgMsgInvitationNotCosignedClaim) Type() string  { return invitationNotCosignedConst }
func (msg MsgMsgInvitationNotCosignedClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddr)}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgMsgInvitationNotCosignedClaim) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgMsgInvitationNotCosignedClaim) ValidateBasic() error {
	if msg.ValidatorAddr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing validator address")
	}
	return nil
}

// TODO: Describe your actions, these will implment the interface of `sdk.Msg`
/*
// verify interface at compile time
var _ sdk.Msg = &Msg<Action>{}

// Msg<Action> - struct for unjailing jailed validator
type Msg<Action> struct {
	ValidatorAddr sdk.ValAddress `json:"address" yaml:"address"` // address of the validator operator
}

// NewMsg<Action> creates a new Msg<Action> instance
func NewMsg<Action>(validatorAddr sdk.ValAddress) Msg<Action> {
	return Msg<Action>{
		ValidatorAddr: validatorAddr,
	}
}

const <action>Const = "<action>"

// nolint
func (msg Msg<Action>) Route() string { return RouterKey }
func (msg Msg<Action>) Type() string  { return <action>Const }
func (msg Msg<Action>) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddr)}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg Msg<Action>) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg Msg<Action>) ValidateBasic() error {
	if msg.ValidatorAddr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing validator address"
	}
	return nil
}
*/
