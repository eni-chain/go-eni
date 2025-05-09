package evmrpc_test

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetBlockByHash(t *testing.T) {
	resObj := sendRequestGood(t, "getBlockByHash", "0x0000000000000000000000000000000000000000000000000000000000000001", true)
	verifyBlockResult(t, resObj)
}

func TestGetEniBlockByHash(t *testing.T) {
	resObj := sendEniRequestGood(t, "getBlockByHash", "0x0000000000000000000000000000000000000000000000000000000000000001", true)
	verifyBlockResult(t, resObj)
}

func TestGetEniBlockByNumberExcludeTraceFail(t *testing.T) {
	resObj := sendEniRequestGood(t, "getBlockByNumberExcludeTraceFail", "0x67", true)
	// first tx is not a panic tx, second tx is a panic tx
	expectedNumTxs := 1
	require.Equal(t, expectedNumTxs, len(resObj["result"].(map[string]interface{})["transactions"].([]interface{})))
}

func TestGetBlockByNumber(t *testing.T) {
	resObjEarliest := sendEniRequestGood(t, "getBlockByNumber", "earliest", true)
	verifyGenesisBlockResult(t, resObjEarliest)
	for _, num := range []string{"0x8", "latest", "pending", "finalized", "safe"} {
		resObj := sendRequestGood(t, "getBlockByNumber", num, true)
		verifyBlockResult(t, resObj)
	}

	resObj := sendRequestBad(t, "getBlockByNumber", "bad_num", true)
	require.Equal(t, "invalid argument 0: hex string without 0x prefix", resObj["error"].(map[string]interface{})["message"])
}

func TestGetEniBlockByNumber(t *testing.T) {
	resObjEarliest := sendEniRequestGood(t, "getBlockByNumber", "earliest", true)
	verifyGenesisBlockResult(t, resObjEarliest)
	for _, num := range []string{"0x8", "latest", "pending", "finalized", "safe"} {
		resObj := sendEniRequestGood(t, "getBlockByNumber", num, true)
		verifyBlockResult(t, resObj)
	}

	resObj := sendEniRequestBad(t, "getBlockByNumber", "bad_num", true)
	require.Equal(t, "invalid argument 0: hex string without 0x prefix", resObj["error"].(map[string]interface{})["message"])
}

func TestGetBlockTransactionCount(t *testing.T) {
	// get by block number
	for _, num := range []string{"0x8", "earliest", "latest", "pending", "finalized", "safe"} {
		resObj := sendRequestGood(t, "getBlockTransactionCountByNumber", num)
		require.Equal(t, "0x1", resObj["result"])
	}

	// get error returns null
	for _, num := range []string{"0x8", "earliest", "latest", "pending", "finalized", "safe", "0x0000000000000000000000000000000000000000000000000000000000000001"} {
		resObj := sendRequestBad(t, "getBlockTransactionCountByNumber", num)
		require.Nil(t, resObj["result"])
	}

	// get by hash
	resObj := sendRequestGood(t, "getBlockTransactionCountByHash", "0x0000000000000000000000000000000000000000000000000000000000000001")
	require.Equal(t, "0x1", resObj["result"])
}

func TestGetBlockReceipts(t *testing.T) {
	// Query by block height
	resObj := sendRequestGood(t, "getBlockReceipts", "0x2")
	result := resObj["result"].([]interface{})
	require.Equal(t, 3, len(result))
	receipt1 := result[0].(map[string]interface{})
	require.Equal(t, "0x2", receipt1["blockNumber"])
	require.Equal(t, multiTxBlockTx1.Hash().Hex(), receipt1["transactionHash"])
	receipt2 := result[1].(map[string]interface{})
	require.Equal(t, "0x2", receipt2["blockNumber"])
	require.Equal(t, multiTxBlockTx2.Hash().Hex(), receipt2["transactionHash"])
	receipt3 := result[2].(map[string]interface{})
	require.Equal(t, "0x2", receipt3["blockNumber"])
	require.Equal(t, multiTxBlockTx3.Hash().Hex(), receipt3["transactionHash"])

	resObjEni := sendEniRequestGood(t, "getBlockReceipts", "0x2")
	result = resObjEni["result"].([]interface{})
	require.Equal(t, 5, len(result))

	// Query by block hash
	resObj2 := sendRequestGood(t, "getBlockReceipts", MultiTxBlockHash)
	result = resObj2["result"].([]interface{})
	require.Equal(t, 3, len(result))
	receipt1 = result[0].(map[string]interface{})
	require.Equal(t, "0x2", receipt1["blockNumber"])
	require.Equal(t, multiTxBlockTx1.Hash().Hex(), receipt1["transactionHash"])
	receipt2 = result[1].(map[string]interface{})
	require.Equal(t, "0x2", receipt2["blockNumber"])
	require.Equal(t, multiTxBlockTx2.Hash().Hex(), receipt2["transactionHash"])
	receipt3 = result[2].(map[string]interface{})
	require.Equal(t, "0x2", receipt3["blockNumber"])
	require.Equal(t, multiTxBlockTx3.Hash().Hex(), receipt3["transactionHash"])

	// Query by tag latest => retrieves block 8
	resObj3 := sendRequestGood(t, "getBlockReceipts", "latest")
	result = resObj3["result"].([]interface{})
	require.Equal(t, 1, len(result))
	receipt1 = result[0].(map[string]interface{})
	require.Equal(t, "0x8", receipt1["blockNumber"])
	require.Equal(t, tx1.Hash().Hex(), receipt1["transactionHash"])
}

func verifyGenesisBlockResult(t *testing.T, resObj map[string]interface{}) {
	resObj = resObj["result"].(map[string]interface{})
	require.Equal(t, "0x0", resObj["baseFeePerGas"])
	require.Equal(t, "0x0", resObj["difficulty"])
	require.Equal(t, "0x", resObj["extraData"])
	require.Equal(t, "0x0", resObj["gasLimit"])
	require.Equal(t, "0x0", resObj["gasUsed"])
	require.Equal(t, "0xF9D3845DF25B43B1C6926F3CEDA6845C17F5624E12212FD8847D0BA01DA1AB9E", resObj["hash"])
	require.Equal(t, "0x0000000000000000", resObj["nonce"])
	require.Equal(t, "0x0", resObj["number"])
}

func verifyBlockResult(t *testing.T, resObj map[string]interface{}) {
	resObj = resObj["result"].(map[string]interface{})
	require.Equal(t, "0x0", resObj["difficulty"])
	require.Equal(t, "0x", resObj["extraData"])
	require.Equal(t, "0xbebc200", resObj["gasLimit"])
	require.Equal(t, "0x5", resObj["gasUsed"])
	require.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000001", resObj["hash"])
	// see setup_tests.go, which have one transaction for block 0x8 (latest)
	require.Equal(t, "0x00002000040000000000000000000080000000200000000000002000000000080000000000000000000000000000000000000000000000000800000000000000001000000000000000000000000000000000020000000000000000000000000100000000000000002000000000200000000000000000000000000000000000100000000000000000000000000400000000000000200000000000000000000000000000000000000100000000000000020000200000000000000000002000000000000000000000000000000000000000000000000000000000000000000200000000010000000002000000000000000000000000000000010200000000000000", resObj["logsBloom"])
	require.Equal(t, "0x0000000000000000000000000000000000000005", resObj["miner"])
	require.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000000", resObj["mixHash"])
	require.Equal(t, "0x0000000000000000", resObj["nonce"])
	require.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000006", resObj["parentHash"])
	require.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000004", resObj["receiptsRoot"])
	require.Equal(t, "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347", resObj["sha3Uncles"])
	require.Equal(t, "0x26b", resObj["size"])
	require.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000003", resObj["stateRoot"])
	require.Equal(t, "0x65254651", resObj["timestamp"])
	tx := resObj["transactions"].([]interface{})[0].(map[string]interface{})
	require.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000001", tx["blockHash"])
	require.Equal(t, "0x5b4eba929f3811980f5ae0c5d04fa200f837df4e", tx["from"])
	require.Equal(t, "0x3e8", tx["gas"])
	require.Equal(t, "0xa", tx["gasPrice"])
	require.Equal(t, "0xa", tx["maxFeePerGas"])
	require.Equal(t, "0x0", tx["maxPriorityFeePerGas"])
	require.Equal(t, "0xf02362077ac075a397344172496b28e913ce5294879d811bb0269b3be20a872e", tx["hash"])
	require.Equal(t, "0x616263", tx["input"])
	require.Equal(t, "0x1", tx["nonce"])
	require.Equal(t, "0x0000000000000000000000000000000000010203", tx["to"])
	require.Equal(t, "0x0", tx["transactionIndex"])
	require.Equal(t, "0x3e8", tx["value"])
	require.Equal(t, "0x2", tx["type"])
	require.Equal(t, []interface{}{}, tx["accessList"])
	require.Equal(t, "0xae3f3", tx["chainId"])
	require.Equal(t, "0x0", tx["v"])
	require.Equal(t, "0xa1ac0e5b8202742e54ae7af350ed855313cc4f9861c2d75a0e541b4aff7c981e", tx["r"])
	require.Equal(t, "0x288b16881aed9640cd360403b9db1ce3961b29af0b00158311856d1446670996", tx["s"])
	require.Equal(t, "0x0", tx["yParity"])
	require.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000002", resObj["transactionsRoot"])
	require.Equal(t, []interface{}{}, resObj["uncles"])
	require.Equal(t, "0x3b9aca00", resObj["baseFeePerGas"])
	require.Equal(t, "0x0", resObj["totalDifficulty"])
}

func TestEncodeTmBlock_EmptyTransactions(t *testing.T) {
	//k := &testkeeper.EVMTestApp.EvmKeeper
	//ctx := testkeeper.EVMTestApp.GetContextForDeliverTx([]byte{}).WithBlockTime(time.Now())
	//block := &coretypes.ResultBlock{
	//	BlockID: MockBlockID,
	//	Block: &tmtypes.Block{
	//		Header: mockBlockHeader(MockHeight8),
	//		Data:   tmtypes.Data{},
	//		LastCommit: &tmtypes.Commit{
	//			Height: MockHeight8 - 1,
	//		},
	//	},
	//}
	//blockRes := &coretypes.ResultBlockResults{
	//	TxsResults: []*abci.ExecTxResult{},
	//	ConsensusParamUpdates: &types2.ConsensusParams{
	//		Block: &types2.BlockParams{
	//			MaxBytes: 100000000,
	//			MaxGas:   200000000,
	//		},
	//	},
	//}

	//// Call EncodeTmBlock with empty transactions
	//result, err := evmrpc.EncodeTmBlock(ctx, block, blockRes, ethtypes.Bloom{}, k, Decoder, true, false, nil)
	//require.Nil(t, err)

	//// Assert txHash is equal to ethtypes.EmptyTxsHash
	//require.Equal(t, ethtypes.EmptyTxsHash, result["transactionsRoot"])
}

func TestEncodeBankMsg(t *testing.T) {
	//k := &testkeeper.EVMTestApp.EvmKeeper
	//ctx := testkeeper.EVMTestApp.GetContextForDeliverTx([]byte{}).WithBlockTime(time.Now())
	//fromEniAddr, _ := testkeeper.MockAddressPair()
	//toEniAddr, _ := testkeeper.MockAddressPair()
	//b := TxConfig.NewTxBuilder()
	//b.SetMsgs(banktypes.NewMsgSend(fromEniAddr, toEniAddr, sdk.NewCoins(sdk.NewCoin("ueni", sdk.NewInt(10)))))
	//tx := b.GetTx()
	//resBlock := coretypes.ResultBlock{
	//	BlockID: MockBlockID,
	//	Block: &tmtypes.Block{
	//		Header: mockBlockHeader(MockHeight8),
	//		Data: tmtypes.Data{
	//			Txs: []tmtypes.Tx{func() []byte {
	//				bz, _ := Encoder(tx)
	//				return bz
	//			}()},
	//		},
	//		LastCommit: &tmtypes.Commit{
	//			Height: MockHeight8 - 1,
	//		},
	//	},
	//}
	//resBlockRes := coretypes.ResultBlockResults{
	//	TxsResults: []*abci.ExecTxResult{
	//		{
	//			Data: func() []byte {
	//				bz, _ := Encoder(tx)
	//				return bz
	//			}(),
	//		},
	//	},
	//	ConsensusParamUpdates: &types2.ConsensusParams{
	//		Block: &types2.BlockParams{
	//			MaxBytes: 100000000,
	//			MaxGas:   200000000,
	//		},
	//	},
	//}
	//res, err := evmrpc.EncodeTmBlock(ctx, &resBlock, &resBlockRes, ethtypes.Bloom{}, k, Decoder, true, false, nil)
	//require.Nil(t, err)
	//txs := res["transactions"].([]interface{})
	//require.Equal(t, 0, len(txs))
}
