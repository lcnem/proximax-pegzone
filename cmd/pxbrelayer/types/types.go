package types

type Event byte

const (
	Unsupported Event = iota
	MsgPegClaim
	MsgUnpegNotCosignedClaim
	MsgInvitationNotCosignedClaim
)

func (d Event) String() string {
	return [...]string{"unsupported", "peg_claim", "unpeg_not_cosigned_claim", "invitation_not_cosigned_claim"}[d]
}

type CosmosMsgAttributeKey int

const (
	UnsupportedAttributeKey CosmosMsgAttributeKey = iota
	CosmosSender
	MainchainTxHash
	ToAddress
	Coin
)

// String returns the event type as a string
func (d CosmosMsgAttributeKey) String() string {
	return [...]string{"unsupported", "cosmos_sender", "mainchain_tx_hash", "to_address", "amount"}[d]
}
