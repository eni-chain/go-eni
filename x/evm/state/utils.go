package state

import (
	"encoding/binary"
	"math/big"

	cosmossdk_io_math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UeniToSweiMultiplier Fields that were denominated in ueni will be converted to swei (1ueni = 10^12swei)
// for existing Ethereum application (which assumes 18 decimal points) to display properly.
var UeniToSweiMultiplier = big.NewInt(1_000_000_000_000)
var SdkUeniToSweiMultiplier = cosmossdk_io_math.NewIntFromBigInt(UeniToSweiMultiplier)

var CoinbaseAddressPrefix = []byte("evm_coinbase")

func GetCoinbaseAddress(txIdx int) sdk.AccAddress {
	txIndexBz := make([]byte, 8)
	binary.BigEndian.PutUint64(txIndexBz, uint64(txIdx))
	return append(CoinbaseAddressPrefix, txIndexBz...)
}
