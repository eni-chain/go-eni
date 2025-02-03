package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eni-chain/go-eni/x/evm/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestEVMAddressToEniAddressKey(t *testing.T) {
	evmAddr := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	expectedPrefix := types.EVMAddressToEniAddressKeyPrefix
	key := types.EVMAddressToEniAddressKey(evmAddr)

	require.Equal(t, expectedPrefix[0], key[0], "Key prefix for evm address to eni address key is incorrect")
	require.Equal(t, append(expectedPrefix, evmAddr.Bytes()...), key, "Generated key format is incorrect")
}

func TestEniAddressToEVMAddressKey(t *testing.T) {
	eniAddr := sdk.AccAddress("eni1234567890abcdef1234567890abcdef12345678")
	expectedPrefix := types.EniAddressToEVMAddressKeyPrefix
	key := types.EniAddressToEVMAddressKey(eniAddr)

	require.Equal(t, expectedPrefix[0], key[0], "Key prefix for eni address to evm address key is incorrect")
	require.Equal(t, append(expectedPrefix, eniAddr...), key, "Generated key format is incorrect")
}

func TestStateKey(t *testing.T) {
	evmAddr := common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	expectedPrefix := types.StateKeyPrefix
	key := types.StateKey(evmAddr)

	require.Equal(t, expectedPrefix[0], key[0], "Key prefix for state key is incorrect")
	require.Equal(t, append(expectedPrefix, evmAddr.Bytes()...), key, "Generated key format is incorrect")
}

func TestBlockBloomKey(t *testing.T) {
	height := int64(123456)
	key := types.BlockBloomKey(height)

	require.Equal(t, types.BlockBloomPrefix[0], key[0], "Key prefix for block bloom key is incorrect")
}
