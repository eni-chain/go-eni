package app_test

import (
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server/api"
	cosmosConfig "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	abci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/eni-chain/go-eni/app"
	"github.com/stretchr/testify/require"
)

func TestEmptyBlockIdempotency(t *testing.T) {
	commitData := [][]byte{}
	tm := time.Now().UTC()
	valPub := secp256k1.GenPrivKey().PubKey()

	for i := 1; i <= 10; i++ {
		testWrapper := app.NewTestWrapper(t, tm, valPub, false)
		res, _ := testWrapper.App.FinalizeBlock(&abci.RequestFinalizeBlock{Height: 1})
		testWrapper.App.Commit()
		data := res.AppHash
		commitData = append(commitData, data)
	}

	referenceData := commitData[0]
	for _, data := range commitData[1:] {
		require.Equal(t, len(referenceData), len(data))
	}
}

func TestInvalidProposalWithExcessiveGasWanted(t *testing.T) {
	tm := time.Now().UTC()
	valPub := secp256k1.GenPrivKey().PubKey()

	testWrapper := app.NewTestWrapper(t, tm, valPub, false)
	ap := testWrapper.App
	testWrapper.Ctx = testWrapper.Ctx.WithConsensusParams(cmtproto.ConsensusParams{
		Block: &cmtproto.BlockParams{MaxGas: 10},
	})
	emptyTxBuilder := app.MakeEncodingConfig().TxConfig.NewTxBuilder()
	txEncoder := app.MakeEncodingConfig().TxConfig.TxEncoder()
	emptyTxBuilder.SetGasLimit(10)
	emptyTx, _ := txEncoder(emptyTxBuilder.GetTx())

	badProposal := abci.RequestProcessProposal{
		Txs:    [][]byte{emptyTx, emptyTx},
		Height: 1,
	}

	res, err := ap.ProcessProposal(&badProposal)
	require.Nil(t, err)
	require.Equal(t, abci.ResponseProcessProposal_REJECT, res.Status)
}

func TestOverflowGas(t *testing.T) {
	tm := time.Now().UTC()
	valPub := secp256k1.GenPrivKey().PubKey()

	testWrapper := app.NewTestWrapper(t, tm, valPub, false)
	ap := testWrapper.App
	testWrapper.Ctx = testWrapper.Ctx.WithConsensusParams(cmtproto.ConsensusParams{
		Block: &cmtproto.BlockParams{MaxGas: math.MaxInt64},
	})
	emptyTxBuilder := app.MakeEncodingConfig().TxConfig.NewTxBuilder()
	txEncoder := app.MakeEncodingConfig().TxConfig.TxEncoder()
	emptyTxBuilder.SetGasLimit(uint64(math.MaxInt64))
	emptyTx, _ := txEncoder(emptyTxBuilder.GetTx())

	secondEmptyTxBuilder := app.MakeEncodingConfig().TxConfig.NewTxBuilder()
	secondEmptyTxBuilder.SetGasLimit(10)
	secondTx, _ := txEncoder(secondEmptyTxBuilder.GetTx())

	proposal := abci.RequestProcessProposal{
		Txs:    [][]byte{emptyTx, secondTx},
		Height: 1,
	}
	res, err := ap.ProcessProposal(&proposal)
	require.Nil(t, err)
	require.Equal(t, abci.ResponseProcessProposal_REJECT, res.Status)
}

func TestApp_RegisterAPIRoutes(t *testing.T) {
	type args struct {
		apiSvr    *api.Server
		apiConfig cosmosConfig.APIConfig
	}
	tests := []struct {
		name        string
		args        args
		wantSwagger bool
	}{
		{
			name: "swagger added to the router if configured",
			args: args{
				apiSvr: &api.Server{
					ClientCtx:         client.Context{},
					Router:            &mux.Router{},
					GRPCGatewayRouter: runtime.NewServeMux(),
				},
				apiConfig: cosmosConfig.APIConfig{
					Swagger: true,
				},
			},
			wantSwagger: true,
		},
		{
			name: "swagger not added to the router if not configured",
			args: args{
				apiSvr: &api.Server{
					ClientCtx:         client.Context{},
					Router:            &mux.Router{},
					GRPCGatewayRouter: runtime.NewServeMux(),
				},
				apiConfig: cosmosConfig.APIConfig{},
			},
			wantSwagger: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eniApp := &app.App{}
			eniApp.RegisterAPIRoutes(tt.args.apiSvr, tt.args.apiConfig)
			routes := tt.args.apiSvr.Router
			gotSwagger := isSwaggerRouteAdded(routes)

			if !reflect.DeepEqual(gotSwagger, tt.wantSwagger) {
				t.Errorf("Run() gotSwagger = %v, want %v", gotSwagger, tt.wantSwagger)
			}
		})

	}
}

func TestGetDeliverTxEntry(t *testing.T) {
	tm := time.Now().UTC()
	valPub := secp256k1.GenPrivKey().PubKey()

	testWrapper := app.NewTestWrapper(t, tm, valPub, false)
	ap := testWrapper.App
	ctx := testWrapper.Ctx.WithConsensusParams(cmtproto.ConsensusParams{
		Block: &cmtproto.BlockParams{MaxGas: 10},
	})
	emptyTxBuilder := app.MakeEncodingConfig().TxConfig.NewTxBuilder()
	txEncoder := app.MakeEncodingConfig().TxConfig.TxEncoder()
	emptyTxBuilder.SetGasLimit(10)
	tx := emptyTxBuilder.GetTx()
	bz, _ := txEncoder(tx)

	require.NotNil(t, ap.GetDeliverTxEntry(ctx, 0, bz, tx))

	require.NotNil(t, ap.GetDeliverTxEntry(ctx, 0, bz, nil))
}

func isSwaggerRouteAdded(router *mux.Router) bool {
	var isAdded bool
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil && pathTemplate == "/swagger/" {
			isAdded = true
		}
		return nil
	})
	if err != nil {
		return false
	}
	return isAdded
}
