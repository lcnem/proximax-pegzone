package types

// proximax-bridge module event types
const (
	// TODO: Create your event types
	EventTypePegClaim              = "pegClaim"
	EventTypeUnpeg                 = "unpeg"
	EventTypeUnpegNotCosignedClaim = "unpegNotCosignedClaim"

	// TODO: Create keys fo your events, the values will be derivided from the msg
	// AttributeKeyAddress  		= "address"

	// TODO: Some events may not have values for that reason you want to emit that something happened.
	// AttributeValueDoubleSign = "double_sign"

	AttributeValueCategory = ModuleName
)
