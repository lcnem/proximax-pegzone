package types

const (
	// ModuleName is the name of the module
	ModuleName = "proximaxbridge"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	StoreKeyForPeg = ModuleName + "_peg"

	StoreKeyForUnpeg = ModuleName + "_unpeg"

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	QuerierRoute = ModuleName
)
