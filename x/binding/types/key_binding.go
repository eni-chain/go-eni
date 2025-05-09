package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// BindingKeyPrefix is the prefix to retrieve all Binding
	BindingKeyPrefix = "Binding/value/"
)

// BindingKey returns the store key to retrieve a Binding from the index fields
func BindingKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
