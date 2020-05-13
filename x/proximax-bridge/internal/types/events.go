package types

// proximax-bridge module event types
const (
	EventTypeCreateClaim    = "create_claim"
	EventTypeProphecyStatus = "prophecy_status"
	EventTypePeg            = "peg"
	EventTypeUnpeg          = "unpeg"
	EventTypeInvitation     = "request_invitation"

	AttributeKeyMainchainTxHash = "mainchain_tx_hash"
	AttributeKeyCosmosReceiver  = "cosmos_receiver"
	AttributeKeyAmount          = "amount"
	AttributeKeyStatus          = "status"
	AttributeKeyClaimType       = "claim_type"

	AttributeKeyMultisigCustodyAddress = "multisig_custody_address"
	AttributeKeyMultisigAccountAddress = "multisig_address"
	AttributeKeyCosmosSender           = "cosmos_sender"
	AttributeKeyCosmosAccount          = "cosmos_account"
	AttributeKeyMainchainAddress       = "mainchain_address"
	AttributeKeyNewCosignerPublicKey   = "new_cosigner_public_key"

	AttributeKeyTxHash                = "tx_hash"
	AttributeKeyNotCosignedValidators = "not_cosigned_validators"

	AttributeKeyFirstCosignerAddress = "first_cosigner_address"

	AttributeValueCategory = ModuleName
)
