package migrations_test

import (
	"testing"
	"time"

	testkeeper "github.com/eni-chain/go-eni/testutil/keeper"
	"github.com/eni-chain/go-eni/utils"
	"github.com/eni-chain/go-eni/x/evm/migrations"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/stretchr/testify/require"
)

func TestMigrateERCNativePointers(t *testing.T) {
	t.Skip("This test is not applicable to the current implementation")
	k := testkeeper.EVMTestApp.EvmKeeper
	ctx := testkeeper.EVMTestApp.GetContextForDeliverTx([]byte{}).WithBlockTime(time.Now())
	var pointerAddr common.Address
	require.Nil(t, k.RunWithOneOffEVMInstance(ctx, func(e *vm.EVM) error {
		a, err := k.UpsertERCNativePointer(ctx, e, "test", utils.ERCMetadata{Name: "name", Symbol: "symbol", Decimals: 6})
		pointerAddr = a
		return err
	}, func(s1, s2 string) {}))
	require.Nil(t, migrations.MigrateERCNativePointers(ctx, &k))
	// address should stay the same
	addr, _, _ := k.GetERC20NativePointer(ctx, "test")
	require.Equal(t, pointerAddr, addr)
}
