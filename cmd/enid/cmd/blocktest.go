package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	pruningtypes "cosmossdk.io/store/pruning/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/testutil/sims"
	evmtypes "github.com/cosmos/cosmos-sdk/x/evm/types"
	"github.com/eni-chain/go-eni/app"
	ethtests "github.com/ethereum/go-ethereum/tests"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	//nolint:gosec
	_ "net/http/pprof"
)

//nolint:gosec
func BlocktestCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blocktest",
		Short: "run EF blocktest",
		Long:  "run EF blocktest",
		RunE: func(cmd *cobra.Command, _ []string) error {
			blockTestFileName, err := cmd.Flags().GetString("block-test")
			if err != nil {
				panic(fmt.Sprintf("Error with retrieving block test path: %v", err.Error()))
			}
			testName, err := cmd.Flags().GetString("test-name")
			if err != nil {
				panic(fmt.Sprintf("Error with retrieving test name: %v", err.Error()))
			}
			if blockTestFileName == "" || testName == "" {
				panic("block test file name or test name not set")
			}

			serverCtx := server.GetServerContextFromCmd(cmd)
			if err := serverCtx.Viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}
			home := serverCtx.Viper.GetString(flags.FlagHome)
			db, err := openDB(home)
			if err != nil {
				return err
			}

			logger := log.NewLogger(os.Stdout)
			cache := store.NewCommitKVStoreCacheManager()
			// turn on Cancun for block test
			evmtypes.CancunTime = 0
			a, err := app.New(
				logger,
				db,
				nil,
				true,
				home,
				sims.EmptyAppOptions{},
				baseapp.SetPruning(pruningtypes.NewPruningOptions(pruningtypes.PruningEverything)),
				baseapp.SetMinGasPrices(cast.ToString(serverCtx.Viper.Get(server.FlagMinGasPrices))),
				baseapp.SetMinRetainBlocks(cast.ToUint64(serverCtx.Viper.Get(server.FlagMinRetainBlocks))),
				baseapp.SetInterBlockCache(cache),
				baseapp.SetChainID("goeni"),
			)
			bt := testIngester(blockTestFileName, testName)
			app.BlockTest(a, bt)
			return nil
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The database home directory")
	cmd.Flags().String(flags.FlagChainID, "go-eni", "chain ID")
	cmd.Flags().String("block-test", "", "path to a block test json file")
	cmd.Flags().String("test-name", "", "individual test name")

	return cmd
}

func testIngester(testFilePath string, testName string) *ethtests.BlockTest {
	file, err := os.Open(testFilePath)
	if err != nil {
		panic(err)
	}
	var tests map[string]ethtests.BlockTest
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tests)
	if err != nil {
		panic(err)
	}

	res, ok := tests[testName]
	if !ok {
		panic(fmt.Sprintf("Unable to find test name %v at test file path %v", testName, testFilePath))
	}

	return &res
}
func openDB(rootDir string) (dbm.DB, error) {
	dataDir := filepath.Join(rootDir, "data")

	return dbm.NewGoLevelDB("application", dataDir, nil)
}
