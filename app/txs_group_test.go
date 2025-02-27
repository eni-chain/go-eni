/**
 * Created by Adwind.
 * User: liuyunlong
 * Date: 2025/2/16
 * Time: 13:37
 */
package app

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestTxGroup_SortGlobalNonce(t1 *testing.T) {
	type fields struct {
		txs            [][]byte
		evmTxs         [][]byte
		evmTxMetas     []*TxMeta
		txDecoder      types.TxDecoder
		serialGroups   [][]byte
		parallelGroups [][]byte
		wg             sync.WaitGroup
		mu             sync.Mutex
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "TestTxGroup_SortGlobalNonce",
			fields: fields{
				evmTxMetas: []*TxMeta{
					{RawTx: []byte("tx3"), Nonce: 3, To: "to3"},
					{RawTx: []byte("tx4"), Nonce: 4, To: "to4"},
					{RawTx: []byte("tx5"), Nonce: 5, To: "to5"},
					{RawTx: []byte("tx1"), Nonce: 1, To: "to1"},
					{RawTx: []byte("tx2"), Nonce: 2, To: "to2"},
					{RawTx: []byte("tx6"), Nonce: 6, To: "to6"},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TxGroup{
				evmTxMetas: tt.fields.evmTxMetas,
			}
			fmt.Println("before sort")
			for key, value := range t.evmTxMetas {
				fmt.Println(key, value.Nonce, value.To, value.RawTx)
			}
			_ = t.SortGlobalNonce()
			fmt.Println("after sort")
			for key, value := range t.evmTxMetas {
				fmt.Println(key, value.Nonce, value.To, value.RawTx)
			}
		})
	}
}
