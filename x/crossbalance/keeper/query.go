package keeper

import (
	"github.com/eni-chain/go-eni/x/crossbalance/types"
)

var _ types.QueryServer = Keeper{}
