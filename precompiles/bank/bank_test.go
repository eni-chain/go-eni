package bank_test

//import (
//	"embed"
//	"encoding/hex"
//	"fmt"
//	"math/big"
//	"strings"
//	"testing"
//	"time"
//
//	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
//
//	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	"github.com/cosmos/cosmos-sdk/types/tx/signing"
//	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
//	"github.com/eni-chain/go-eni/precompiles/bank"
//	pcommon "github.com/eni-chain/go-eni/precompiles/common"
//	testkeeper "github.com/eni-chain/go-eni/testutil/keeper"
//	"github.com/cosmos/cosmos-sdk/x/evm/ante"
//	"github.com/cosmos/cosmos-sdk/x/evm/keeper"
//	"github.com/cosmos/cosmos-sdk/x/evm/state"
//	"github.com/cosmos/cosmos-sdk/x/evm/types"
//	//"github.com/cosmos/cosmos-sdk/x/evm/types/ethtx"
//	"github.com/ethereum/go-ethereum/common"
//	ethtypes "github.com/ethereum/go-ethereum/core/types"
//	"github.com/ethereum/go-ethereum/core/vm"
//	"github.com/ethereum/go-ethereum/crypto"
//	"github.com/stretchr/testify/require"
//	tmtypes "github.com/tendermint/tendermint/proto/tendermint/types"
//)
//
////go:embed abi.json
//var f embed.FS
//
//type mockTx struct {
//	msgs    []sdk.Msg
//	signers []sdk.AccAddress
//}
//
//func (tx mockTx) GetMsgs() []sdk.Msg                              { return tx.msgs }
//func (tx mockTx) ValidateBasic() error                            { return nil }
//func (tx mockTx) GetSigners() []sdk.AccAddress                    { return tx.signers }
//func (tx mockTx) GetPubKeys() ([]cryptotypes.PubKey, error)       { return nil, nil }
//func (tx mockTx) GetSignaturesV2() ([]signing.SignatureV2, error) { return nil, nil }
//
//func TestRun(t *testing.T) {
//	testApp := testkeeper.EVMTestApp
//	ctx := testApp.NewContext(false, tmtypes.Header{}).WithBlockHeight(2)
//	k := &testApp.EvmKeeper
//
//	// Setup sender addresses and environment
//	privKey := testkeeper.MockPrivateKey()
//	testPrivHex := hex.EncodeToString(privKey.Bytes())
//	senderAddr, senderEVMAddr := testkeeper.PrivateKeyToAddresses(privKey)
//	k.SetAddressMapping(ctx, senderAddr, senderEVMAddr)
//	err := k.BankKeeper().MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin("ufoo", sdk.NewInt(10000000))))
//	require.Nil(t, err)
//	err = k.BankKeeper().SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, sdk.NewCoins(sdk.NewCoin("ufoo", sdk.NewInt(10000000))))
//	require.Nil(t, err)
//	err = k.BankKeeper().MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin("ueni", sdk.NewInt(10000000))))
//	require.Nil(t, err)
//	err = k.BankKeeper().SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, sdk.NewCoins(sdk.NewCoin("ueni", sdk.NewInt(10000000))))
//	require.Nil(t, err)
//
//	// Setup receiving addresses
//	eniAddr, evmAddr := testkeeper.MockAddressPair()
//	k.SetAddressMapping(ctx, eniAddr, evmAddr)
//	p, err := bank.NewPrecompile(k.BankKeeper(), bankkeeper.NewMsgServerImpl(k.BankKeeper()), k, k.AccountKeeper())
//	require.Nil(t, err)
//	statedb := state.NewDBImpl(ctx, k, true)
//	evm := vm.EVM{
//		StateDB:   statedb,
//		TxContext: vm.TxContext{Origin: senderEVMAddr},
//	}
//
//	// Precompile send test
//	send, err := p.ABI.MethodById(p.GetExecutor().(*bank.PrecompileExecutor).SendID)
//	require.Nil(t, err)
//	args, err := send.Inputs.Pack(senderEVMAddr, evmAddr, "ueni", big.NewInt(25))
//	require.Nil(t, err)
//	_, _, err = p.RunAndCalculateGas(&evm, senderEVMAddr, senderEVMAddr, append(p.GetExecutor().(*bank.PrecompileExecutor).SendID, args...), 100000, nil, nil, true, false) // should error because of read only call
//	require.NotNil(t, err)
//	_, _, err = p.RunAndCalculateGas(&evm, senderEVMAddr, senderEVMAddr, append(p.GetExecutor().(*bank.PrecompileExecutor).SendID, args...), 100000, big.NewInt(1), nil, false, false) // should error because it's not payable
//	require.NotNil(t, err)
//	_, _, err = p.RunAndCalculateGas(&evm, senderEVMAddr, senderEVMAddr, append(p.GetExecutor().(*bank.PrecompileExecutor).SendID, args...), 100000, nil, nil, false, false) // should error because address is not whitelisted
//	require.NotNil(t, err)
//	invalidDenomArgs, err := send.Inputs.Pack(senderEVMAddr, evmAddr, "", big.NewInt(25))
//	require.Nil(t, err)
//	_, _, err = p.RunAndCalculateGas(&evm, senderEVMAddr, senderEVMAddr, append(p.GetExecutor().(*bank.PrecompileExecutor).SendID, invalidDenomArgs...), 100000, nil, nil, false, false) // should error because denom is empty
//	require.NotNil(t, err)
//
//	// Precompile sendNative test error
//	sendNative, err := p.ABI.MethodById(p.GetExecutor().(*bank.PrecompileExecutor).SendNativeID)
//	require.Nil(t, err)
//	eniAddrString := eniAddr.String()
//	argsNativeError, err := sendNative.Inputs.Pack(eniAddrString)
//	require.Nil(t, err)
//	// 0 amount disallowed
//	_, _, err = p.RunAndCalculateGas(&evm, senderEVMAddr, senderEVMAddr, append(p.GetExecutor().(*bank.PrecompileExecutor).SendNativeID, argsNativeError...), 100000, big.NewInt(0), nil, false, false)
//	require.NotNil(t, err)
//	argsNativeError, err = sendNative.Inputs.Pack("")
//	require.Nil(t, err)
//	_, _, err = p.RunAndCalculateGas(&evm, senderEVMAddr, senderEVMAddr, append(p.GetExecutor().(*bank.PrecompileExecutor).SendNativeID, argsNativeError...), 100000, big.NewInt(100), nil, false, false)
//	require.NotNil(t, err)
//	argsNativeError, err = sendNative.Inputs.Pack("invalidaddr")
//	require.Nil(t, err)
//	_, _, err = p.RunAndCalculateGas(&evm, senderEVMAddr, senderEVMAddr, append(p.GetExecutor().(*bank.PrecompileExecutor).SendNativeID, argsNativeError...), 100000, big.NewInt(100), nil, false, false)
//	require.NotNil(t, err)
//	argsNativeError, err = sendNative.Inputs.Pack(senderAddr.String())
//	require.Nil(t, err)
//	_, _, err = p.RunAndCalculateGas(&evm, evmAddr, evmAddr, append(p.GetExecutor().(*bank.PrecompileExecutor).SendNativeID, argsNativeError...), 100000, big.NewInt(100), nil, false, false)
//	require.NotNil(t, err)
//	_, _, err = p.RunAndCalculateGas(&evm, evmAddr, evmAddr, append(p.GetExecutor().(*bank.PrecompileExecutor).SendNativeID, argsNativeError...), 100000, big.NewInt(100), nil, true, false)
//	require.NotNil(t, err)
//	_, _, err = p.RunAndCalculateGas(&evm, evmAddr, evmAddr, append(p.GetExecutor().(*bank.PrecompileExecutor).SendNativeID, argsNativeError...), 100000, big.NewInt(100), nil, false, true)
//	require.NotNil(t, err)
//
//	// Send native 10_000_000_000_100, split into 10 ueni 100wei
//	// Test payable with eth LegacyTx
//	abi := pcommon.MustGetABI(f, "abi.json")
//	argsNative, err := abi.Pack(bank.SendNativeMethod, eniAddr.String())
//	require.Nil(t, err)
//	require.Nil(t, err)
//	key, _ := crypto.HexToECDSA(testPrivHex)
//	addr := common.HexToAddress(bank.BankAddress)
//	txData := ethtypes.LegacyTx{
//		GasPrice: big.NewInt(1000000000000),
//		Gas:      200000,
//		To:       &addr,
//		Value:    big.NewInt(10_000_000_000_100),
//		Data:     argsNative,
//		Nonce:    0,
//	}
//	chainID := k.ChainID(ctx)
//	chainCfg := types.DefaultChainConfig()
//	ethCfg := chainCfg.EthereumConfig(chainID)
//	blockNum := big.NewInt(ctx.BlockHeight())
//	signer := ethtypes.MakeSigner(ethCfg, blockNum, uint64(ctx.BlockTime().Unix()))
//	tx, err := ethtypes.SignTx(ethtypes.NewTx(&txData), signer, key)
//	require.Nil(t, err)
//	txwrapper, err := ethtx.NewLegacyTx(tx)
//	require.Nil(t, err)
//	req, err := types.NewMsgEVMTransaction(txwrapper)
//	require.Nil(t, err)
//
//	// send the transaction
//	msgServer := keeper.NewMsgServerImpl(k)
//	ante.Preprocess(ctx, req)
//	ctx = ctx.WithEventManager(sdk.NewEventManager())
//	ctx, err = ante.NewEVMFeeCheckDecorator(k).AnteHandle(ctx, mockTx{msgs: []sdk.Msg{req}}, false, func(sdk.Context, sdk.Tx, bool) (sdk.Context, error) {
//		return ctx, nil
//	})
//	require.Nil(t, err)
//	res, err := msgServer.EVMTransaction(sdk.WrapSDKContext(ctx), req)
//	require.Nil(t, err)
//	require.Empty(t, res.VmError)
//
//	evts := ctx.EventManager().ABCIEvents()
//
//	for _, evt := range evts {
//		var lines []string
//		for _, attr := range evt.Attributes {
//			lines = append(lines, fmt.Sprintf("%s=%s", string(attr.Key), string(attr.Value)))
//		}
//		fmt.Printf("type=%s\t%s\n", evt.Type, strings.Join(lines, "\t"))
//	}
//
//	var expectedEvts sdk.Events = []sdk.Event{
//		// gas is sent from sender
//		banktypes.NewCoinSpentEvent(senderAddr, sdk.NewCoins(sdk.NewCoin("ueni", sdk.NewInt(200000)))),
//		// wei events
//		banktypes.NewWeiSpentEvent(senderAddr, sdk.NewInt(100)),
//		banktypes.NewWeiReceivedEvent(eniAddr, sdk.NewInt(100)),
//		sdk.NewEvent(
//			banktypes.EventTypeWeiTransfer,
//			sdk.NewAttribute(banktypes.AttributeKeyRecipient, eniAddr.String()),
//			sdk.NewAttribute(banktypes.AttributeKeySender, senderAddr.String()),
//			sdk.NewAttribute(sdk.AttributeKeyAmount, sdk.NewInt(100).String()),
//		),
//		// sender sends coin to the receiver
//		banktypes.NewCoinSpentEvent(senderAddr, sdk.NewCoins(sdk.NewCoin("ueni", sdk.NewInt(10)))),
//		banktypes.NewCoinReceivedEvent(eniAddr, sdk.NewCoins(sdk.NewCoin("ueni", sdk.NewInt(10)))),
//		sdk.NewEvent(
//			banktypes.EventTypeTransfer,
//			sdk.NewAttribute(banktypes.AttributeKeyRecipient, eniAddr.String()),
//			sdk.NewAttribute(banktypes.AttributeKeySender, senderAddr.String()),
//			sdk.NewAttribute(sdk.AttributeKeyAmount, sdk.NewCoin("ueni", sdk.NewInt(10)).String()),
//		),
//		sdk.NewEvent(
//			sdk.EventTypeMessage,
//			sdk.NewAttribute(banktypes.AttributeKeySender, senderAddr.String()),
//		),
//		// gas refund to the sender
//		banktypes.NewCoinReceivedEvent(senderAddr, sdk.NewCoins(sdk.NewCoin("ueni", sdk.NewInt(132401)))),
//		// tip is paid to the validator
//		banktypes.NewCoinReceivedEvent(sdk.MustAccAddressFromBech32("eni1v4mx6hmrda5kucnpwdjsqqqqqqqqqqqqzws0wt"), sdk.NewCoins(sdk.NewCoin("ueni", sdk.NewInt(67599)))),
//	}
//	require.EqualValues(t, expectedEvts.ToABCIEvents(), evts)
//
//	// Use precompile balance to verify sendNative ueni amount succeeded
//	balance, err := p.ABI.MethodById(p.GetExecutor().(*bank.PrecompileExecutor).BalanceID)
//	require.Nil(t, err)
//	args, err = balance.Inputs.Pack(evmAddr, "ueni")
//	require.Nil(t, err)
//	precompileRes, _, err := p.RunAndCalculateGas(&evm, common.Address{}, common.Address{}, append(p.GetExecutor().(*bank.PrecompileExecutor).BalanceID, args...), 100000, nil, nil, false, false)
//	require.Nil(t, err)
//	is, err := balance.Outputs.Unpack(precompileRes)
//	require.Nil(t, err)
//	require.Equal(t, 1, len(is))
//	require.Equal(t, big.NewInt(10), is[0].(*big.Int))
//	weiBalance := k.BankKeeper().GetWeiBalance(ctx, eniAddr)
//	require.Equal(t, big.NewInt(100), weiBalance.BigInt())
//
//	newAddr, _ := testkeeper.MockAddressPair()
//	require.Nil(t, k.AccountKeeper().GetAccount(ctx, newAddr))
//	argsNewAccount, err := sendNative.Inputs.Pack(newAddr.String())
//	require.Nil(t, err)
//	require.Nil(t, k.BankKeeper().SendCoins(ctx, eniAddr, k.GetEniAddressOrDefault(ctx, p.Address()), sdk.NewCoins(sdk.NewCoin("ueni", sdk.OneInt()))))
//	_, _, err = p.RunAndCalculateGas(&evm, evmAddr, evmAddr, append(p.GetExecutor().(*bank.PrecompileExecutor).SendNativeID, argsNewAccount...), 100000, big.NewInt(1), nil, false, false)
//	require.Nil(t, err)
//	// should create account if not exists
//	require.NotNil(t, k.AccountKeeper().GetAccount(statedb.Ctx(), newAddr))
//
//	// test get all balances
//	allBalances, err := p.ABI.MethodById(p.GetExecutor().(*bank.PrecompileExecutor).AllBalancesID)
//	require.Nil(t, err)
//	args, err = allBalances.Inputs.Pack(senderEVMAddr)
//	require.Nil(t, err)
//	precompileRes, _, err = p.RunAndCalculateGas(&evm, common.Address{}, common.Address{}, append(p.GetExecutor().(*bank.PrecompileExecutor).AllBalancesID, args...), 100000, nil, nil, false, false)
//	require.Nil(t, err)
//	balances, err := allBalances.Outputs.Unpack(precompileRes)
//	require.Nil(t, err)
//	require.Equal(t, 1, len(balances))
//	parsedBalances := balances[0].([]struct {
//		Amount *big.Int `json:"amount"`
//		Denom  string   `json:"denom"`
//	})
//
//	// TODO: Not sure why changing the keyword causes this index to change
//	require.Equal(t, 2, len(parsedBalances))
//	require.Equal(t, bank.CoinBalance{
//		Amount: big.NewInt(10000000),
//		Denom:  "ufoo",
//	}, bank.CoinBalance(parsedBalances[1]))
//	require.Equal(t, bank.CoinBalance{
//		Amount: big.NewInt(9932390),
//		Denom:  "ueni",
//	}, bank.CoinBalance(parsedBalances[0]))
//
//	// Verify errors properly raised on bank balance calls with incorrect inputs
//	_, _, err = p.RunAndCalculateGas(&evm, common.Address{}, common.Address{}, append(p.GetExecutor().(*bank.PrecompileExecutor).BalanceID, args[:1]...), 100000, nil, nil, false, false)
//	require.NotNil(t, err)
//	args, err = balance.Inputs.Pack(evmAddr, "")
//	require.Nil(t, err)
//	_, _, err = p.RunAndCalculateGas(&evm, common.Address{}, common.Address{}, append(p.GetExecutor().(*bank.PrecompileExecutor).BalanceID, args...), 100000, nil, nil, false, false)
//	require.NotNil(t, err)
//
//	// invalid input
//	_, _, err = p.RunAndCalculateGas(&evm, common.Address{}, common.Address{}, []byte{1, 2, 3, 4}, 100000, nil, nil, false, false)
//	require.NotNil(t, err)
//}
//
//func TestSendForUnlinkedReceiver(t *testing.T) {
//	testApp := testkeeper.EVMTestApp
//	ctx := testApp.NewContext(false, tmtypes.Header{}).WithBlockHeight(2)
//	k := &testApp.EvmKeeper
//
//	// Setup sender addresses and environment
//	privKey := testkeeper.MockPrivateKey()
//	// testPrivHex := hex.EncodeToString(privKey.Bytes())
//	senderAddr, senderEVMAddr := testkeeper.PrivateKeyToAddresses(privKey)
//	k.SetAddressMapping(ctx, senderAddr, senderEVMAddr)
//	err := k.BankKeeper().MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin("ufoo", sdk.NewInt(10000000))))
//	require.Nil(t, err)
//	err = k.BankKeeper().SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, sdk.NewCoins(sdk.NewCoin("ufoo", sdk.NewInt(10000000))))
//	require.Nil(t, err)
//	err = k.BankKeeper().MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin("ueni", sdk.NewInt(10000000))))
//	require.Nil(t, err)
//	err = k.BankKeeper().SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, sdk.NewCoins(sdk.NewCoin("ueni", sdk.NewInt(10000000))))
//	require.Nil(t, err)
//
//	_, pointerAddr := testkeeper.MockAddressPair()
//	k.SetERC20NativePointer(ctx, "ufoo", pointerAddr)
//
//	// Setup receiving addresses - NOT linked
//	_, evmAddr := testkeeper.MockAddressPair()
//	p, err := bank.NewPrecompile(k.BankKeeper(), bankkeeper.NewMsgServerImpl(k.BankKeeper()), k, k.AccountKeeper())
//	require.Nil(t, err)
//	statedb := state.NewDBImpl(ctx, k, true)
//	evm := vm.EVM{
//		StateDB:   statedb,
//		TxContext: vm.TxContext{Origin: senderEVMAddr},
//	}
//
//	// Precompile send test
//	send, err := p.ABI.MethodById(p.GetExecutor().(*bank.PrecompileExecutor).SendID)
//	require.Nil(t, err)
//	args, err := send.Inputs.Pack(senderEVMAddr, evmAddr, "ufoo", big.NewInt(100))
//	require.Nil(t, err)
//	_, _, err = p.RunAndCalculateGas(&evm, pointerAddr, pointerAddr, append(p.GetExecutor().(*bank.PrecompileExecutor).SendID, args...), 100000, nil, nil, false, false) // should not error
//	require.Nil(t, err)
//
//	// Use precompile balance to verify sendNative ueni amount succeeded
//	balance, err := p.ABI.MethodById(p.GetExecutor().(*bank.PrecompileExecutor).BalanceID)
//	require.Nil(t, err)
//	args, err = balance.Inputs.Pack(evmAddr, "ufoo")
//	require.Nil(t, err)
//	precompileRes, _, err := p.RunAndCalculateGas(&evm, common.Address{}, common.Address{}, append(p.GetExecutor().(*bank.PrecompileExecutor).BalanceID, args...), 100000, nil, nil, false, false)
//	require.Nil(t, err)
//	is, err := balance.Outputs.Unpack(precompileRes)
//	require.Nil(t, err)
//	require.Equal(t, 1, len(is))
//	require.Equal(t, big.NewInt(100), is[0].(*big.Int))
//
//	// test get all balances
//	allBalances, err := p.ABI.MethodById(p.GetExecutor().(*bank.PrecompileExecutor).AllBalancesID)
//	require.Nil(t, err)
//	args, err = allBalances.Inputs.Pack(senderEVMAddr)
//	require.Nil(t, err)
//	precompileRes, _, err = p.RunAndCalculateGas(&evm, common.Address{}, common.Address{}, append(p.GetExecutor().(*bank.PrecompileExecutor).AllBalancesID, args...), 100000, nil, nil, false, false)
//	require.Nil(t, err)
//	balances, err := allBalances.Outputs.Unpack(precompileRes)
//	require.Nil(t, err)
//	require.Equal(t, 1, len(balances))
//	parsedBalances := balances[0].([]struct {
//		Amount *big.Int `json:"amount"`
//		Denom  string   `json:"denom"`
//	})
//
//	require.Equal(t, 2, len(parsedBalances))
//	require.Equal(t, bank.CoinBalance{
//		Amount: big.NewInt(9999900),
//		Denom:  "ufoo",
//	}, bank.CoinBalance(parsedBalances[1]))
//	require.Equal(t, bank.CoinBalance{
//		Amount: big.NewInt(10000000),
//		Denom:  "ueni",
//	}, bank.CoinBalance(parsedBalances[0]))
//
//	// Verify errors properly raised on bank balance calls with incorrect inputs
//	_, _, err = p.RunAndCalculateGas(&evm, common.Address{}, common.Address{}, append(p.GetExecutor().(*bank.PrecompileExecutor).BalanceID, args[:1]...), 100000, nil, nil, false, false)
//	require.NotNil(t, err)
//	args, err = balance.Inputs.Pack(evmAddr, "")
//	require.Nil(t, err)
//	_, _, err = p.RunAndCalculateGas(&evm, common.Address{}, common.Address{}, append(p.GetExecutor().(*bank.PrecompileExecutor).BalanceID, args...), 100000, nil, nil, false, false)
//	require.NotNil(t, err)
//
//	// invalid input
//	_, _, err = p.RunAndCalculateGas(&evm, common.Address{}, common.Address{}, []byte{1, 2, 3, 4}, 100000, nil, nil, false, false)
//	require.NotNil(t, err)
//}
//
//func TestMetadata(t *testing.T) {
//	k := &testkeeper.EVMTestApp.EvmKeeper
//	ctx := testkeeper.EVMTestApp.GetContextForDeliverTx([]byte{}).WithBlockTime(time.Now())
//	k.BankKeeper().SetDenomMetaData(ctx, banktypes.Metadata{Name: "ENI", Symbol: "ueni", Base: "ueni"})
//	p, err := bank.NewPrecompile(k.BankKeeper(), bankkeeper.NewMsgServerImpl(k.BankKeeper()), k, k.AccountKeeper())
//	require.Nil(t, err)
//	statedb := state.NewDBImpl(ctx, k, true)
//	evm := vm.EVM{
//		StateDB: statedb,
//	}
//	name, err := p.ABI.MethodById(p.GetExecutor().(*bank.PrecompileExecutor).NameID)
//	require.Nil(t, err)
//	args, err := name.Inputs.Pack("ueni")
//	require.Nil(t, err)
//	res, _, err := p.RunAndCalculateGas(&evm, common.Address{}, common.Address{}, append(p.GetExecutor().(*bank.PrecompileExecutor).NameID, args...), 100000, nil, nil, false, false)
//	require.Nil(t, err)
//	outputs, err := name.Outputs.Unpack(res)
//	require.Nil(t, err)
//	require.Equal(t, "ENI", outputs[0])
//
//	symbol, err := p.ABI.MethodById(p.GetExecutor().(*bank.PrecompileExecutor).SymbolID)
//	require.Nil(t, err)
//	args, err = symbol.Inputs.Pack("ueni")
//	require.Nil(t, err)
//	res, _, err = p.RunAndCalculateGas(&evm, common.Address{}, common.Address{}, append(p.GetExecutor().(*bank.PrecompileExecutor).SymbolID, args...), 100000, nil, nil, false, false)
//	require.Nil(t, err)
//	outputs, err = symbol.Outputs.Unpack(res)
//	require.Nil(t, err)
//	require.Equal(t, "ueni", outputs[0])
//
//	decimal, err := p.ABI.MethodById(p.GetExecutor().(*bank.PrecompileExecutor).DecimalsID)
//	require.Nil(t, err)
//	args, err = decimal.Inputs.Pack("ueni")
//	require.Nil(t, err)
//	res, _, err = p.RunAndCalculateGas(&evm, common.Address{}, common.Address{}, append(p.GetExecutor().(*bank.PrecompileExecutor).DecimalsID, args...), 100000, nil, nil, false, false)
//	require.Nil(t, err)
//	outputs, err = decimal.Outputs.Unpack(res)
//	require.Nil(t, err)
//	require.Equal(t, uint8(0), outputs[0])
//}
//
//func TestAddress(t *testing.T) {
//	k := &testkeeper.EVMTestApp.EvmKeeper
//	p, err := bank.NewPrecompile(k.BankKeeper(), bankkeeper.NewMsgServerImpl(k.BankKeeper()), k, k.AccountKeeper())
//	require.Nil(t, err)
//	require.Equal(t, common.HexToAddress(bank.BankAddress), p.Address())
//}
