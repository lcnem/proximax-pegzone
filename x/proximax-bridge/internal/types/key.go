package types

const (
	// ModuleName is the name of the module
	ModuleName = "proximaxbridge"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	StoreKeyForPeg = ModuleName + "_peg"

	StoreKeyForCosign = ModuleName + "_cosign"

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	QuerierRoute = ModuleName
)
