package types

import (
	"errors"
	"fmt"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"github.com/eni-chain/go-eni/x/evm/types/ethtx"
	"github.com/gogo/protobuf/proto"
	// this line is used by starport scaffolding # 1
)

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgEVMTransaction{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSend{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRegisterPointer{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAssociateContractAddress{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAssociate{},
	)
	// this line is used by starport scaffolding # 3

	//registry.RegisterImplementations((*sdk.Msg)(nil),
	//	&MsgUpdateParams{},
	//)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

func UnpackTxData(any *cdctypes.Any) (ethtx.TxData, error) {
	if any == nil {
		return nil, errors.New("protobuf Any message cannot be nil")
	}

	txData, ok := any.GetCachedValue().(ethtx.TxData)
	if !ok {
		ltx := ethtx.LegacyTx{}
		if proto.Unmarshal(any.Value, &ltx) == nil {
			// value is a legacy tx
			return &ltx, nil
		}
		atx := ethtx.AccessListTx{}
		if proto.Unmarshal(any.Value, &atx) == nil {
			// value is a accesslist tx
			return &atx, nil
		}
		dtx := ethtx.DynamicFeeTx{}
		if proto.Unmarshal(any.Value, &dtx) == nil {
			// value is a dynamic fee tx
			return &dtx, nil
		}
		btx := ethtx.BlobTx{}
		if proto.Unmarshal(any.Value, &btx) == nil {
			// value is a blob tx
			return &btx, nil
		}
		astx := ethtx.AssociateTx{}
		if proto.Unmarshal(any.Value, &astx) == nil {
			// value is an associate tx
			return &astx, nil
		}
		return nil, fmt.Errorf("cannot unpack Any into TxData %T", any)
	}

	return txData, nil
}
