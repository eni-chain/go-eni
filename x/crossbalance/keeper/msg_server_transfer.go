package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	//"github.com/cosmos/cosmos-sdk/utils/helpers"
	"github.com/eni-chain/go-eni/x/crossbalance/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func (k msgServer) TransferEniToEvm(goCtx context.Context, msg *types.MsgTransferEniToEvm) (*types.MsgTransferEniToEvmResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Convert sender address
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	// Convert EVM address
	evmAddr := common.HexToAddress(msg.To)

	// Step 1: Transfer from ENI address to ENI address (burn)
	// Create a burn address (you can use a specific burn address or generate one)
	burnAddr := sdk.AccAddress(crypto.Keccak256Hash([]byte("burn")).Bytes())

	// Check if sender has enough balance
	if !k.BankKeeper.HasBalance(ctx, sender, msg.Amount) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds, "insufficient balance")
	}

	// Transfer to burn address
	if err := k.BankKeeper.SendCoins(ctx, sender, burnAddr, sdk.NewCoins(msg.Amount)); err != nil {
		return nil, errorsmod.Wrapf(err, "failed to send coins to burn address")
	}

	// Step 2: Transfer from EVM address to EVM address (mint)
	// Convert EVM address to Cosmos address for the recipient
	recipientAddr := sdk.AccAddress(evmAddr.Bytes())

	// Mint coins to the recipient's EVM address
	// todo evm bankeeper mintcoins
	if err := k.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(msg.Amount)); err != nil {
		return nil, errorsmod.Wrapf(err, "failed to mint coins")
	}

	// todo evm bankeeper sendcoinsfrommoduletoaccount
	if err := k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipientAddr, sdk.NewCoins(msg.Amount)); err != nil {
		return nil, errorsmod.Wrapf(err, "failed to send coins to recipient")
	}

	return &types.MsgTransferEniToEvmResponse{}, nil
}

func (k msgServer) TransferEvmToEni(goCtx context.Context, msg *types.MsgTransferEvmToEni) (*types.MsgTransferEvmToEniResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Convert EVM address
	//evmAddr := common.HexToAddress(msg.From)

	// Convert public key to Cosmos address
	pubKey, err := crypto.UnmarshalPubkey(msg.Ak)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidPubKey, "invalid public key (%s)", err)
	}

	//evmAddr, eniAddr, eniPubkey, err := helpers.GetAddressesFromPubkeyBytes(msg.Ak)
	// Get the Cosmos address from the public key
	//todo helpers.GetAddressesFromPubkeyBytes
	//1. evmaddress burn amount to evm address
	//2. eniaddr mint amount to eni address
	cosmosFrom := sdk.AccAddress(crypto.PubkeyToAddress(*pubKey).Bytes())

	// Convert destination address
	to, err := sdk.AccAddressFromBech32(msg.To)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid destination address (%s)", err)
	}

	// Check if sender has enough balance
	if !k.BankKeeper.HasBalance(ctx, cosmosFrom, msg.Amount) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds, "insufficient balance")
	}

	// Deduct from sender's balance and add to recipient's balance
	if err := k.BankKeeper.SendCoins(ctx, cosmosFrom, to, sdk.NewCoins(msg.Amount)); err != nil {
		return nil, errorsmod.Wrapf(err, "failed to send coins")
	}

	return &types.MsgTransferEvmToEniResponse{}, nil
}
