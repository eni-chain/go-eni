package app

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"

	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	//ibcante "github.com/cosmos/ibc-go/v3/modules/core/ante"
	//ibckeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"

	evmante "github.com/eni-chain/go-eni/x/evm/ante"
	evmkeeper "github.com/eni-chain/go-eni/x/evm/keeper"
	//"github.com/eni-chain/go-eni/x/oracle"
	//oraclekeeper "github.com/eni-chain/go-eni/x/oracle/keeper"
)

// HandlerOptions extend the SDK's AnteHandler options by requiring the IBC
// channel keeper.
type HandlerOptions struct {
	ante.HandlerOptions

	EVMKeeper        *evmkeeper.Keeper
	AppAccountKeeper *authkeeper.AccountKeeper
	LatestCtxGetter  func() sdk.Context
}

func NewAnteHandlerAndDepGenerator(options HandlerOptions) (sdk.AnteHandler, error) {
	if options.AccountKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "account keeper is required for ante builder")
	}

	if options.BankKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "bank keeper is required for ante builder")
	}

	if options.SignModeHandler == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	}

	if options.EVMKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "evm keeper is required for ante builder")
	}
	if options.LatestCtxGetter == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "latest context getter is required for ante builder")
	}
	if options.AppAccountKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "account keeper is required for ante builder")
	}
	sigGasConsumer := options.SigGasConsumer
	if sigGasConsumer == nil {
		sigGasConsumer = ante.DefaultSigVerificationGasConsumer
	}

	anteDecorators := []sdk.AnteDecorator{
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker),
		ante.NewSetPubKeyDecorator(options.AccountKeeper), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		evmante.NewEVMAddressDecorator(options.EVMKeeper, options.AppAccountKeeper),
	}

	anteHandler := sdk.ChainAnteDecorators(anteDecorators...)
	evmAnteDecorators := []sdk.AnteDecorator{
		evmante.NewEVMPreprocessDecorator(options.EVMKeeper, options.AppAccountKeeper),
		evmante.NewBasicDecorator(options.EVMKeeper),
		evmante.NewEVMFeeCheckDecorator(options.EVMKeeper),
		evmante.NewEVMSigVerifyDecorator(options.EVMKeeper, options.LatestCtxGetter),
		evmante.NewGasLimitDecorator(options.EVMKeeper),
	}
	evmAnteHandler := sdk.ChainAnteDecorators(evmAnteDecorators...)

	router := evmante.NewEVMRouterDecorator(anteHandler, evmAnteHandler)

	return router.AnteHandle, nil
}
