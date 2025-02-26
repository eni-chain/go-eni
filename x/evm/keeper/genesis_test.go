package keeper_test

//
//import (
//	"bytes"
//	"testing"
//
//	testkeeper "github.com/eni-chain/go-eni/testutil/keeper"
//	"github.com/eni-chain/go-eni/x/evm/keeper"
//	"github.com/stretchr/testify/require"
//)
//
//func TestInitGenesis(t *testing.T) {
//	k := &testkeeper.EVMTestApp.EvmKeeper
//	ctx := testkeeper.EVMTestApp.GetContextForDeliverTx([]byte{})
//	// coinbase address must be associated
//	coinbaseEniAddr, associated := k.GetEniAddress(ctx, keeper.GetCoinbaseAddress())
//	require.True(t, associated)
//	require.True(t, bytes.Equal(coinbaseEniAddr, k.AccountKeeper().GetModuleAddress("fee_collector")))
//}
