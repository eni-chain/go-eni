package artifacts

import (
	"fmt"

	"github.com/eni-chain/go-eni/x/evm/artifacts/native"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

func GetParsedABI(typ string) *abi.ABI {
	switch typ {
	case "native":
		return native.GetParsedABI()

	default:
		panic(fmt.Sprintf("unknown artifact type %s", typ))
	}
}

func GetBin(typ string) []byte {
	switch typ {
	case "native":
		return native.GetBin()
	default:
		panic(fmt.Sprintf("unknown artifact type %s", typ))
	}
}
