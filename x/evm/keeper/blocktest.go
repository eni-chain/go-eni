package keeper

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (k *Keeper) VerifyAccount(ctx sdk.Context, addr common.Address, accountData types.Account) {
	// we no longer check eth balance due to limiting EVM max refund to 150% of used gas (https://github.com/eni-protocol/go-ethereum/pull/32)
	code := accountData.Code
	for key, expectedState := range accountData.Storage {
		actualState := k.GetState(ctx, addr, key)
		if !bytes.Equal(actualState.Bytes(), expectedState.Bytes()) {
			panic(fmt.Sprintf("storage mismatch for address %s: expected %X, got %X", addr.Hex(), expectedState, actualState))
		}
	}
	nonce := accountData.Nonce
	if !bytes.Equal(code, k.GetCode(ctx, addr)) {
		panic(fmt.Sprintf("code mismatch for address %s", addr))
	}
	if nonce != k.GetNonce(ctx, addr) {
		panic(fmt.Sprintf("nonce mismatch for address %s: expected %d, got %d", addr.Hex(), nonce, k.GetNonce(ctx, addr)))
	}
}
