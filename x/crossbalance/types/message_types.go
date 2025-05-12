package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgCrossBalanceTransfer 定义了跨账户转账的消息结构
type MsgCrossBalanceTransfer struct {
	FromAddress string    `protobuf:"bytes,1,opt,name=from_address,json=fromAddress,proto3" json:"from_address,omitempty"`
	ToAddress   string    `protobuf:"bytes,2,opt,name=to_address,json=toAddress,proto3" json:"to_address,omitempty"`
	Amount      sdk.Coins `protobuf:"bytes,3,rep,name=amount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"amount"`
}

// MsgCrossBalanceTransferResponse 定义了跨账户转账的响应结构
type MsgCrossBalanceTransferResponse struct{}
