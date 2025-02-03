package state_test

import (
	"math/big"
	"testing"

	"github.com/eni-chain/go-eni/x/evm/state"
	"github.com/stretchr/testify/require"
)

func TestGetCoinbaseAddress(t *testing.T) {
	coinbaseAddr := state.GetCoinbaseAddress(1).String()
	require.Equal(t, coinbaseAddr, "eni1v4mx6hmrda5kucnpwdjsqqqqqqqqqqqpz6djs7")
}

func TestSplitUeniWeiAmount(t *testing.T) {
	for _, test := range []struct {
		amt         *big.Int
		expectedEni *big.Int
		expectedWei *big.Int
	}{
		{
			amt:         big.NewInt(0),
			expectedEni: big.NewInt(0),
			expectedWei: big.NewInt(0),
		}, {
			amt:         big.NewInt(1),
			expectedEni: big.NewInt(0),
			expectedWei: big.NewInt(1),
		}, {
			amt:         big.NewInt(999_999_999_999),
			expectedEni: big.NewInt(0),
			expectedWei: big.NewInt(999_999_999_999),
		}, {
			amt:         big.NewInt(1_000_000_000_000),
			expectedEni: big.NewInt(1),
			expectedWei: big.NewInt(0),
		}, {
			amt:         big.NewInt(1_000_000_000_001),
			expectedEni: big.NewInt(1),
			expectedWei: big.NewInt(1),
		}, {
			amt:         big.NewInt(123_456_789_123_456_789),
			expectedEni: big.NewInt(123456),
			expectedWei: big.NewInt(789_123_456_789),
		},
	} {
		ueni, wei := state.SplitUeniWeiAmount(test.amt)
		require.Equal(t, test.expectedEni, ueni.BigInt())
		require.Equal(t, test.expectedWei, wei.BigInt())
	}
}
