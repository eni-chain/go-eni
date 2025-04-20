package main

import (
	"bufio"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"os"
	"sync/atomic"
	"testing"
)

func TestUnmarshal(t *testing.T) {

}

func BenchmarkName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		filePath := "./scripts-test-1000/newtx_batch1000_11.txt"
		file, err := os.Open(filePath)
		if err != nil {
			panic(err.Error())
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {

			messageType := "offlineTx"

			var signedTx SignedTx
			// Sign EVM and Cosmos TX differently
			line := scanner.Text()
			if len(line) == 0 {
				break
			}
			tx := &ethtypes.Transaction{}
			tx.UnmarshalBinary(common.FromHex(line))
			signedTx = SignedTx{EvmTx: tx, MsgType: messageType}
			EvmTxHashes = append(EvmTxHashes, signedTx.EvmTx.Hash())
			atomic.AddInt64(ReadTxCount, 1)
		}
	}
}
