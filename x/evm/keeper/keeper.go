package keeper

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"slices"
	"sort"

	"cosmossdk.io/store/prefix"
	abci "github.com/cometbft/cometbft/abci/types"
	tmtypes "github.com/cometbft/cometbft/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	"github.com/eni-chain/go-eni/utils"
	"github.com/eni-chain/go-eni/x/evm/blocktest"
	"github.com/eni-chain/go-eni/x/evm/querier"
	"github.com/eni-chain/go-eni/x/evm/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/misc/eip4844"
	"github.com/ethereum/go-ethereum/core"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"

	//enidbtypes "github.com/eni-chain/eni-db/ss/types"

	"sync"

	"github.com/ethereum/go-ethereum/tests"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eni-chain/go-eni/x/evm/types"
)

type (
	AddressNoncePair struct {
		Address common.Address
		Nonce   uint64
	}

	PendingTx struct {
		Key      tmtypes.TxKey
		Nonce    uint64
		Priority int64
	}

	Keeper struct {
		cdc codec.BinaryCodec
		//storeService store.KVStoreService
		logger log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string

		storeKey          storetypes.StoreKey
		transientStoreKey storetypes.StoreKey

		Paramstore paramtypes.Subspace

		txResults []*abci.ExecTxResult
		msgs      []*types.MsgEVMTransaction

		bankKeeper     bankkeeper.Keeper
		accountKeeper  types.AccountKeeper
		stakingKeeper  *stakingkeeper.Keeper
		transferKeeper ibctransferkeeper.Keeper
		//wasmKeeper     *wasmkeeper.PermissionedKeeper
		//wasmViewKeeper *wasmkeeper.Keeper

		cachedFeeCollectorAddressMtx *sync.RWMutex
		cachedFeeCollectorAddress    *common.Address
		nonceMx                      *sync.RWMutex
		pendingTxs                   map[string][]*PendingTx
		keyToNonce                   map[tmtypes.TxKey]*AddressNoncePair

		QueryConfig *querier.Config

		// only used during blocktest. Not used in chain critical path.
		EthBlockTestConfig blocktest.Config
		BlockTest          *tests.BlockTest

		// used for both ETH replay and block tests. Not used in chain critical path.
		Trie        ethstate.Trie
		DB          ethstate.Database
		Root        common.Hash
		ReplayBlock *ethtypes.Block

		//receiptStore enidbtypes.StateStore
	}
)

func NewKeeper(
	storeKey storetypes.StoreKey,
	transientStoreKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	//storeService store.KVStoreService,
	logger log.Logger,
	authority string,

	accountKeeper types.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	stakingKeeper *stakingkeeper.Keeper,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		storeKey:          storeKey,
		transientStoreKey: transientStoreKey,
		cdc:               cdc,
		//storeService:      storeService,
		authority: authority,
		logger:    logger,

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
	}
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) AccountKeeper() types.AccountKeeper {
	return k.accountKeeper
}

func (k *Keeper) BankKeeper() bankkeeper.Keeper {
	return k.bankKeeper
}

//func (k *Keeper) WasmKeeper() *wasmkeeper.PermissionedKeeper {
//	return k.wasmKeeper
//}

func (k *Keeper) GetStoreKey() storetypes.StoreKey {
	return k.storeKey
}

func (k *Keeper) IterateAll(ctx sdk.Context, pref []byte, cb func(key, val []byte) bool) {
	iter := k.PrefixStore(ctx, pref).Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		if cb(iter.Key(), iter.Value()) {
			break
		}
	}
}

func (k *Keeper) PrefixStore(ctx sdk.Context, pref []byte) storetypes.KVStore {
	store := ctx.KVStore(k.GetStoreKey())
	return prefix.NewStore(store, pref)
}

func (k *Keeper) PurgePrefix(ctx sdk.Context, pref []byte) {
	store := k.PrefixStore(ctx, pref)
	////if err := store.DeleteAll(nil, nil); err != nil {
	//	panic(err)
	//}
	store.Delete(nil)
}

func (k *Keeper) GetVMBlockContext(ctx sdk.Context, gp core.GasPool) (*vm.BlockContext, error) {
	if k.EthBlockTestConfig.Enabled {
		return k.getBlockTestBlockCtx(ctx)
	}
	coinbase, err := k.GetFeeCollectorAddress(ctx)
	if err != nil {
		return nil, err
	}

	// Use hash of block timestamp as info for PREVRANDAO
	r, err := ctx.BlockHeader().Time.MarshalBinary()
	if err != nil {
		return nil, err
	}
	rh := crypto.Keccak256Hash(r)

	txfer := func(db vm.StateDB, sender, recipient common.Address, amount *uint256.Int) {
		if IsPayablePrecompile(&recipient) {
			state.TransferWithoutEvents(db, sender, recipient, amount)
		} else {
			core.Transfer(db, sender, recipient, amount)
		}
	}

	return &vm.BlockContext{
		CanTransfer: core.CanTransfer,
		Transfer:    txfer,
		GetHash:     k.GetHashFn(ctx),
		Coinbase:    coinbase,
		GasLimit:    gp.Gas(),
		BlockNumber: big.NewInt(ctx.BlockHeight()),
		Time:        uint64(ctx.BlockHeader().Time.Unix()),
		Difficulty:  utils.Big0, // only needed for PoW
		//BaseFee:     k.GetCurrBaseFeePerGas(ctx).TruncateInt().BigInt(),
		//BaseFee:     k.GetCurrBaseFeePerGas(ctx),
		BlobBaseFee: utils.Big1, // Cancun not enabled
		Random:      &rh,
	}, nil
}

// returns a function that provides block header hash based on block number
func (k *Keeper) GetHashFn(ctx sdk.Context) vm.GetHashFunc {
	return func(height uint64) common.Hash {
		if height > math.MaxInt64 {
			ctx.Logger().Error("Eni block height is bounded by int64 range")
			return common.Hash{}
		}
		h := int64(height)
		if ctx.BlockHeight() == h {
			// current header hash is in the context already
			return common.BytesToHash(ctx.HeaderHash())
		}
		if ctx.BlockHeight() < h {
			// future block doesn't have a hash yet
			return common.Hash{}
		}
		// fetch historical hash from historical info
		return k.getHistoricalHash(ctx, h)
	}
}

func (k *Keeper) getHistoricalHash(ctx sdk.Context, h int64) common.Hash {
	//histInfo, found := k.stakingKeeper.GetHistoricalInfo(ctx, h)
	histInfo, err := k.stakingKeeper.GetHistoricalInfo(ctx, h)
	//if !found {
	if err != nil {
		// too old, already pruned
		return common.Hash{}
	}
	header, _ := tmtypes.HeaderFromProto(&histInfo.Header)

	return common.BytesToHash(header.Hash())
}

// CalculateNextNonce calculates the next nonce for an address
// If includePending is true, it will consider pending nonces
// If includePending is false, it will only return the next nonce from GetNonce
func (k *Keeper) CalculateNextNonce(ctx sdk.Context, addr common.Address, includePending bool) uint64 {
	k.nonceMx.Lock()
	defer k.nonceMx.Unlock()

	nextNonce := k.GetNonce(ctx, addr)

	// we only want the latest nonce if we're not including pending
	if !includePending {
		return nextNonce
	}

	// get the pending nonces (nil is fine)
	pending := k.pendingTxs[addr.Hex()]

	// Check each nonce starting from latest until we find a gap
	// That gap is the next nonce we should use.
	for ; ; nextNonce++ {
		// if it's not in pending, then it's the next nonce
		if _, found := sort.Find(len(pending), func(i int) int { return uint64Cmp(nextNonce, pending[i].Nonce) }); !found {
			return nextNonce
		}
	}
}

// AddPendingNonce adds a pending nonce to the keeper
func (k *Keeper) AddPendingNonce(key tmtypes.TxKey, addr common.Address, nonce uint64, priority int64) {
	k.nonceMx.Lock()
	defer k.nonceMx.Unlock()

	addrStr := addr.Hex()
	if existing, ok := k.keyToNonce[key]; ok {
		if existing.Nonce != nonce {
			fmt.Printf("Seeing transactions with the same hash %X but different nonces (%d vs. %d), which should be impossible\n", key, nonce, existing.Nonce)
		}
		if existing.Address != addr {
			fmt.Printf("Seeing transactions with the same hash %X but different addresses (%s vs. %s), which should be impossible\n", key, addr.Hex(), existing.Address.Hex())
		}
		// we want to no-op whether it's a genuine duplicate or not
		return
	}
	for _, pendingTx := range k.pendingTxs[addrStr] {
		if pendingTx.Nonce == nonce {
			if priority > pendingTx.Priority {
				// replace existing tx
				delete(k.keyToNonce, pendingTx.Key)
				pendingTx.Priority = priority
				pendingTx.Key = key
				k.keyToNonce[key] = &AddressNoncePair{
					Address: addr,
					Nonce:   nonce,
				}
			}
			// we don't need to return error here if priority is lower.
			// Tendermint will take care of rejecting the tx from mempool
			return
		}
	}
	k.keyToNonce[key] = &AddressNoncePair{
		Address: addr,
		Nonce:   nonce,
	}
	k.pendingTxs[addrStr] = append(k.pendingTxs[addrStr], &PendingTx{
		Key:      key,
		Nonce:    nonce,
		Priority: priority,
	})
	slices.SortStableFunc(k.pendingTxs[addrStr], func(a, b *PendingTx) int {
		if a.Nonce < b.Nonce {
			return -1
		} else if a.Nonce > b.Nonce {
			return 1
		}
		return 0
	})
}

// RemovePendingNonce removes a pending nonce from the keeper but leaves a hole
// so that a future transaction must use this nonce.
func (k *Keeper) RemovePendingNonce(key tmtypes.TxKey) {
	k.nonceMx.Lock()
	defer k.nonceMx.Unlock()
	tx, ok := k.keyToNonce[key]

	if !ok {
		return
	}

	delete(k.keyToNonce, key)

	addr := tx.Address.Hex()
	pendings := k.pendingTxs[addr]
	firstMatch, found := sort.Find(len(pendings), func(i int) int { return uint64Cmp(tx.Nonce, pendings[i].Nonce) })
	if !found {
		fmt.Printf("Removing tx %X without a corresponding pending nonce, which should not happen\n", key)
		return
	}
	k.pendingTxs[addr] = append(k.pendingTxs[addr][:firstMatch], k.pendingTxs[addr][firstMatch+1:]...)
	if len(k.pendingTxs[addr]) == 0 {
		delete(k.pendingTxs, addr)
	}
}

func (k *Keeper) SetTxResults(txResults []*abci.ExecTxResult) {
	k.txResults = txResults
}

func (k *Keeper) SetMsgs(msgs []*types.MsgEVMTransaction) {
	k.msgs = msgs
}

// Test use only
func (k *Keeper) GetPendingTxs() map[string][]*PendingTx {
	return k.pendingTxs
}

// Test use only
func (k *Keeper) GetKeysToNonces() map[tmtypes.TxKey]*AddressNoncePair {
	return k.keyToNonce
}

func (k *Keeper) GetBaseFee(ctx sdk.Context) *big.Int {

	if k.EthBlockTestConfig.Enabled {
		bb := k.BlockTest.Json.Blocks[ctx.BlockHeight()-1]
		b, err := bb.Decode()
		if err != nil {
			panic(err)
		}
		return b.Header().BaseFee
	}
	return nil
}

func (k *Keeper) GetReplayedHeight(ctx sdk.Context) int64 {
	return k.getInt64State(ctx, types.ReplayedHeight)
}

func (k *Keeper) SetReplayedHeight(ctx sdk.Context) {
	k.setInt64State(ctx, types.ReplayedHeight, ctx.BlockHeight())
}

func (k *Keeper) GetReplayInitialHeight(ctx sdk.Context) int64 {
	return k.getInt64State(ctx, types.ReplayInitialHeight)
}

func (k *Keeper) SetReplayInitialHeight(ctx sdk.Context, h int64) {
	k.setInt64State(ctx, types.ReplayInitialHeight, h)
}

func (k *Keeper) setInt64State(ctx sdk.Context, key []byte, val int64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(val))
	store.Set(key, bz)
}

func (k *Keeper) getInt64State(ctx sdk.Context, key []byte) int64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(key)
	if bz == nil {
		return 0
	}
	return int64(binary.BigEndian.Uint64(bz))
}

func (k *Keeper) getBlockTestBlockCtx(ctx sdk.Context) (*vm.BlockContext, error) {
	bb := k.BlockTest.Json.Blocks[ctx.BlockHeight()-1]
	b, err := bb.Decode()
	if err != nil {
		return nil, err
	}
	header := b.Header()
	getHash := func(height uint64) common.Hash {
		height = height + 1
		for i := 0; i < len(k.BlockTest.Json.Blocks); i++ {
			if k.BlockTest.Json.Blocks[i].BlockHeader.Number.Uint64() == height {
				return k.BlockTest.Json.Blocks[i].BlockHeader.Hash
			}
		}
		panic(fmt.Sprintf("block hash not found for height %d", height))
	}
	var (
		baseFee     *big.Int
		blobBaseFee *big.Int
		random      *common.Hash
	)
	if header.BaseFee != nil {
		baseFee = new(big.Int).Set(header.BaseFee)
	}
	if header.ExcessBlobGas != nil {
		blobBaseFee = eip4844.CalcBlobFeeOld(*header.ExcessBlobGas)
	} else {
		blobBaseFee = eip4844.CalcBlobFeeOld(0)
	}
	if header.Difficulty.Cmp(common.Big0) == 0 {
		random = &header.MixDigest
	}
	return &vm.BlockContext{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash:     getHash,
		Coinbase:    header.Coinbase,
		GasLimit:    header.GasLimit,
		BlockNumber: new(big.Int).Set(header.Number),
		Time:        header.Time,
		Difficulty:  new(big.Int).Set(header.Difficulty),
		BaseFee:     baseFee,
		BlobBaseFee: blobBaseFee,
		Random:      random,
	}, nil
}

func uint64Cmp(a, b uint64) int {
	if a < b {
		return -1
	} else if a == b {
		return 0
	}
	return 1
}
