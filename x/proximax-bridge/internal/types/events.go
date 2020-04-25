package types

// proximax-bridge module event types
const (
	EventTypeCreateClaim    = "create_claim"
	EventTypeProphecyStatus = "prophecy_status"
	EventTypeUnpeg          = "unpeg"

	AttributeKeyMainchainTxHash = "mainchain_tx_hash"
	AttributeKeyCosmosReceiver  = "cosmos_receiver"
	AttributeKeyAmount          = "amount"
	AttributeKeyStatus          = "status"
	AttributeKeyClaimType       = "claim_type"

	AttributeKeyMultisigCustodyAddress = "multisig_custody_address"
	AttributeKeyCosmosSender           = "cosmos_sender"
	AttributeKeyMainchainReceiver      = "mainchain_receiver"

	AttributeKeyTxHash                = "tx_hash"
	AttributeKeyNotCosignedValidators = "not_cosigned_validators"

	AttributeValueCategory = ModuleName
)
