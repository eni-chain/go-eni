package types

const (
	// ModuleName defines the module name
	ModuleName = "binding"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_binding"
)

var (
	ParamsKey = []byte("p_binding")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
