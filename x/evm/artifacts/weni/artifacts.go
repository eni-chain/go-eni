package weni

import (
	"embed"
	"encoding/hex"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const CurrentVersion uint16 = 1

//go:embed WENI.abi
//go:embed WENI.bin
var f embed.FS

var cachedBin []byte
var cachedABI *abi.ABI

func GetABI() []byte {
	bz, err := f.ReadFile("WENI.abi")
	if err != nil {
		panic("failed to read WENI contract ABI")
	}
	return bz
}

func GetParsedABI() *abi.ABI {
	if cachedABI != nil {
		return cachedABI
	}
	parsedABI, err := abi.JSON(strings.NewReader(string(GetABI())))
	if err != nil {
		panic(err)
	}
	cachedABI = &parsedABI
	return cachedABI
}

func GetBin() []byte {
	if cachedBin != nil {
		return cachedBin
	}
	code, err := f.ReadFile("WENI.bin")
	if err != nil {
		panic("failed to read WENI contract binary")
	}
	bz, err := hex.DecodeString(string(code))
	if err != nil {
		panic("failed to decode WENI contract binary")
	}
	cachedBin = bz
	return bz
}
