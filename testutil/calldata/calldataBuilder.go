package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"golang.org/x/crypto/ed25519"
	"io/ioutil"
	"math/big"
	"math/rand"
	"strings"
	"time"
)

func calldataF(caller string, name, method string, args ...interface{}) (string, error) {
	//system contract dir
	path := "/Users/moses/workspace/go-eni/syscontract/genesis/dpos/"
	abiJson, err := ioutil.ReadFile(path + name + ".abi")
	if err != nil {
		fmt.Printf("Read ABI file failed, err:%v", err.Error())
	}

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	if err != nil {
		fmt.Printf("constrcut ABI obj failed, err:%v", err.Error())
	}

	dataByte, err := myAbi.Pack(method, args...)
	if err != nil {
		fmt.Printf("create ABI data failed, err:%v", err.Error())
	}

	dataString := hex.EncodeToString(dataByte)
	fmt.Printf("\nFor:%s, contract:%s, method:%s, calldata:\n%s\n", caller, name, method, dataString)

	return dataString, err
}

func base64Decode(dataString string) (string, error) {
	base64Bytes, err := base64.StdEncoding.DecodeString(dataString)
	if err != nil {
		return "", err
	}
	hexString := hex.EncodeToString(base64Bytes)
	fmt.Printf("\nbase64 to hex:\n%s\n", hexString)
	return string(base64Bytes), nil
}

func getRand() []byte {
	rand.Seed(time.Now().UnixNano())

	bytes := make([]byte, 32)
	for i := 0; i < 32; i++ {
		bytes[i] = byte(rand.Intn(255))
	}
	return bytes
}

const AddressLength = 20

type Address [AddressLength]byte

func getAddrByPriKey(priBytes []byte) (Address, ed25519.PublicKey) {
	priKey := ed25519.PrivateKey(priBytes)
	pubKey := priKey.Public().(ed25519.PublicKey)

	hash := sha256.Sum256(pubKey)
	var addr Address
	//addr := hash[:20]
	copy(addr[:], hash[:20])
	return addr, pubKey
}

func getAddrByHexPriKey(hexPriKey string) (Address, ed25519.PublicKey, ed25519.PrivateKey) {
	priKey, err := hex.DecodeString(hexPriKey)
	if err != nil {
		fmt.Printf("hex decode failed, err:%v", err.Error())
	}
	addr, pubKey := getAddrByPriKey(priKey)
	return addr, pubKey, priKey
}

// get address and public key by private key encode by base64
func getAddrByBase64PriKey(privateKey string) (Address, ed25519.PublicKey, ed25519.PrivateKey) {
	priBytes, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		fmt.Printf("base64 decode failed, err:%v", err.Error())
	}

	addr, pubKey := getAddrByPriKey(priBytes)
	return addr, pubKey, priBytes
}

// get standard address by hex str address
func getAddrByHexStr(str string) Address {
	addrBytes, err := hex.DecodeString(str)
	if err != nil {
		fmt.Printf("hex decode failed, err:%v", err.Error())
	}

	var addr Address
	copy(addr[:], addrBytes[:20])
	return addr
}

func getAddrAndSecKeys() (addr Address, pubKey ed25519.PublicKey, priKey ed25519.PrivateKey) {
	priKey = ed25519.NewKeyFromSeed(getRand())
	addr, pubKey = getAddrByPriKey(priKey)

	//fmt.Printf("addr1:%x\npubKey1:%x\npriKey1:%x\n", addr, pubKey, priKey)
	return addr, pubKey, priKey
}

//alice
//addr：0xD5dd01FFCC547734bCe4Df0F5acf4A5407c275d3
//priKey： dc9bb398d00f7778a61dcbb7e90cfe527b7e7b69ce9d557a08d5e32ea8d3eac0

//bob
//addr：0x89b7145BeF43EfAfaf076C181b124d22f9D218e5
//priKey：d731aab93f2ba8503c351bf67eea13277cf2c35075ff6d9ae1bfa23d40c26501

//clare
//addr：0xF87A299e6bC7bEba58dbBe5a5Aa21d49bCD16D52
//priKey：0x57acb95d82739866a5c29e40b0aa2590742ae50425b7dd5b5d279a986370189e

func main() {
	//get default validator addr, pubKey, and priKey, bond alice
	alice := "251604eBfD1ddeef1F4f40b8F9Fc425538BE1339"
	valAddr, valPubKey, valPriKey := getAddrByBase64PriKey("ooX0ThgTQSWWrH+gVy1w1esHCbLmi9+FyPWPn140F9iIujcEkgvQ9s4XJyoq99AHfqn3DCvRNCO8auNrpn0AEQ==")
	fmt.Printf("\naddr:%x\npubKey:%x\npriKey:%x\n", valAddr, valPubKey, valPriKey)

	////generate ed25519 secret key and addr
	//addr1, pubKey1, priKey1 := getAddrAndSecKeys()
	//fmt.Printf("addr1:%x\npubKey1:%x\npriKey1:%x\n", addr1, pubKey1, priKey1)

	//bond clare
	clare := "F87A299e6bC7bEba58dbBe5a5Aa21d49bCD16D52"
	addr2, pubKey2, priKey2 := getAddrByHexPriKey("b7cef85c61f7a973896d5b12f7de5020453dde51c19e4693fb6a55dfa8e3e64080f123e970b41abe1709ff176fc6bcd673afd41456e064e337f8713f8bcd068e")
	fmt.Printf("\naddr2:%x\npubKey2:%x\npriKey:%x\n", addr2, pubKey2, priKey2)

	//calldata for apply for validator
	calldataF(alice, "hub", "applyForValidator", valAddr, valAddr, "node1", "node1 apply for validator", []byte(valPubKey))
	calldataF(clare, "hub", "applyForValidator", addr2, addr2, "node2", "node2 apply for validator", []byte(pubKey2))

	//admin(alice) audit
	calldataF(alice, "hub", "auditPass", getAddrByHexStr(alice)) //alice
	calldataF(clare, "hub", "auditPass", getAddrByHexStr(clare)) //clare

	//admin(alice) set init seed
	hexInitSeed := "ba1aa46438a7b446c0a6f1ca54d04eccda80fed5f1460be9e17cd6931eaef64c1f1cbe714c603521c2f06a4a39cd8d50015068890aaaf04d92d9ed997f9c0689"
	initSeed, _ := hex.DecodeString(hexInitSeed)
	calldataF(alice, "vrf", "init", initSeed, big.NewInt(1))

	//alice send random
	valSignature := ed25519.Sign(valPriKey, initSeed)
	calldataF(alice, "vrf", "sendRandom", valSignature, big.NewInt(2))

	//clare send random
	valSignature2 := ed25519.Sign(priKey2, initSeed)
	calldataF(clare, "vrf", "sendRandom", valSignature2, big.NewInt(2))
}
