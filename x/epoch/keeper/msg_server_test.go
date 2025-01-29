package keeper_test

import (
	"testing"

	keepertest "github.com/eni-chain/go-eni/testutil/keeper"
	"github.com/eni-chain/go-eni/x/epoch/keeper"
	"github.com/stretchr/testify/require"
)

func TestSetupMsgServer(t *testing.T) {
	k, _ := keepertest.EpochKeeper(t)
	msg := keeper.NewMsgServerImpl(*k)
	require.NotNil(t, msg)
}
