package evmrpc

import (
	"errors"

	"github.com/eni-chain/go-eni/x/evm/artifacts/cw20"
)

type TestAPI struct{}

func NewTestAPI() *TestAPI {
	return &TestAPI{}
}

func (a *TestAPI) IncrementPointerVersion(pointerType string, offset int16) error {
	switch pointerType {
	case "cw20":
		cw20.SetVersionWithOffset(offset)
	default:
		return errors.New("invalid pointer type")
	}
	return nil
}
