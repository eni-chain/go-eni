package app_test

import (
	"context"
	"encoding/hex"
	"math/big"
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/testutil"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/utils/tracing"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/eni-chain/go-eni/app"
	"go.opentelemetry.io/otel"

	abci "github.com/cometbft/cometbft/abci/types"
	testkeeper "github.com/cosmos/cosmos-sdk/testutil/keeper"
	testutilmod "github.com/cosmos/cosmos-sdk/types/module/testutil"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/cosmos-sdk/x/evm/types"
	evmtypes "github.com/cosmos/cosmos-sdk/x/evm/types"
	"github.com/cosmos/cosmos-sdk/x/evm/types/ethtx"
	"github.com/eni-chain/go-eni/app/apptesting"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// AnteTestSuite is a test suite to be used with ante handler tests.
type AnteTestSuite struct {
	suite.Suite

	apptesting.KeeperTestHelper

	encCfg      testutilmod.TestEncodingConfig
	anteHandler sdk.AnteHandler
	clientCtx   client.Context
	txBuilder   client.TxBuilder
	testAcc     sdk.AccAddress
	testAccPriv cryptotypes.PrivKey
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(AnteTestSuite))
}

// SetupTest setups a new test, with new app, context, and anteHandler.
func (suite *AnteTestSuite) SetupTest(isCheckTx bool) {
	suite.Setup()

	// keys and addresses
	suite.testAccPriv, _, suite.testAcc = testdata.KeyTestPubAddr()
	initalBalance := sdk.Coins{sdk.NewInt64Coin("atom", 100000000000)}
	suite.FundAcc(suite.testAcc, initalBalance)

	suite.Ctx = suite.Ctx.WithBlockHeight(1)

	// We're using TestMsg encoding in some tests, so register it here.
	suite.encCfg = testutilmod.MakeTestEncodingConfig(staking.AppModuleBasic{})

	suite.encCfg.Amino.RegisterConcrete(&testdata.TestMsg{}, "testdata.TestMsg", nil)
	testdata.RegisterInterfaces(suite.encCfg.InterfaceRegistry)

	suite.clientCtx = client.Context{}.
		WithTxConfig(suite.encCfg.TxConfig)

	defaultTracer, _ := tracing.DefaultTracerProvider()
	otel.SetTracerProvider(defaultTracer)
	tr := defaultTracer.Tracer("component-main")

	tracingInfo := &tracing.Info{
		Tracer: &tr,
	}
	tracingInfo.SetContext(context.Background())
	anteHandler, err := app.NewAnteHandlerAndDepGenerator(app.HandlerOptions{
		HandlerOptions: ante.HandlerOptions{
			AccountKeeper:   suite.App.AccountKeeper,
			BankKeeper:      suite.App.BankKeeper,
			FeegrantKeeper:  suite.App.FeeGrantKeeper,
			SignModeHandler: suite.clientCtx.TxConfig.SignModeHandler(),
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
		},
	})

	suite.Require().NoError(err)
	suite.anteHandler = anteHandler
}

// CreateTestTx is a helper function to create a tx given multiple inputs.
func (suite *AnteTestSuite) CreateTestTx(privs []cryptotypes.PrivKey, accNums []uint64, accSeqs []uint64, chainID string) (xauthsigning.Tx, error) {
	// First round: we gather all the signer infos. We use the "set empty
	// signature" hack to do that.
	var sigsV2 []signing.SignatureV2
	for i, priv := range privs {
		sigV2 := signing.SignatureV2{
			PubKey: priv.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode:  signing.SignMode_SIGN_MODE_DIRECT,
				Signature: nil,
			},
			Sequence: accSeqs[i],
		}

		sigsV2 = append(sigsV2, sigV2)
	}
	err := suite.txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return nil, err
	}

	// Second round: all signer infos are set, so each signer can sign.
	sigsV2 = []signing.SignatureV2{}
	for i, priv := range privs {
		signerData := xauthsigning.SignerData{
			ChainID:       chainID,
			AccountNumber: accNums[i],
			Sequence:      accSeqs[i],
		}
		sigV2, err := tx.SignWithPrivKey(context.Background(), signing.SignMode_SIGN_MODE_DIRECT, signerData,
			suite.txBuilder, priv, suite.clientCtx.TxConfig, accSeqs[i])
		if err != nil {
			return nil, err
		}

		sigsV2 = append(sigsV2, sigV2)
	}
	err = suite.txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return nil, err
	}

	return suite.txBuilder.GetTx(), nil
}

func TestEvmAnteErrorHandler(t *testing.T) {

	app, ctx := testkeeper.NewMockApp(t, false)

	privKey := testkeeper.MockPrivateKey()
	testPrivHex := hex.EncodeToString(privKey.Bytes())
	key, _ := crypto.HexToECDSA(testPrivHex)
	txData := ethtypes.LegacyTx{
		GasPrice: big.NewInt(1000000000000),
		Gas:      200000,
		To:       nil,
		Value:    big.NewInt(0),
		Data:     []byte{},
		Nonce:    1, // will cause ante error
	}
	chainID := app.EVMKeeper.ChainID(ctx)
	chainCfg := evmtypes.DefaultChainConfig()
	ethCfg := chainCfg.EthereumConfig(chainID)
	blockNum := big.NewInt(ctx.BlockHeight())
	signer := ethtypes.MakeSigner(ethCfg, blockNum, uint64(ctx.BlockTime().Unix()))
	tx, err := ethtypes.SignTx(ethtypes.NewTx(&txData), signer, key)
	require.Nil(t, err)
	txwrapper, err := ethtx.NewLegacyTx(tx)
	require.Nil(t, err)
	req, err := types.NewMsgEVMTransaction(txwrapper)
	require.Nil(t, err)
	interfaceRegistry := testutil.CodecOptions{}.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)
	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)
	builder := txConfig.NewTxBuilder()
	builder.SetMsgs(req)
	txToSend := builder.GetTx()
	encodedTx, err := app.App.TxEncode(txToSend)
	require.Nil(t, err)
	_, err = app.App.TxDecode(encodedTx)
	require.Nil(t, err)

	addr, _ := testkeeper.PrivateKeyToAddresses(privKey)
	app.GetEVMKeeper().BankKeeper().AddCoins(ctx, addr, sdk.NewCoins(sdk.NewCoin("ueni", math.NewInt(100000000000))))

	block, err := app.App.FinalizeBlock(&abci.RequestFinalizeBlock{
		Txs:    [][]byte{encodedTx},
		Height: 1,
	})
	require.NoError(t, err)

	msg := convertEVMMsg(txToSend)
	app.GetEVMKeeper().SetTxResults(block.TxResults)
	app.GetEVMKeeper().SetMsgs([]*types.MsgEVMTransaction{msg})

	deferredInfo := app.GetEVMKeeper().GetAllEVMTxDeferredInfo(ctx)
	require.Equal(t, 1, len(deferredInfo))
}

func convertEVMMsg(tx sdk.Tx) (res *types.MsgEVMTransaction) {
	defer func() {
		if err := recover(); err != nil {
			res = nil
		}
	}()
	if tx == nil {
		return nil
	} else if emsg := types.GetEVMTransactionMessage(tx); emsg != nil && !emsg.IsAssociateTx() {
		return emsg
	} else {
		return nil
	}
}
