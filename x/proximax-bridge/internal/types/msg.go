package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// debug
var _ sdk.Msg = &MsgPeg{}

// MsgUnpeg - struct for unjailing jailed validator
type MsgPeg struct {
	Address         sdk.AccAddress `json:"address" yaml:"address"`
	MainchainTxHash string         `json:"mainchain_tx_hash" yaml:"mainchain_tx_hash"`
	Amount          sdk.Coins      `json:"amount" yaml:"amount"`
}

// NewMsgUnpeg creates a new MsgUnpeg instance
func NewMsgPeg(address sdk.AccAddress, mainchainTxHash string, amount sdk.Coins) MsgPeg {
	return MsgPeg{
		Address:         address,
		MainchainTxHash: mainchainTxHash,
		Amount:          amount,
	}
}

const pegConst = "peg"

// nolint
func (msg MsgPeg) Route() string { return RouterKey }
func (msg MsgPeg) Type() string  { return unpegConst }
func (msg MsgPeg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Address)}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgPeg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgPeg) ValidateBasic() error {
	if msg.Address.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing validator address")
	}
	return nil
}

//

// TODO: Describe your actions, these will implment the interface of `sdk.Msg`
// verify interface at compile time
var _ sdk.Msg = &MsgPegClaim{}

// MsgPegClaim - struct for unjailing jailed validator
type MsgPegClaim struct {
	Address          sdk.AccAddress `json:"address" yaml:"address"`
	MainchainTxHash  string         `json:"mainchain_tx_hash" yaml:"mainchain_tx_hash"`
	Amount           sdk.Coins      `json:"amount" yaml:"amount"`
	ValidatorAddress sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
}

// NewMsgPegClaim creates a new MsgPegClaim instance
func NewMsgPegClaim(address sdk.AccAddress, mainchainTxHash string, amount sdk.Coins, validatorAddress sdk.ValAddress) MsgPegClaim {
	return MsgPegClaim{
		Address:          address,
		MainchainTxHash:  mainchainTxHash,
		Amount:           amount,
		ValidatorAddress: validatorAddress,
	}
}

const pegClaimConst = "peg_claim"

// nolint
func (msg MsgPegClaim) Route() string { return RouterKey }
func (msg MsgPegClaim) Type() string  { return pegClaimConst }
func (msg MsgPegClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddress)}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgPegClaim) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgPegClaim) ValidateBasic() error {
	if msg.Address.Empty() {
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
	Amount               sdk.Coins      `json:"amount" yaml:"amount"`
	FirstCosignerAddress sdk.ValAddress `json:"first_cosigner_address" yaml:"first_cosigner_address"`
}

// NewMsgUnpeg creates a new MsgUnpeg instance
func NewMsgUnpeg(address sdk.AccAddress, mainchainAddress string, amount sdk.Coins, firstCosignerAddress sdk.ValAddress) MsgUnpeg {
	return MsgUnpeg{
		Address:              address,
		MainchainAddress:     mainchainAddress,
		Amount:               amount,
		FirstCosignerAddress: firstCosignerAddress,
	}
}

const unpegConst = "unpeg"

// nolint
func (msg MsgUnpeg) Route() string { return RouterKey }
func (msg MsgUnpeg) Type() string  { return unpegConst }
func (msg MsgUnpeg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgUnpeg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgUnpeg) ValidateBasic() error {
	if msg.Address.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing validator address")
	}
	return nil
}

var _ sdk.Msg = &MsgRecordUnpeg{}

// MsgUnpeg - struct for unjailing jailed validator
type MsgRecordUnpeg struct {
	Address                sdk.AccAddress `json:"address" yaml:"address"`
	MainchainTxHash        string         `json:"mainchain_tx_hash" yaml:"mainchain_tx_hash"`
	Amount                 sdk.Coins      `json:"amount" yaml:"amount"`
	FirstCosignerPublicKey string         `json:"first_cosigner_public_key" yaml:"first_cosigner_public_key"`
	ValidatorAddress       sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
}

// NewMsgUnpeg creates a new MsgUnpeg instance
func NewMsgRecordUnpeg(address sdk.AccAddress, mainchainTxHash string, amount sdk.Coins, firstCosignerPublicKey string, validatorAddress sdk.ValAddress) MsgRecordUnpeg {
	return MsgRecordUnpeg{
		Address:                address,
		MainchainTxHash:        mainchainTxHash,
		Amount:                 amount,
		FirstCosignerPublicKey: firstCosignerPublicKey,
		ValidatorAddress:       validatorAddress,
	}
}

const recordUnpegConst = "record_unpeg"

// nolint
func (msg MsgRecordUnpeg) Route() string { return RouterKey }
func (msg MsgRecordUnpeg) Type() string  { return recordUnpegConst }
func (msg MsgRecordUnpeg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddress)}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgRecordUnpeg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgRecordUnpeg) ValidateBasic() error {
	if msg.Address.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing validator address")
	}
	return nil
}

var _ sdk.Msg = &MsgNotifyCosigned{}

// MsgUnpeg - struct for unjailing jailed validator
type MsgNotifyCosigned struct {
	Address           sdk.ValAddress `json:"address" yaml:"address"`
	MainchainTxHash   string         `json:"mainchain_tx_hash" yaml:"mainchain_tx_hash"`
	CosignerPublicKey string         `json:"cosigner_public_key" yaml:"cosigner_public_key"`
}

// NewMsgUnpeg creates a new MsgUnpeg instance
func NewMsgNotifyCosigned(address sdk.ValAddress, mainchainTxHash, cosignerPublicKey string) MsgNotifyCosigned {
	return MsgNotifyCosigned{
		Address:           address,
		MainchainTxHash:   mainchainTxHash,
		CosignerPublicKey: cosignerPublicKey,
	}
}

const notifyCosignedConst = "notify_cosigned"

// nolint
func (msg MsgNotifyCosigned) Route() string { return RouterKey }
func (msg MsgNotifyCosigned) Type() string  { return notifyCosignedConst }
func (msg MsgNotifyCosigned) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Address)}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgNotifyCosigned) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgNotifyCosigned) ValidateBasic() error {
	if msg.Address.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing validator address")
	}
	return nil
}

// TODO: Describe your actions, these will implment the interface of `sdk.Msg`
// verify interface at compile time
var _ sdk.Msg = &MsgNotCosignedClaim{}

// MsgNotCosignedClaim - struct for unjailing jailed validator
type MsgNotCosignedClaim struct {
	Address sdk.ValAddress `json:"address" yaml:"address"`
	TxHash  string         `json:"tx_hash" yaml:"tx_hash"`
}

// NewMsgNotCosignedClaim creates a new MsgNotCosignedClaim instance
func NewMsgNotCosignedClaim(address sdk.ValAddress, txHash string) MsgNotCosignedClaim {
	return MsgNotCosignedClaim{
		Address: address,
		TxHash:  txHash,
	}
}

const notCosignedClaimConst = "not_cosigned_claim"

// nolint
func (msg MsgNotCosignedClaim) Route() string { return RouterKey }
func (msg MsgNotCosignedClaim) Type() string  { return notCosignedClaimConst }
func (msg MsgNotCosignedClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Address)}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgNotCosignedClaim) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgNotCosignedClaim) ValidateBasic() error {
	if msg.Address.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing validator address")
	}
	return nil
}

// TODO: Describe your actions, these will implment the interface of `sdk.Msg`

// verify interface at compile time
var _ sdk.Msg = &MsgRequestInvitation{}

// MsgRequestInvitation - struct for unjailing jailed validator
type MsgRequestInvitation struct {
	Address              sdk.ValAddress `json:"address" yaml:"address"`
	NewCosignerPublicKey string         `json:"new_cosigner_public_key" yaml:"new_cosigner_public_key"`
	FirstCosignerAddress sdk.ValAddress `json:"first_cosigner_address" yaml:"first_cosigner_address"`
}

// NewMsgRequestInvitation creates a new MsgRequestInvitation instance
func NewMsgRequestInvitation(address sdk.ValAddress, newCosignerPublicKey string, firstCosignerAddress sdk.ValAddress) MsgRequestInvitation {
	return MsgRequestInvitation{
		Address:              address,
		NewCosignerPublicKey: newCosignerPublicKey,
		FirstCosignerAddress: firstCosignerAddress,
	}
}

const requestInvitationConst = "request_invitation"

// nolint
func (msg MsgRequestInvitation) Route() string { return RouterKey }
func (msg MsgRequestInvitation) Type() string  { return requestInvitationConst }
func (msg MsgRequestInvitation) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Address)}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgRequestInvitation) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgRequestInvitation) ValidateBasic() error {
	if msg.Address.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing validator address")
	}
	return nil
}

var _ sdk.Msg = &MsgPendingRequestInvitation{}

// MsgRequestInvitation - struct for unjailing jailed validator
type MsgPendingRequestInvitation struct {
	Address                sdk.ValAddress `json:"address" yaml:"address"`
	NewCosignerPublicKey   string         `json:"new_cosigner_public_key" yaml:"new_cosigner_public_key"`
	FirstCosignerAddress   sdk.ValAddress `json:"first_cosigner_address" yaml:"first_cosigner_address"`
	FirstCosignerPublicKey string         `json:"first_cosigner_public_key" yaml:"first_cosigner_public_key"`
	TxHash                 string         `json:"tx_hash" yaml:"tx_hash"`
}

// NewMsgRequestInvitation creates a new MsgRequestInvitation instance
func NewMsgPendingRequestInvitation(address sdk.ValAddress, newCosignerPublicKey string, firstCosignerAddress sdk.ValAddress, firstCosignerPublicKey, txHash string) MsgPendingRequestInvitation {
	return MsgPendingRequestInvitation{
		Address:                address,
		NewCosignerPublicKey:   newCosignerPublicKey,
		FirstCosignerAddress:   firstCosignerAddress,
		FirstCosignerPublicKey: firstCosignerPublicKey,
		TxHash:                 txHash,
	}
}

const pendingRequestInvitationConst = "pending_request_invitation"

// nolint
func (msg MsgPendingRequestInvitation) Route() string { return RouterKey }
func (msg MsgPendingRequestInvitation) Type() string  { return pendingRequestInvitationConst }
func (msg MsgPendingRequestInvitation) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.FirstCosignerAddress)}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgPendingRequestInvitation) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgPendingRequestInvitation) ValidateBasic() error {
	if msg.Address.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing validator address")
	}
	return nil
}

var _ sdk.Msg = &MsgNewCosignerInvited{}

// MsgProximaXTransactionStatus - struct for unjailing jailed validator
type MsgNewCosignerInvited struct {
	Address            sdk.ValAddress `json:"address" yaml:"address"`
	TxHash             string         `json:"mainchain_tx_hash" yaml:"mainchain_tx_hash"`
	MainchainPublicKey string         `json:"mainchain_public_key" yaml:"mainchain_public_key"`
}

// NewMsgRequestInvitation creates a new MsgRequestInvitation instance
func NewMsgNewCosignerInvited(address sdk.ValAddress, txHash, pubKey string) MsgNewCosignerInvited {
	return MsgNewCosignerInvited{
		Address:            address,
		TxHash:             txHash,
		MainchainPublicKey: pubKey,
	}
}

const newCosignerInvitedConst = "new_cosigner_invited"

// nolint
func (msg MsgNewCosignerInvited) Route() string { return RouterKey }
func (msg MsgNewCosignerInvited) Type() string  { return newCosignerInvitedConst }
func (msg MsgNewCosignerInvited) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Address)}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgNewCosignerInvited) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgNewCosignerInvited) ValidateBasic() error {
	if msg.Address.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing validator address")
	}
	return nil
}

//hoge
var _ sdk.Msg = &MsgConfirmedInvitation{}

// MsgProximaXTransactionStatus - struct for unjailing jailed validator
type MsgConfirmedInvitation struct {
	Address sdk.ValAddress `json:"address" yaml:"address"`
	TxHash  string         `json:"mainchain_tx_hash" yaml:"mainchain_tx_hash"`
}

// NewMsgRequestInvitation creates a new MsgRequestInvitation instance
func NewMsgConfirmedInvitation(address sdk.ValAddress, txHash string) MsgConfirmedInvitation {
	return MsgConfirmedInvitation{
		Address: address,
		TxHash:  txHash,
	}
}

const confirmedInvitationConst = "confirmed_invitation"

// nolint
func (msg MsgConfirmedInvitation) Route() string { return RouterKey }
func (msg MsgConfirmedInvitation) Type() string  { return confirmedInvitationConst }
func (msg MsgConfirmedInvitation) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Address)}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgConfirmedInvitation) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgConfirmedInvitation) ValidateBasic() error {
	if msg.Address.Empty() {
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
