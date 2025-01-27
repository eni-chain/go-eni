package keeper

import (
	"github.com/eni-chain/go-eni/x/epoch/types"
)

var _ types.QueryServer = Keeper{}
