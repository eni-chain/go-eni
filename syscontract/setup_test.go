package syscontract

import (
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"testing"
)

func TestSetupSystemContracts(t *testing.T) {
	type args struct {
		blockNumber *big.Int
		statedb     vm.StateDB
		logger      log.Logger
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestSetupSystemContracts",
			args: args{
				blockNumber: big.NewInt(0),
				statedb:     nil,
				logger:      log.New(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetupSystemContracts(tt.args.blockNumber, tt.args.statedb, tt.args.logger)
		})
	}
}
