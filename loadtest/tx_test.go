package main

import (
	"bufio"
	"os"
	"testing"

	evmante "github.com/cosmos/cosmos-sdk/x/evm/ante"
	"github.com/cosmos/cosmos-sdk/x/evm/types"
	"github.com/cosmos/cosmos-sdk/x/evm/types/ethtx"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

func TestConvertEvmTx2EniTx(t *testing.T) {
	//read txt file
	filePath := "./scripts/10000.txt"
	file, err := os.Open(filePath)
	file.w
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//convert line to EniTx
		// Sign EVM and Cosmos TX differently
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		tx := &ethtypes.Transaction{}
		tx.UnmarshalBinary(common.FromHex(line))
		txData, err := ethtx.NewTxDataFromTx(tx)
		if err != nil {
			return
		}
		msg, err := types.NewMsgEVMTransaction(txData)
		if err != nil {
			return
		}
		err = evmante.PreprocessMsgSender(msg)
		if err != nil {
			return
		}
		txBuilder := s.txConfig.NewTxBuilder()
		if err = txBuilder.SetMsgs(msg); err != nil {
			return
		}
		txbz, encodeErr := s.txConfig.TxEncoder()(txBuilder.GetTx())
	}
}
