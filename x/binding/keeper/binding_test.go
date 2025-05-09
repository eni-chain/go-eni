package keeper_test

import (
	"context"
	"strconv"
	"testing"

	keepertest "github.com/eni-chain/go-eni/testutil/keeper"
	"github.com/eni-chain/go-eni/testutil/nullify"
	"github.com/eni-chain/go-eni/x/binding/keeper"
	"github.com/eni-chain/go-eni/x/binding/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNBinding(keeper keeper.Keeper, ctx context.Context, n int) []types.Binding {
	items := make([]types.Binding, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)

		keeper.SetBinding(ctx, items[i])
	}
	return items
}

func TestBindingGet(t *testing.T) {
	keeper, ctx := keepertest.BindingKeeper(t)
	items := createNBinding(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetBinding(ctx,
			item.Index,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestBindingRemove(t *testing.T) {
	keeper, ctx := keepertest.BindingKeeper(t)
	items := createNBinding(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveBinding(ctx,
			item.Index,
		)
		_, found := keeper.GetBinding(ctx,
			item.Index,
		)
		require.False(t, found)
	}
}

func TestBindingGetAll(t *testing.T) {
	keeper, ctx := keepertest.BindingKeeper(t)
	items := createNBinding(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllBinding(ctx)),
	)
}
