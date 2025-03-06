package ante

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/eni-chain/go-eni/utils/helpers"
	"github.com/eni-chain/go-eni/x/evm/types/ethtx"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	coserrors "github.com/cosmos/cosmos-sdk/types/errors"
	accountkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	//"github.com/eni-chain/go-eni/app/antedecorators"
	"github.com/eni-chain/go-eni/utils"
	//"github.com/eni-chain/go-eni/utils/metrics"
	"github.com/eni-chain/go-eni/x/evm/derived"
	evmkeeper "github.com/eni-chain/go-eni/x/evm/keeper"
	evmtypes "github.com/eni-chain/go-eni/x/evm/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

// Accounts need to have at least 1Eni to force association. Note that account won't be charged.
const BalanceThreshold uint64 = 1000000

var BigBalanceThreshold *big.Int = new(big.Int).SetUint64(BalanceThreshold)
var BigBalanceThresholdMinus1 *big.Int = new(big.Int).SetUint64(BalanceThreshold - 1)

var SignerMap = map[derived.SignerVersion]func(*big.Int) ethtypes.Signer{
	derived.London: ethtypes.NewLondonSigner,
	derived.Cancun: ethtypes.NewCancunSigner,
}
var AllowedTxTypes = map[derived.SignerVersion][]uint8{
	derived.London: {ethtypes.LegacyTxType, ethtypes.AccessListTxType, ethtypes.DynamicFeeTxType},
	derived.Cancun: {ethtypes.LegacyTxType, ethtypes.AccessListTxType, ethtypes.DynamicFeeTxType, ethtypes.BlobTxType},
}

type EVMPreprocessDecorator struct {
	evmKeeper     *evmkeeper.Keeper
	accountKeeper ante.AccountKeeper
}

// todo authkeeper.AccountKeeper auth.AccountKeeper need check
func NewEVMPreprocessDecorator(evmKeeper *evmkeeper.Keeper, accountKeeper ante.AccountKeeper) *EVMPreprocessDecorator {
	return &EVMPreprocessDecorator{evmKeeper: evmKeeper, accountKeeper: accountKeeper}
}

//nolint:revive
func (p *EVMPreprocessDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	msg := evmtypes.MustGetEVMTransactionMessage(tx)
	if err := Preprocess(ctx, msg); err != nil {
		return ctx, err
	}

	// use infinite gas meter for EVM transaction because EVM handles gas checking from within
	if ctx.GasMeter() == nil {
		ctx = ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
	}

	derived := msg.Derived
	eniAddr := derived.SenderEniAddr
	evmAddr := derived.SenderEVMAddr
	ctx.EventManager().EmitEvent(sdk.NewEvent(evmtypes.EventTypeSigner,
		sdk.NewAttribute(evmtypes.AttributeKeyEvmAddress, evmAddr.Hex()),
		sdk.NewAttribute(evmtypes.AttributeKeyEniAddress, eniAddr.String())))

	// todo Need to deal with later
	// evm and eni address Associate implement
	//pubkey := derived.PubKey
	//isAssociateTx := derived.IsAssociate
	//associateHelper := helpers.NewAssociationHelper(p.evmKeeper, p.evmKeeper.BankKeeper(), p.accountKeeper)
	//_, isAssociated := p.evmKeeper.GetEVMAddress(ctx, eniAddr)
	//if isAssociateTx && isAssociated {
	//	return ctx, sdkerrors.Wrap(coserrors.ErrInvalidRequest, "account already has association set")
	//} else if isAssociateTx {
	//	// check if the account has enough balance (without charging)
	//	if !p.IsAccountBalancePositive(ctx, eniAddr, evmAddr) {
	//		metrics.IncrementAssociationError("associate_tx_insufficient_funds", evmtypes.NewAssociationMissingErr(eniAddr.String()))
	//		return ctx, sdkerrors.Wrap(coserrors.ErrInsufficientFunds, "account needs to have at least 1 wei to force association")
	//	}
	//	if err := associateHelper.AssociateAddresses(ctx, eniAddr, evmAddr, pubkey); err != nil {
	//		return ctx, err
	//	}
	//
	//	return ctx.WithPriority(antedecorators.EVMAssociatePriority), nil // short-circuit without calling next
	//} else if isAssociated {
	//	// noop; for readability
	//} else {
	//	// not associatedTx and not already associated
	//	if err := associateHelper.AssociateAddresses(ctx, eniAddr, evmAddr, pubkey); err != nil {
	//		return ctx, err
	//	}
	//}

	return next(ctx, tx, simulate)
}

func (p *EVMPreprocessDecorator) IsAccountBalancePositive(ctx sdk.Context, eniAddr sdk.AccAddress, evmAddr common.Address) bool {
	baseDenom := p.evmKeeper.GetBaseDenom(ctx)
	if amt := p.evmKeeper.BankKeeper().GetBalance(ctx, eniAddr, baseDenom).Amount; amt.IsPositive() {
		return true
	}
	if amt := p.evmKeeper.BankKeeper().GetBalance(ctx, sdk.AccAddress(evmAddr[:]), baseDenom).Amount; amt.IsPositive() {
		return true
	}
	//if amt := p.evmKeeper.BankKeeper().GetWeiBalance(ctx, eniAddr); amt.IsPositive() {
	//	return true
	//}
	//return p.evmKeeper.BankKeeper().GetWeiBalance(ctx, sdk.AccAddress(evmAddr[:])).IsPositive()
	return false
}

// stateless
func Preprocess(ctx sdk.Context, msgEVMTransaction *evmtypes.MsgEVMTransaction) error {
	if msgEVMTransaction.Derived != nil {
		if msgEVMTransaction.Derived.PubKey == nil {
			// this means the message has `Derived` set from the outside, in which case we should reject
			return coserrors.ErrInvalidPubKey
		}
		// already preprocessed
		return nil
	}
	txData, err := evmtypes.UnpackTxData(msgEVMTransaction.Data)
	if err != nil {
		return err
	}

	if atx, ok := txData.(*ethtx.AssociateTx); ok {
		V, R, S := atx.GetRawSignatureValues()
		V = new(big.Int).Add(V, utils.Big27)
		// Hash custom message passed in
		customMessageHash := crypto.Keccak256Hash([]byte(atx.CustomMessage))
		evmAddr, eniAddr, pubkey, err := helpers.GetAddresses(V, R, S, customMessageHash)
		if err != nil {
			return err
		}
		msgEVMTransaction.Derived = &derived.Derived{
			SenderEVMAddr: evmAddr,
			SenderEniAddr: eniAddr,
			PubKey:        &secp256k1.PubKey{Key: pubkey.Bytes()},
			Version:       derived.Cancun,
			IsAssociate:   true,
		}
		return nil
	}

	ethTx := ethtypes.NewTx(txData.AsEthereumData())
	chainID := ethTx.ChainId()
	chainCfg := evmtypes.DefaultChainConfig()
	ethCfg := chainCfg.EthereumConfig(chainID)
	version := GetVersion(ctx, ethCfg)
	signer := SignerMap[version](chainID)
	if !isTxTypeAllowed(version, ethTx.Type()) {
		return ethtypes.ErrInvalidChainId
	}

	var txHash common.Hash
	V, R, S := ethTx.RawSignatureValues()
	if ethTx.Protected() {
		V = AdjustV(V, ethTx.Type(), ethCfg.ChainID)
		txHash = signer.Hash(ethTx)
	} else {
		txHash = ethtypes.FrontierSigner{}.Hash(ethTx)
	}
	evmAddr, eniAddr, eniPubkey, err := helpers.GetAddresses(V, R, S, txHash)
	if err != nil {
		return err
	}
	msgEVMTransaction.Derived = &derived.Derived{
		SenderEVMAddr: evmAddr,
		SenderEniAddr: eniAddr,
		PubKey:        &secp256k1.PubKey{Key: eniPubkey.Bytes()},
		Version:       version,
		IsAssociate:   false,
	}
	return nil
}

//func (p *EVMPreprocessDecorator) AnteDeps(txDeps []sdkacltypes.AccessOperation, tx sdk.Tx, txIndex int, next sdk.AnteDepGenerator) (newTxDeps []sdkacltypes.AccessOperation, err error) {
//	msg := evmtypes.MustGetEVMTransactionMessage(tx)
//	return next(append(txDeps, sdkacltypes.AccessOperation{
//		AccessType:         sdkacltypes.AccessType_READ,
//		ResourceType:       sdkacltypes.ResourceType_KV_EVM_S2E,
//		IdentifierTemplate: hex.EncodeToString(evmtypes.EniAddressToEVMAddressKey(msg.Derived.SenderEniAddr)),
//	}, sdkacltypes.AccessOperation{
//		AccessType:         sdkacltypes.AccessType_WRITE,
//		ResourceType:       sdkacltypes.ResourceType_KV_EVM_S2E,
//		IdentifierTemplate: hex.EncodeToString(evmtypes.EniAddressToEVMAddressKey(msg.Derived.SenderEniAddr)),
//	}, sdkacltypes.AccessOperation{
//		AccessType:         sdkacltypes.AccessType_WRITE,
//		ResourceType:       sdkacltypes.ResourceType_KV_EVM_E2S,
//		IdentifierTemplate: hex.EncodeToString(evmtypes.EVMAddressToEniAddressKey(msg.Derived.SenderEVMAddr)),
//	}, sdkacltypes.AccessOperation{
//		AccessType:         sdkacltypes.AccessType_READ,
//		ResourceType:       sdkacltypes.ResourceType_KV_BANK_BALANCES,
//		IdentifierTemplate: hex.EncodeToString(banktypes.CreateAccountBalancesPrefix(msg.Derived.SenderEniAddr)),
//	}, sdkacltypes.AccessOperation{
//		AccessType:         sdkacltypes.AccessType_WRITE,
//		ResourceType:       sdkacltypes.ResourceType_KV_BANK_BALANCES,
//		IdentifierTemplate: hex.EncodeToString(banktypes.CreateAccountBalancesPrefix(msg.Derived.SenderEniAddr)),
//	}, sdkacltypes.AccessOperation{
//		AccessType:         sdkacltypes.AccessType_READ,
//		ResourceType:       sdkacltypes.ResourceType_KV_BANK_BALANCES,
//		IdentifierTemplate: hex.EncodeToString(banktypes.CreateAccountBalancesPrefix(msg.Derived.SenderEVMAddr[:])),
//	}, sdkacltypes.AccessOperation{
//		AccessType:         sdkacltypes.AccessType_WRITE,
//		ResourceType:       sdkacltypes.ResourceType_KV_BANK_BALANCES,
//		IdentifierTemplate: hex.EncodeToString(banktypes.CreateAccountBalancesPrefix(msg.Derived.SenderEVMAddr[:])),
//	}, sdkacltypes.AccessOperation{
//		AccessType:         sdkacltypes.AccessType_READ,
//		ResourceType:       sdkacltypes.ResourceType_KV_AUTH_ADDRESS_STORE,
//		IdentifierTemplate: hex.EncodeToString(authtypes.AddressStoreKey(msg.Derived.SenderEniAddr)),
//	}, sdkacltypes.AccessOperation{
//		AccessType:         sdkacltypes.AccessType_WRITE,
//		ResourceType:       sdkacltypes.ResourceType_KV_AUTH_ADDRESS_STORE,
//		IdentifierTemplate: hex.EncodeToString(authtypes.AddressStoreKey(msg.Derived.SenderEniAddr)),
//	}, sdkacltypes.AccessOperation{
//		AccessType:         sdkacltypes.AccessType_READ,
//		ResourceType:       sdkacltypes.ResourceType_KV_AUTH_ADDRESS_STORE,
//		IdentifierTemplate: hex.EncodeToString(authtypes.AddressStoreKey(msg.Derived.SenderEVMAddr[:])),
//	}, sdkacltypes.AccessOperation{
//		AccessType:         sdkacltypes.AccessType_WRITE,
//		ResourceType:       sdkacltypes.ResourceType_KV_AUTH_ADDRESS_STORE,
//		IdentifierTemplate: hex.EncodeToString(authtypes.AddressStoreKey(msg.Derived.SenderEVMAddr[:])),
//	}, sdkacltypes.AccessOperation{
//		AccessType:         sdkacltypes.AccessType_READ,
//		ResourceType:       sdkacltypes.ResourceType_KV_EVM_NONCE,
//		IdentifierTemplate: hex.EncodeToString(append(evmtypes.NonceKeyPrefix, msg.Derived.SenderEVMAddr[:]...)),
//	}), tx, txIndex)
//}

func isTxTypeAllowed(version derived.SignerVersion, txType uint8) bool {
	for _, t := range AllowedTxTypes[version] {
		if t == txType {
			return true
		}
	}
	return false
}

func AdjustV(V *big.Int, txType uint8, chainID *big.Int) *big.Int {
	// Non-legacy TX always needs to be bumped by 27
	if txType != ethtypes.LegacyTxType {
		return new(big.Int).Add(V, utils.Big27)
	}

	// legacy TX needs to be adjusted based on chainID
	V = new(big.Int).Sub(V, new(big.Int).Mul(chainID, utils.Big2))
	return V.Sub(V, utils.Big8)
}

func GetVersion(ctx sdk.Context, ethCfg *params.ChainConfig) derived.SignerVersion {
	blockNum := big.NewInt(ctx.BlockHeight())
	ts := uint64(ctx.BlockTime().Unix())
	switch {
	case ethCfg.IsCancun(blockNum, ts):
		return derived.Cancun
	default:
		return derived.London
	}
}

type EVMAddressDecorator struct {
	evmKeeper     *evmkeeper.Keeper
	accountKeeper *accountkeeper.AccountKeeper
}

func NewEVMAddressDecorator(evmKeeper *evmkeeper.Keeper, accountKeeper *accountkeeper.AccountKeeper) *EVMAddressDecorator {
	return &EVMAddressDecorator{evmKeeper: evmKeeper, accountKeeper: accountKeeper}
}

//nolint:revive
func (p *EVMAddressDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	//sigTx, ok := tx.(authsigning.SigVerifiableTx)
	//if !ok {
	//	return ctx, sdkerrors.Wrap(cosmoserrors.ErrTxDecode, "invalid tx type")
	//}
	//signers, err := sigTx.GetSigners()
	//if err != nil {
	//	ctx.Logger().Error(fmt.Sprintf("get signers failed: %s", err))
	//	return ctx, sdkerrors.Wrap(cosmoserrors.ErrTxDecode, "invalid tx signers")
	//}
	//for _, signer := range signers {
	//if evmAddr, associated := p.evmKeeper.GetEVMAddress(ctx, signer); associated {
	//	ctx.EventManager().EmitEvent(sdk.NewEvent(evmtypes.EventTypeSigner,
	//		sdk.NewAttribute(evmtypes.AttributeKeyEvmAddress, evmAddr.Hex()),
	//		sdk.NewAttribute(evmtypes.AttributeKeyEniAddress, signer.String())))
	//	continue
	//}
	//acc := p.accountKeeper.GetAccount(ctx, signer)
	//if acc.GetPubKey() == nil {
	//	ctx.Logger().Error(fmt.Sprintf("missing pubkey for %s", signer.String()))
	//	ctx.EventManager().EmitEvent(sdk.NewEvent(evmtypes.EventTypeSigner,
	//		sdk.NewAttribute(evmtypes.AttributeKeyEniAddress, signer.String())))
	//	continue
	//}
	//pk, err := btcec.ParsePubKey(acc.GetPubKey().Bytes(), btcec.S256())
	//if err != nil {
	//	ctx.Logger().Debug(fmt.Sprintf("failed to parse pubkey for %s, likely due to the fact that it isn't on secp256k1 curve", acc.GetPubKey()), "err", err)
	//	ctx.EventManager().EmitEvent(sdk.NewEvent(evmtypes.EventTypeSigner,
	//		sdk.NewAttribute(evmtypes.AttributeKeyEniAddress, signer.String())))
	//	continue
	//}
	//evmAddr, err := helpers.PubkeyToEVMAddress(pk.SerializeUncompressed())
	//if err != nil {
	//	ctx.Logger().Error(fmt.Sprintf("failed to get EVM address from pubkey due to %s", err))
	//	ctx.EventManager().EmitEvent(sdk.NewEvent(evmtypes.EventTypeSigner,
	//		sdk.NewAttribute(evmtypes.AttributeKeyEniAddress, signer.String())))
	//	continue
	//}
	//ctx.EventManager().EmitEvent(sdk.NewEvent(evmtypes.EventTypeSigner,
	//	sdk.NewAttribute(evmtypes.AttributeKeyEvmAddress, evmAddr.Hex()),
	//	sdk.NewAttribute(evmtypes.AttributeKeyEniAddress, signer.String())))
	//p.evmKeeper.SetAddressMapping(ctx, signer, evmAddr)
	//associationHelper := helpers.NewAssociationHelper(p.evmKeeper, p.evmKeeper.BankKeeper(), p.accountKeeper)
	//if err := associationHelper.MigrateBalance(ctx, evmAddr, signer); err != nil {
	//	ctx.Logger().Error(fmt.Sprintf("failed to migrate EVM address balance (%s) %s", evmAddr.Hex(), err))
	//	return ctx, err
	//}
	//if evmtypes.IsTxMsgAssociate(tx) {
	//	// check if there is non-zero balance
	//	if !p.evmKeeper.BankKeeper().GetBalance(ctx, signer, sdk.MustGetBaseDenom()).IsPositive() && !p.evmKeeper.BankKeeper().GetWeiBalance(ctx, signer).IsPositive() {
	//		return ctx, sdkerrors.Wrap(cosmoserrors.ErrInsufficientFunds, "account needs to have at least 1 wei to force association")
	//	}
	//}
	//}
	return next(ctx, tx, simulate)
}
