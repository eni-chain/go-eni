package epoch_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"

	"github.com/eni-chain/go-eni/app"
	"github.com/eni-chain/go-eni/x/epoch"
	"github.com/eni-chain/go-eni/x/epoch/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestNewHandler(t *testing.T) {
	app := app.Setup(false, false) // Your setup function here
	handler := epoch.NewHandler(app.EpochKeeper)

	// Test unrecognized message type
	testMsg := testdata.NewTestMsg()
	_, err := handler(app.BaseApp.NewContext(false, tmproto.Header{}), testMsg)
	require.Error(t, err)

	expectedErrMsg := fmt.Sprintf("unrecognized %s message type", types.ModuleName)
	require.ErrorContains(t, err, expectedErrMsg)
}
