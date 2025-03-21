package evmrpc

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/evm/keeper"
	"time"
	//"github.com/tendermint/tendermint/libs/time"
	//rpcclient "github.com/tendermint/tendermint/rpc/client"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
)

type NetAPI struct {
	tmClient       rpcclient.Client
	keeper         *keeper.Keeper
	ctxProvider    func(int64) sdk.Context
	txDecoder      sdk.TxDecoder
	connectionType ConnectionType
}

func NewNetAPI(tmClient rpcclient.Client, k *keeper.Keeper, ctxProvider func(int64) sdk.Context, txDecoder sdk.TxDecoder, connectionType ConnectionType) *NetAPI {
	return &NetAPI{tmClient: tmClient, keeper: k, ctxProvider: ctxProvider, txDecoder: txDecoder, connectionType: connectionType}
}

func (i *NetAPI) Version() string {
	//startTime := time.Now()
	startTime := time.Now().Round(0).UTC()

	defer recordMetrics("net_version", i.connectionType, startTime, true)
	return fmt.Sprintf("%d", i.keeper.ChainID(i.ctxProvider(LatestCtxHeight)).Uint64())
}
