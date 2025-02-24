package keeper

import (
	"github.com/eni-chain/go-eni/x/goeni/types"
)

var _ types.QueryServer = Keeper{}
