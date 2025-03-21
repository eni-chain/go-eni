package app

import (
	"encoding/binary"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/cosmos/cosmos-sdk/x/evm/types"
	"github.com/cosmos/cosmos-sdk/x/evm/types/ethtx"
	genutilstypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/eni-chain/go-eni/utils"
	"github.com/ethereum/go-ethereum/common"
	ethcore "github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	ethtests "github.com/ethereum/go-ethereum/tests"
)

func BlockTest(a *App, bt *ethtests.BlockTest) {
	a.EvmKeeper.BlockTest = bt
	a.EvmKeeper.EthBlockTestConfig.Enabled = true

	gendoc, err := genutilstypes.AppGenesisFromFile(filepath.Join(DefaultNodeHome, "config/genesis.json"))
	if err != nil {
		panic(err)
	}
	_, err = a.InitChain(&abci.RequestInitChain{
		Time:          time.Now(),
		ChainId:       gendoc.ChainID,
		AppStateBytes: gendoc.AppState,
	})
	if err != nil {
		panic(err)
	}

	for addr, genesisAccount := range a.EvmKeeper.BlockTest.Json.Pre {
		a.Logger().Debug("blocktest Pre SetAccount", "address", addr, "account", genesisAccount)
		ueni := math.NewIntFromBigInt(genesisAccount.Balance)
		eniAddr := a.EvmKeeper.GetEniAddressOrDefault(a.GetContextForFinalizeBlock([]byte{}), addr)
		err := a.EvmKeeper.BankKeeper().AddCoins(a.GetContextForFinalizeBlock([]byte{}), eniAddr, sdk.NewCoins(sdk.NewCoin("ueni", ueni)))
		if err != nil {
			panic(err)
		}

		a.EvmKeeper.SetNonce(a.GetContextForFinalizeBlock([]byte{}), addr, genesisAccount.Nonce)
		a.EvmKeeper.SetCode(a.GetContextForFinalizeBlock([]byte{}), addr, genesisAccount.Code)
		for key, value := range genesisAccount.Storage {
			a.EvmKeeper.SetState(a.GetContextForFinalizeBlock([]byte{}), addr, key, value)
		}
		params := a.EvmKeeper.GetParams(a.GetContextForFinalizeBlock([]byte{}))
		params.MinimumFeePerGas = math.LegacyNewDecFromInt(math.NewInt(0))
		a.EvmKeeper.SetParams(a.GetContextForFinalizeBlock([]byte{}), params)
	}

	if len(bt.Json.Blocks) == 0 {
		panic("no blocks found")
	}

	ethblocks := make([]*ethtypes.Block, 0)
	for i, btBlock := range bt.Json.Blocks {
		h := int64(i + 1)
		b, err := btBlock.Decode()
		if err != nil {
			panic(err)
		}
		ethblocks = append(ethblocks, b)
		hash := make([]byte, 8)
		binary.BigEndian.PutUint64(hash, uint64(h))
		a.Logger().Debug("blocktest FinalizeBlock", "height", b.Number().Uint64(), "hash", hash)
		_, err = a.FinalizeBlock(&abci.RequestFinalizeBlock{
			Txs:               utils.Map(b.Transactions(), func(tx *ethtypes.Transaction) []byte { return encodeTx(tx, a.TxConfig()) }),
			ProposerAddress:   a.EvmKeeper.GetEniAddressOrDefault(a.GetContextForCheckTx(nil), b.Coinbase()),
			DecidedLastCommit: abci.CommitInfo{Votes: []abci.VoteInfo{}},
			Height:            h,
			Hash:              hash,
			Time:              time.Now(),
		})
		if err != nil {
			panic(err)
		}
		a.Logger().Debug("blocktest Commit", "height", b.Number().Uint64(), "hash", hash)
		_, err = a.Commit()
		if err != nil {
			panic(err)
		}
	}

	// Check post-state after all blocks are run
	ctx := a.GetContextForCheckTx(nil)
	for addr, accountData := range bt.Json.Post {
		if IsWithdrawalAddress(addr, ethblocks) {
			fmt.Println("Skipping withdrawal address: ", addr)
			continue
		}
		// Not checking compliance with EIP-4788
		if addr == params.BeaconRootsAddress {
			fmt.Println("Skipping beacon roots storage address: ", addr)
			continue
		}
		a.EvmKeeper.VerifyAccount(ctx, addr, accountData)
		fmt.Println("Successfully verified account: ", addr)
	}
}

func encodeTx(tx *ethtypes.Transaction, txConfig client.TxConfig) []byte {
	var txData ethtx.TxData
	var err error
	switch tx.Type() {
	case ethtypes.LegacyTxType:
		txData, err = ethtx.NewLegacyTx(tx)
	case ethtypes.DynamicFeeTxType:
		txData, err = ethtx.NewDynamicFeeTx(tx)
	case ethtypes.AccessListTxType:
		txData, err = ethtx.NewAccessListTx(tx)
	case ethtypes.BlobTxType:
		txData, err = ethtx.NewBlobTx(tx)
	}
	if err != nil {
		if strings.Contains(err.Error(), ethcore.ErrTipAboveFeeCap.Error()) {
			return nil
		}
		panic(err)
	}
	msg, err := evmtypes.NewMsgEVMTransaction(txData)
	signer := ethtypes.LatestSignerForChainID(tx.ChainId())
	sender, _ := ethtypes.Sender(signer, tx)
	eniAddr := sdk.AccAddress(sender[:])
	msg.Sender = eniAddr.String()
	if err != nil {
		panic(err)
	}
	//ante.Preprocess2(msg)
	txBuilder := txConfig.NewTxBuilder()
	if err = txBuilder.SetMsgs(msg); err != nil {
		panic(err)
	}
	txbz, encodeErr := txConfig.TxEncoder()(txBuilder.GetTx())
	if encodeErr != nil {
		panic(encodeErr)
	}
	return txbz
}

func IsWithdrawalAddress(addr common.Address, blocks []*ethtypes.Block) bool {
	for _, block := range blocks {
		for _, w := range block.Withdrawals() {
			if w.Address == addr {
				return true
			}
		}
	}
	return false
}
