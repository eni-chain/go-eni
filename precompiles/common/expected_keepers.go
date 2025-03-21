package common

import (
	"context"
	"cosmossdk.io/math"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"

	//connectiontypes "github.com/cosmos/ibc-go/v3/modules/core/03-connection/types"
	//"github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	//"github.com/cosmos/ibc-go/v3/modules/core/exported"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	//ibctypes "//github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	"github.com/eni-chain/go-eni/utils"
	//oracletypes "github.com/eni-chain/go-eni/x/oracle/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

type BankKeeper interface {
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	//SendCoinsAndWei(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amt math.Int, wei math.Int) error
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	//GetWeiBalance(ctx sdk.Context, addr sdk.AccAddress) math.Int
	GetDenomMetaData(ctx context.Context, denom string) (banktypes.Metadata, bool)
	GetSupply(ctx context.Context, denom string) sdk.Coin
	LockedCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
}

type BankMsgServer interface {
	Send(goCtx context.Context, msg *banktypes.MsgSend) (*banktypes.MsgSendResponse, error)
}

type EVMKeeper interface {
	GetEniAddress(sdk.Context, common.Address) (sdk.AccAddress, bool)
	GetEniAddressOrDefault(ctx sdk.Context, evmAddress common.Address) sdk.AccAddress // only used for getting precompile Eni addresses
	GetEVMAddress(sdk.Context, sdk.AccAddress) (common.Address, bool)
	SetAddressMapping(sdk.Context, sdk.AccAddress, common.Address)
	GetCodeHash(sdk.Context, common.Address) common.Hash
	GetPriorityNormalizer(ctx sdk.Context) math.LegacyDec
	GetBaseDenom(ctx sdk.Context) string
	SetERC20NativePointer(ctx sdk.Context, token string, addr common.Address) error
	GetERC20NativePointer(ctx sdk.Context, token string) (addr common.Address, version uint16, exists bool)
	//SetERC20CW20Pointer(ctx sdk.Context, cw20Address string, addr common.Address) error
	//GetERC20CW20Pointer(ctx sdk.Context, cw20Address string) (addr common.Address, version uint16, exists bool)
	//SetERC721CW721Pointer(ctx sdk.Context, cw721Address string, addr common.Address) error
	//GetERC721CW721Pointer(ctx sdk.Context, cw721Address string) (addr common.Address, version uint16, exists bool)
	//SetERC1155CW1155Pointer(ctx sdk.Context, cw1155Address string, addr common.Address) error
	//GetERC1155CW1155Pointer(ctx sdk.Context, cw1155Address string) (addr common.Address, version uint16, exists bool)
	SetCode(ctx sdk.Context, addr common.Address, code []byte)
	UpsertERCNativePointer(
		ctx sdk.Context, evm *vm.EVM, token string, metadata utils.ERCMetadata,
	) (contractAddr common.Address, err error)
	//UpsertERCCW20Pointer(
	//	ctx sdk.Context, evm *vm.EVM, cw20Addr string, metadata utils.ERCMetadata,
	//) (contractAddr common.Address, err error)
	//UpsertERCCW721Pointer(
	//	ctx sdk.Context, evm *vm.EVM, cw721Addr string, metadata utils.ERCMetadata,
	//) (contractAddr common.Address, err error)
	//UpsertERCCW1155Pointer(
	//	ctx sdk.Context, evm *vm.EVM, cw1155Addr string, metadata utils.ERCMetadata,
	//) (contractAddr common.Address, err error)
	GetEVMGasLimitFromCtx(ctx sdk.Context) uint64
	GetCosmosGasLimitFromEVMGas(ctx sdk.Context, evmGas uint64) uint64
}

type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	HasAccount(ctx context.Context, addr sdk.AccAddress) bool
	SetAccount(ctx context.Context, acc sdk.AccountI)
	RemoveAccount(ctx context.Context, acc sdk.AccountI)
	NewAccountWithAddress(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
}

type StakingKeeper interface {
	Delegate(goCtx context.Context, msg *stakingtypes.MsgDelegate) (*stakingtypes.MsgDelegateResponse, error)
	BeginRedelegate(goCtx context.Context, msg *stakingtypes.MsgBeginRedelegate) (*stakingtypes.MsgBeginRedelegateResponse, error)
	Undelegate(goCtx context.Context, msg *stakingtypes.MsgUndelegate) (*stakingtypes.MsgUndelegateResponse, error)
}

type StakingQuerier interface {
	Delegation(c context.Context, req *stakingtypes.QueryDelegationRequest) (*stakingtypes.QueryDelegationResponse, error)
}

//type GovKeeper interface {
//	AddVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress, options govtypes.WeightedVoteOptions) error
//	AddDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress, depositAmount sdk.Coins) (bool, error)
//}

type DistributionKeeper interface {
	SetWithdrawAddr(ctx sdk.Context, delegatorAddr sdk.AccAddress, withdrawAddr sdk.AccAddress) error
	WithdrawDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error)
	DelegationTotalRewards(c context.Context, req *distrtypes.QueryDelegationTotalRewardsRequest) (*distrtypes.QueryDelegationTotalRewardsResponse, error)
}

//type TransferKeeper interface {
//	Transfer(goCtx context.Context, msg *ibctypes.MsgTransfer) (*ibctypes.MsgTransferResponse, error)
//}

type ClientKeeper interface {
	GetClientState(ctx sdk.Context, clientID string) (exported.ClientState, bool)
	GetClientConsensusState(ctx sdk.Context, clientID string, height exported.Height) (exported.ConsensusState, bool)
}

type ConnectionKeeper interface {
	//GetConnection(ctx sdk.Context, connectionID string) (connectiontypes.ConnectionEnd, bool)
}

//type ChannelKeeper interface {
//	GetChannel(ctx sdk.Context, portID, channelID string) (types.Channel, bool)
//}
