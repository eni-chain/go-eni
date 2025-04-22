package app

import (
	//"context"
	//"encoding/json"
	tmtypes "github.com/cometbft/cometbft/types"
	protov2 "google.golang.org/protobuf/proto"
	"os"

	//"github.com/cosmos/cosmos-sdk/codec"
	//distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	genutilstypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	//"os"
	"path/filepath"
	"testing"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	crptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	//"github.com/cosmos/cosmos-sdk/simapp"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	//"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	//"github.com/tendermint/tendermint/config"
	//"github.com/cometbft/cometbft/config"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"
	//"github.com/tendermint/tendermint/libs/log"
	//"github.com/cometbft/cometbft/libs/log"
	"cosmossdk.io/log"
	//tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	//dbm "github.com/tendermint/tm-db"
	//dbm "github.com/cometbft/cometbft-db"
	dbm "github.com/cosmos/cosmos-db"

	//minttypes "github.com/eni-chain/go-eni/x/mint/types"
	//servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/testutil/sims"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

const TestContract = "TEST"
const TestUser = "eni1jdppe6fnj2q7hjsepty5crxtrryzhuqsjrj95y"

var DefaultConsensusParams = &tmproto.ConsensusParams{
	Block: &tmproto.BlockParams{
		MaxBytes: 200000,
		MaxGas:   2000000,
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

type TestTx struct {
	msgs []sdk.Msg
}

func (t TestTx) GetMsgsV2() ([]protov2.Message, error) {
	//TODO implement me
	panic("implement me")
}

func NewTestTx(msgs []sdk.Msg) TestTx {
	return TestTx{msgs: msgs}
}

func (t TestTx) GetMsgs() []sdk.Msg {
	return t.msgs
}

func (t TestTx) ValidateBasic() error {
	return nil
}

type TestAppOpts struct {
	useSc bool
}

func (t TestAppOpts) Get(s string) interface{} {
	if s == "chain-id" {
		return "eni-test"
	}
	if s == FlagSCEnable {
		return t.useSc
	}
	return nil
}

type TestWrapper struct {
	suite.Suite

	App *App
	Ctx sdk.Context
}

func NewTestWrapper(t *testing.T, tm time.Time, valPub crptotypes.PubKey, enableEVMCustomPrecompiles bool, baseAppOptions ...func(*baseapp.BaseApp)) *TestWrapper {
	return newTestWrapper(t, tm, valPub, enableEVMCustomPrecompiles, false, baseAppOptions...)
}

func NewTestWrapperWithSc(t *testing.T, tm time.Time, valPub crptotypes.PubKey, enableEVMCustomPrecompiles bool, baseAppOptions ...func(*baseapp.BaseApp)) *TestWrapper {
	return newTestWrapper(t, tm, valPub, enableEVMCustomPrecompiles, true, baseAppOptions...)
}

func newTestWrapper(t *testing.T, tm time.Time, valPub crptotypes.PubKey, enableEVMCustomPrecompiles bool, useSc bool, baseAppOptions ...func(*baseapp.BaseApp)) *TestWrapper {
	var appPtr *App
	if useSc {
		appPtr = SetupWithSc(false, enableEVMCustomPrecompiles, baseAppOptions...)
	} else {
		appPtr = Setup(false, enableEVMCustomPrecompiles, baseAppOptions...)
	}
	ctx := appPtr.BaseApp.NewContext(false)
	ctx.WithBlockHeader(tmproto.Header{Height: 1, ChainID: "eni-test", Time: tm})
	wrapper := &TestWrapper{
		App: appPtr,
		Ctx: ctx,
	}
	wrapper.SetT(t)
	wrapper.setupValidator(stakingtypes.Unbonded, valPub)
	return wrapper
}

func (s *TestWrapper) FundAcc(acc sdk.AccAddress, amounts sdk.Coins) {
	err := s.App.BankKeeper.MintCoins(s.Ctx, minttypes.ModuleName, amounts)
	s.Require().NoError(err)

	err = s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, minttypes.ModuleName, acc, amounts)
	s.Require().NoError(err)
}

func (s *TestWrapper) setupValidator(bondStatus stakingtypes.BondStatus, valPub crptotypes.PubKey) sdk.ValAddress {
	valAddr := sdk.ValAddress(valPub.Address())
	//bondDenom := s.App.StakingKeeper.GetParams(s.Ctx).BondDenom
	params, _ := s.App.StakingKeeper.GetParams(s.Ctx)
	selfBond := sdk.NewCoins(sdk.Coin{Amount: math.NewInt(100), Denom: params.BondDenom})

	s.FundAcc(sdk.AccAddress(valAddr), selfBond)

	//sh := teststaking.NewHelper(s.Suite.T(), s.Ctx, s.App.StakingKeeper)
	//msg := sh.CreateValidatorMsg(valAddr, valPub, selfBond[0].Amount)
	//sh.Handle(msg, true)

	val, err := s.App.StakingKeeper.GetValidator(s.Ctx, valAddr)
	s.Require().True(err == nil, "")

	val = val.UpdateStatus(bondStatus)
	s.App.StakingKeeper.SetValidator(s.Ctx, val)

	consAddr, err := val.GetConsAddr()
	s.Suite.Require().NoError(err)

	signingInfo := slashingtypes.NewValidatorSigningInfo(
		consAddr,
		s.Ctx.BlockHeight(),
		0,
		time.Unix(0, 0),
		false,
		0,
	)
	s.App.SlashingKeeper.SetValidatorSigningInfo(s.Ctx, consAddr, signingInfo)

	return valAddr
}

func (s *TestWrapper) BeginBlock() {
	var proposer sdk.ValAddress

	validators, _ := s.App.StakingKeeper.GetAllValidators(s.Ctx)
	s.Require().Equal(1, len(validators))

	valAddrFancy, err := validators[0].GetConsAddr()
	s.Require().NoError(err)
	proposer = valAddrFancy

	validator, err := s.App.StakingKeeper.GetValidator(s.Ctx, proposer)
	s.Assert().True(err == nil, "")

	valConsAddr, err := validator.GetConsAddr()

	s.Require().NoError(err)

	valAddr := valConsAddr

	newBlockTime := s.Ctx.BlockTime().Add(2 * time.Second)

	header := tmproto.Header{Height: s.Ctx.BlockHeight() + 1, Time: newBlockTime}
	newCtx := s.Ctx.WithBlockTime(newBlockTime).WithBlockHeight(s.Ctx.BlockHeight() + 1)
	s.Ctx = newCtx
	lastCommitInfo := abci.CommitInfo{
		Votes: []abci.VoteInfo{{
			Validator: abci.Validator{Address: valAddr, Power: 1000},
			//SignedLastBlock: true,
			BlockIdFlag: tmproto.BlockIDFlagCommit,
		}},
	}
	//reqBeginBlock := abci.RequestBeginBlock{Header: header, LastCommitInfo: lastCommitInfo}
	//s.App.BeginBlocker(s.Ctx, reqBeginBlock)
	s.Ctx.WithVoteInfos(lastCommitInfo.Votes)

	s.Ctx.WithBlockHeader(header)
	s.App.BeginBlocker(s.Ctx)
}

func (s *TestWrapper) EndBlock() {
	//reqEndBlock := abci.RequestEndBlock{Height: s.Ctx.BlockHeight()}
	//s.App.EndBlocker(s.Ctx, reqEndBlock)
	s.App.EndBlocker(s.Ctx)
}

func Setup(isCheckTx bool, enableEVMCustomPrecompiles bool, baseAppOptions ...func(*baseapp.BaseApp)) (res *App) {
	db := dbm.NewMemDB()
	//encodingConfig := MakeEncodingConfig()
	//cdc := encodingConfig.Marshaler

	//options := []servertypes.AppOptions{
	//	func(app *App) {
	//		app.receiptStore = NewInMemoryStateStore()
	//	},
	//}

	res, _ = New(
		log.NewNopLogger(),
		db,
		nil,
		true,
		DefaultNodeHome,
		sims.EmptyAppOptions{},
		baseAppOptions...,
	)
	baseapp.SetChainID("goeni")(res.BaseApp)
	if !isCheckTx {
		//genesisState := NewDefaultGenesisState(cdc)
		//stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		//if err != nil {
		//	panic(err)
		//}

		dir, _ := os.Getwd()
		genDoc, err := genutilstypes.AppGenesisFromFile(filepath.Join(dir, "../eni-node/config/genesis.json"))
		req := abci.RequestInitChain{
			ChainId:         "goeni",
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: DefaultConsensusParams,
			//AppStateBytes:   stateBytes,
			AppStateBytes: genDoc.AppState,
		}

		_, err = res.InitChain(&req)
		if err != nil {
			panic(err)
		}
	}

	return res
}

func SetupWithSc(isCheckTx bool, enableEVMCustomPrecompiles bool, baseAppOptions ...func(*baseapp.BaseApp)) (res *App) {
	db := dbm.NewMemDB()
	//encodingConfig := MakeEncodingConfig()
	//cdc := encodingConfig.Marshaler

	//options := []AppOption{
	//	func(app *App) {
	//		app.receiptStore = NewInMemoryStateStore()
	//	},
	//}

	res, _ = New(
		log.NewNopLogger(),
		db,
		nil,
		true,
		DefaultNodeHome,
		sims.EmptyAppOptions{},
		baseAppOptions...,
	)
	if !isCheckTx {
		//genesisState := NewDefaultGenesisState(cdc)
		//stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		//if err != nil {
		//	panic(err)
		//}

		// TODO: remove once init chain works with SC
		defer func() { _ = recover() }()

		genDoc, err := genutilstypes.AppGenesisFromFile(filepath.Join(DefaultNodeHome, "config/genesis.json"))
		_, err = res.InitChain(
			&abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: DefaultConsensusParams,
				//AppStateBytes:   stateBytes,
				AppStateBytes: genDoc.AppState,
			},
		)
		if err != nil {
			panic(err)
		}
	}

	return res
}

//func SetupTestingAppWithLevelDb(isCheckTx bool, enableEVMCustomPrecompiles bool) (*App, func()) {
//	dir := "eni_testing"
//	db, err := sdk.NewLevelDB("eni_leveldb_testing", dir)
//	if err != nil {
//		panic(err)
//	}
//	encodingConfig := MakeEncodingConfig()
//	cdc := encodingConfig.Marshaler
//	app := New(
//		log.NewNopLogger(),
//		db,
//		nil,
//		true,
//		map[int64]bool{},
//		DefaultNodeHome,
//		5,
//		enableEVMCustomPrecompiles,
//		nil,
//		encodingConfig,
//		TestAppOpts{},
//		EmptyACLOpts,
//		nil,
//	)
//	if !isCheckTx {
//		genesisState := NewDefaultGenesisState(cdc)
//		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
//		if err != nil {
//			panic(err)
//		}
//
//		_, err = app.InitChain(
//			context.Background(), &abci.RequestInitChain{
//				Validators:      []abci.ValidatorUpdate{},
//				ConsensusParams: DefaultConsensusParams,
//				AppStateBytes:   stateBytes,
//			},
//		)
//		if err != nil {
//			panic(err)
//		}
//	}
//
//	cleanupFn := func() {
//		db.Close()
//		err = os.RemoveAll(dir)
//		if err != nil {
//			panic(err)
//		}
//	}
//
//	return app, cleanupFn
//}
