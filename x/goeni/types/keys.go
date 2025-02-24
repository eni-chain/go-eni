package types

const (
	// ModuleName defines the module name
	ModuleName = "goeni"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_goeni"
)

var (
	ParamsKey = []byte("p_goeni")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
