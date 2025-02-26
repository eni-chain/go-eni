package keeper

import (
	"github.com/ethereum/go-ethereum/common"
)

// add any payable precompiles here
// these will suppress transfer events to/from the precompile address
var payablePrecompiles = map[string]struct{}{
	//bank.BankAddress:       {},
	//staking.StakingAddress: {},
	//gov.GovAddress:         {},
}

func IsPayablePrecompile(addr *common.Address) bool {
	if addr == nil {
		return false
	}
	_, ok := payablePrecompiles[addr.Hex()]
	return ok
}
