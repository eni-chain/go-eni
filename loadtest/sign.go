package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	evmcrypto "github.com/ethereum/go-ethereum/crypto"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	cryptokeys "github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/sr25519"
	"github.com/cosmos/cosmos-sdk/client"
	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/ethereum/go-ethereum/common"

	//"github.com/cosmos/cosmos-sdk/crypto/keys/sr25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type AccountInfo struct {
	Address  string `json:"address"`
	Mnemonic string `json:"mnemonic"`
}

type SignerInfo struct {
	AccountNumber  uint64
	SequenceNumber uint64
	mutex          *sync.Mutex
}

func NewSignerInfo(accountNumber uint64, sequenceNumber uint64) *SignerInfo {
	return &SignerInfo{
		AccountNumber:  accountNumber,
		SequenceNumber: sequenceNumber,
		mutex:          &sync.Mutex{},
	}
}

func (si *SignerInfo) IncrementAccountNumber() {
	si.mutex.Lock()
	defer si.mutex.Unlock()
	si.AccountNumber++
}

type SignerClient struct {
	CachedAccountSeqNum *sync.Map
	CachedAccountKey    *sync.Map
	NodeURI             string
}

func NewSignerClient(nodeURI string) *SignerClient {
	return &SignerClient{
		CachedAccountSeqNum: &sync.Map{},
		CachedAccountKey:    &sync.Map{},
		NodeURI:             nodeURI,
	}
}

func (sc *SignerClient) GetTestAccountsKeys(maxAccounts int) []cryptotypes.PrivKey {
	//userHomeDir, _ := os.UserHomeDir()
	// todo check code
	userHomeDir, _ := os.Getwd()
	files, _ := os.ReadDir(filepath.Join(userHomeDir, "loadtest/test_accounts"))
	var testAccountsKeys = make([]cryptotypes.PrivKey, int(math.Min(float64(len(files)), float64(maxAccounts))))
	var wg sync.WaitGroup
	fmt.Printf("Loading accounts\n")
	for i, file := range files {
		if i >= maxAccounts {
			break
		}
		wg.Add(1)
		go func(i int, fileName string) {
			defer wg.Done()
			key := sc.GetKey(fmt.Sprint(i), "test", filepath.Join(userHomeDir, "loadtest/test_accounts", fileName))
			testAccountsKeys[i] = key
		}(i, file.Name())
	}
	wg.Wait()
	fmt.Printf("Finished loading %d accounts \n", len(testAccountsKeys))
	printEvmAddress(testAccountsKeys)
	return testAccountsKeys
}

func printEvmAddress(keys []cryptotypes.PrivKey) {
	for _, key := range keys {
		//Calculate EVM address
		evmAddress := common.Address{}
		if strings.Contains(key.PubKey().Type(), "secp256k1") {
			pubKey := key.PubKey().Bytes()
			if len(pubKey) == 33 {
				pubK, err1 := btcec.ParsePubKey(pubKey)
				if err1 != nil {
					panic(err1)
				}
				pubKey = pubK.SerializeUncompressed()
			}
			hash := evmcrypto.Keccak256(pubKey[1:])
			evmAddress = common.BytesToAddress(hash[len(hash)-20:])
		}
		fmt.Printf("EVM Address: %s Private Key: %x\n", evmAddress.Hex(), key.Bytes())
	}
}

func (sc *SignerClient) GetAdminAccountKeyPath() string {
	userHomeDir, _ := os.UserHomeDir()
	return filepath.Join(userHomeDir, ".eni", "config", "admin_key.json")
}

func (sc *SignerClient) GetAdminKey() cryptotypes.PrivKey {
	return sc.GetKey("admin", "os", sc.GetAdminAccountKeyPath())
}

func (sc *SignerClient) GetKey(accountID, backend, accountKeyFilePath string) cryptotypes.PrivKey {
	if val, ok := sc.CachedAccountKey.Load(accountID); ok {
		privKey := val.(cryptotypes.PrivKey)
		return privKey
	}
	userHomeDir, _ := os.UserHomeDir()
	jsonFile, err := os.Open(accountKeyFilePath)
	if err != nil {
		panic(err)
	}
	var accountInfo AccountInfo
	byteVal, err := io.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}
	jsonFile.Close()
	if err := json.Unmarshal(byteVal, &accountInfo); err != nil {
		panic(err)
	}

	encodingConfig := moduletestutil.MakeTestEncodingConfig()
	kr, _ := keyring.New(sdk.KeyringServiceName(), backend, filepath.Join(userHomeDir, ".eni"), os.Stdin, encodingConfig.Codec)
	keyringAlgos, _ := kr.SupportedAlgorithms()
	algoStr := string(hd.Secp256k1Type)
	algo, _ := keyring.NewSigningAlgoFromString(algoStr, keyringAlgos)
	hdpath := hd.CreateHDPath(sdk.GetConfig().GetCoinType(), 0, 0).String()
	derivedPriv, _ := algo.Derive()(accountInfo.Mnemonic, "", hdpath)
	privKey := algo.Generate()(derivedPriv)

	// Cache this so we don't need to regenerate it
	sc.CachedAccountKey.Store(accountID, privKey)
	return privKey
}

func (sc *SignerClient) GetValKeys() []cryptokeys.PrivKey {
	valKeys := []cryptokeys.PrivKey{}
	userHomeDir, _ := os.UserHomeDir()
	valKeysFilePath := filepath.Join(userHomeDir, "exported_keys")
	files, _ := os.ReadDir(valKeysFilePath)
	for _, fn := range files {
		// we dont expect subdirectories, so we can just handle files
		valKeyFile := filepath.Join(valKeysFilePath, fn.Name())
		privKeyBz, err := os.ReadFile(valKeyFile)
		if err != nil {
			panic(err)
		}

		privKeyBytes, algo, err := crypto.UnarmorDecryptPrivKey(string(privKeyBz), "12345678")
		if err != nil {
			panic(err)
		}

		var privKey cryptokeys.PrivKey
		var ok bool
		if algo == string(hd.Sr25519Type) {
			typedKey := &sr25519.PrivKey{}
			if err := typedKey.UnmarshalJSON(privKeyBytes.Bytes()); err != nil {
				panic(err)
			}
			privKey = typedKey
		} else {
			key, err := legacy.PrivKeyFromBytes(privKeyBytes.Bytes())
			if err != nil {
				panic(err)
			}

			if privKey, ok = key.(cryptokeys.PrivKey); !ok {
				panic("invalid private key type")
			}
		}

		valKeys = append(valKeys, privKey)
	}
	return valKeys
}

func (sc *SignerClient) SignTx(chainID string, txBuilder *client.TxBuilder, privKey cryptotypes.PrivKey, seqDelta uint64) {
	var sigsV2 []signing.SignatureV2
	signerInfo := sc.GetAccountNumberSequenceNumber(privKey)
	accountNum := signerInfo.AccountNumber
	seqNum := signerInfo.SequenceNumber

	seqNum += seqDelta
	sigV2 := signing.SignatureV2{
		PubKey: privKey.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  signing.SignMode(TestConfig.TxConfig.SignModeHandler().DefaultMode()),
			Signature: nil,
		},
		Sequence: seqNum,
	}
	sigsV2 = append(sigsV2, sigV2)
	_ = (*txBuilder).SetSignatures(sigsV2...)
	sigsV2 = []signing.SignatureV2{}
	signerData := xauthsigning.SignerData{
		ChainID:       chainID,
		AccountNumber: accountNum,
		Sequence:      seqNum,
	}
	sigV2, _ = clienttx.SignWithPrivKey(
		context.TODO(),
		signing.SignMode(TestConfig.TxConfig.SignModeHandler().DefaultMode()),
		signerData,
		*txBuilder,
		privKey,
		TestConfig.TxConfig,
		seqNum,
	)
	sigsV2 = append(sigsV2, sigV2)
	_ = (*txBuilder).SetSignatures(sigsV2...)
}

func (sc *SignerClient) GetAccountNumberSequenceNumber(privKey cryptotypes.PrivKey) SignerInfo {
	if val, ok := sc.CachedAccountSeqNum.Load(privKey); ok {
		signerinfo := val.(SignerInfo)
		return signerinfo
	}

	hexAccount := privKey.PubKey().Address()
	address, err := AccAddressFromHex(hexAccount.String())
	if err != nil {
		panic(err)
	}
	accountRetriever := authtypes.AccountRetriever{}
	cl, err := client.NewClientFromNode(sc.NodeURI)
	if err != nil {
		panic(err)
	}
	context := client.Context{}
	context = context.WithNodeURI(sc.NodeURI)
	context = context.WithClient(cl)
	context = context.WithInterfaceRegistry(TestConfig.InterfaceRegistry)
	userHomeDir, _ := os.UserHomeDir()
	encodingConfig := moduletestutil.MakeTestEncodingConfig()

	kr, _ := keyring.New(sdk.KeyringServiceName(), "test", filepath.Join(userHomeDir, ".eni"), os.Stdin, encodingConfig.Codec)
	context = context.WithKeyring(kr)
	account, seq, err := accountRetriever.GetAccountNumberSequence(context, address)
	if err != nil {
		time.Sleep(5 * time.Second)
		// retry once after 5 seconds
		account, seq, err = accountRetriever.GetAccountNumberSequence(context, address)
		if err != nil {
			panic(err)
		}
	}

	signerInfo := *NewSignerInfo(account, seq)
	sc.CachedAccountSeqNum.Store(privKey, signerInfo)
	return signerInfo
}

func addressBytesFromHexString(address string) ([]byte, error) {
	if len(address) == 0 {
		return nil, errors.New("decoding Bech32 address failed: must provide an address")
	}

	return hex.DecodeString(address)
}

type AccAddress []byte

// AccAddressFromHex creates an AccAddress from a hex string.
func AccAddressFromHex(address string) (addr sdk.AccAddress, err error) {
	bz, err := addressBytesFromHexString(address)
	return bz, err
}
