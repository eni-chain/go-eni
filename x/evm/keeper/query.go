package keeper

import (
	"github.com/eni-chain/go-eni/x/evm/types"
)

var _ types.QueryServer = Keeper{}
