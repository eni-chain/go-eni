package app

import (
	"fmt"
	"sort"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	evmante "github.com/cosmos/cosmos-sdk/x/evm/ante"
	evmtypes "github.com/cosmos/cosmos-sdk/x/evm/types"
	"github.com/cosmos/cosmos-sdk/x/evm/types/ethtx"
)

type TxsGrouper interface {
	FilterEvmTxs(typedTx sdk.Tx, encodedTx []byte) (interface{}, error)
	DecodeEvmTxs(msgData interface{}, idx int, rawTx []byte) (*TxMeta, error)
	SortGlobalNonce() error
	GroupByAddressTxs() error
	GroupSequentialTxs() error
}

type TxGroup struct {
	evmTxs       [][]byte
	otherTxs     [][]byte
	evmTxMetas   []*TxMeta
	txDecoder    sdk.TxDecoder
	groups       map[string][]*TxMeta
	serialGroups [][]byte
	otherGroups  [][]byte
}

type TxMeta struct {
	RawTx []byte
	To    string
	Nonce uint64
	Data  []byte
	idx   int
}

func NewTxGroup(txDecoder sdk.TxDecoder) *TxGroup {
	return &TxGroup{
		evmTxs:       make([][]byte, 0),
		otherTxs:     make([][]byte, 0),
		evmTxMetas:   make([]*TxMeta, 0),
		txDecoder:    txDecoder,
		groups:       make(map[string][]*TxMeta),
		serialGroups: make([][]byte, 0),
		otherGroups:  make([][]byte, 0),
	}
}

func (app *App) GroupByTxs(ctx sdk.Context, txs [][]byte) (*TxGroup, error) {
	typedTxs := make([]sdk.Tx, len(txs))
	//txGroup := NewTxGroup(app.txDecoder)
	txGroup := NewTxGroup(app.TxDecode)
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	for i, tx := range txs {
		wg.Add(1)
		go func(idx int, encodedTx []byte) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					ctx.Logger().Error(fmt.Sprintf("encountered panic during transaction decoding: %s", err))
					mu.Lock()
					typedTxs[idx] = nil
					mu.Unlock()
				}
			}()

			// Decode transaction
			//typedTx, err := app.txDecoder(encodedTx)
			typedTx, err := app.TxDecode(encodedTx)
			if err != nil {
				ctx.Logger().Error(fmt.Sprintf("error decoding transaction at index %d due to %s", idx, err))
				mu.Lock()
				typedTxs[idx] = nil
				mu.Unlock()
				return
			}

			// check is evm transaction
			if isEVM, _ := evmante.IsEVMMessage(typedTx); isEVM {
				// get cached value from transaction
				msgData, err := txGroup.FilterEvmTxs(typedTx, encodedTx)
				if err != nil {
					ctx.Logger().Error(fmt.Sprintf("error getting cached value from transaction at index %d", idx))
					mu.Lock()
					typedTxs[idx] = nil
					mu.Unlock()
					return
				}

				// get tx meta and nonce from transaction
				txMeta, err := txGroup.DecodeEvmTxs(msgData, idx, encodedTx)
				if err != nil {
					ctx.Logger().Error(fmt.Sprintf("error getting tx meta and nonce from transaction at index %d due to %s", idx, err))
					mu.Lock()
					typedTxs[idx] = nil
					mu.Unlock()
					return
				}

				mu.Lock()
				txGroup.evmTxMetas = append(txGroup.evmTxMetas, txMeta)
				mu.Unlock()
			} else {
				txGroup.otherTxs = append(txGroup.otherTxs, encodedTx)
			}
			typedTxs[idx] = typedTx
		}(i, tx)
	}
	wg.Wait()

	// sort evm transactions by global nonce
	if err := txGroup.SortGlobalNonce(); err != nil {
		return nil, err
	}

	// group by address
	if err := txGroup.GroupByAddressTxs(); err != nil {
		return nil, err
	}

	// group sequential transactions
	if err := txGroup.GroupSequentialTxs(); err != nil {
		return nil, err
	}
	return txGroup, nil
}

// FilterEvmTxs filter evm transactions and cache value
func (t *TxGroup) FilterEvmTxs(typedTx sdk.Tx, encodedTx []byte) (interface{}, error) {
	t.evmTxs = append(t.evmTxs, encodedTx)
	msg := evmtypes.MustGetEVMTransactionMessage(typedTx)
	cachedValue := msg.Data.GetCachedValue()
	// todo: not deal with preprocessing EVM  tx for now
	if cachedValue == nil {
		return nil, fmt.Errorf("error getting cached value")
	}
	return cachedValue, nil
}

// DecodeEvmTxs decode evm transactions and get tx meta and nonce
func (t *TxGroup) DecodeEvmTxs(msgData interface{}, idx int, rawTx []byte) (*TxMeta, error) {
	txMeta := &TxMeta{
		RawTx: rawTx,
		idx:   idx,
	}
	switch tx := msgData.(type) {
	case *ethtx.DynamicFeeTx:
		txMeta.To = tx.GetTo().Hex()
		txMeta.Nonce = tx.GetNonce()
		txMeta.Data = tx.GetData()
	case *ethtx.AccessListTx:
		txMeta.To = tx.GetTo().Hex()
		txMeta.Nonce = tx.GetNonce()
		txMeta.Data = tx.GetData()

	case *ethtx.BlobTx:
		txMeta.To = tx.GetTo().Hex()
		txMeta.Nonce = tx.GetNonce()
		txMeta.Data = tx.GetData()

	case *ethtx.LegacyTx:
		txMeta.To = tx.GetTo().Hex()
		txMeta.Nonce = tx.GetNonce()
		txMeta.Data = tx.GetData()

	default:
		return nil, fmt.Errorf("unsupported transaction type: %T", tx)
	}

	return txMeta, nil
}

// SortGlobalNonce sort evm transactions by global nonce
func (t *TxGroup) SortGlobalNonce() error {
	sort.Slice(t.evmTxMetas, func(i, j int) bool {
		return t.evmTxMetas[i].Nonce < t.evmTxMetas[j].Nonce
	})
	return nil
}

// GroupByAddressTxs group evm transactions by address and sort by nonce
func (t *TxGroup) GroupByAddressTxs() error {
	groups := make(map[string][]*TxMeta)
	for _, meta := range t.evmTxMetas {
		groups[meta.To] = append(groups[meta.To], meta)
	}

	// Sort by nonce in each group
	for to := range groups {
		sort.Slice(groups[to], func(i, j int) bool {
			return groups[to][i].Nonce < groups[to][j].Nonce
		})
	}
	t.groups = groups
	return nil
}

// GroupSequentialTxs group sequential evm transactions
func (t *TxGroup) GroupSequentialTxs() error {
	// Store grouped transactions into serialGroups (assuming that each group requires serial processing)
	t.serialGroups = make([][]byte, 0)
	for _, group := range t.groups {
		var groupTxs [][]byte
		// if group has only one transaction, it can be processed in parallel
		if len(group) == 1 {
			t.otherGroups = append(t.otherGroups, group[0].RawTx)
		}
		// group sequential transactions
		if len(group) > 1 {
			for _, meta := range group {
				groupTxs = append(groupTxs, meta.RawTx)
			}
			t.serialGroups = append(t.serialGroups, groupTxs...)
		}
	}

	// other tx can be processed in parallel
	t.otherGroups = append(t.otherGroups, t.otherTxs...)
	return nil
}
