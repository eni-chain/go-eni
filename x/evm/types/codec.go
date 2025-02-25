package types

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
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
