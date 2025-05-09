package keeper

import (
	"github.com/eni-chain/go-eni/x/binding/types"
)

var _ types.QueryServer = Keeper{}
