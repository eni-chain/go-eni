package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/eni-chain/go-eni/testutil/keeper"
	"github.com/eni-chain/go-eni/x/epoch/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.EpochKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
