package keeper_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/eni-chain/go-eni/precompiles/bank"
	"github.com/eni-chain/go-eni/precompiles/gov"
	"github.com/eni-chain/go-eni/precompiles/staking"
	"github.com/eni-chain/go-eni/testutil/keeper"
	evmkeeper "github.com/eni-chain/go-eni/x/evm/keeper"
)

func toAddr(addr string) *common.Address {
	ca := common.HexToAddress(addr)
	return &ca
}

func TestIsPayablePrecompile(t *testing.T) {
	_, evmAddr := keeper.MockAddressPair()
	require.False(t, evmkeeper.IsPayablePrecompile(&evmAddr))
	require.False(t, evmkeeper.IsPayablePrecompile(nil))

	require.True(t, evmkeeper.IsPayablePrecompile(toAddr(bank.BankAddress)))
	require.True(t, evmkeeper.IsPayablePrecompile(toAddr(staking.StakingAddress)))
	require.True(t, evmkeeper.IsPayablePrecompile(toAddr(gov.GovAddress)))
}
