package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "binding"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_binding"

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName
)

var (
	ParamsKey = []byte("p_binding")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// GetBindingKey 返回绑定关系的存储键
func GetBindingKey(cosmosAddr sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("binding/%s", cosmosAddr.String()))
}
