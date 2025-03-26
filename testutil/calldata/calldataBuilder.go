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

func solo() {
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

	//test verifyEd25519Sign method
	calldataF(clare, "vrf", "verifyEd25519Sign", []byte(pubKey2), valSignature2, initSeed)
}

func multi() {
	////generate ed25519 secret key and addr
	//addr, pubKey, priKey := getAddrAndSecKeys()
	//fmt.Printf("addr:%x\npubKey:%x\npriKey:%x\n", addr, pubKey, priKey)

	//get default validator1 addr, pubKey, and priKey
	validator1 := "3140aedbf686A3150060Cb946893b0598b266f5C"
	v1Addr, v1PubKey, v1PriKey := getAddrByBase64PriKey("OVd/AbN5k5xJxE6XFfQwnzvKZw6BDOYPZfuyHmsBcBsM0r5RCbInFMqmNESJn3xZh4S2S18BChhmijAvT9b7Rg==")
	fmt.Printf("\naddr1:%x\npubKey1:%x\npriKey1:%x\n", v1Addr, v1PubKey, v1PriKey)

	validator2 := "B680152a597c937164941e77cbF4a6b2F866675c"
	v2Addr, v2PubKey, v2PriKey := getAddrByBase64PriKey("jo+5N22PB1zAUuW3110Bq8VJlfPpB92innuusH+6AR5wZe3gwp5jWWn1TQ5XTfbr6oV5dNrahgJh52NWgMrCug==")
	fmt.Printf("\naddr2:%x\npubKey2:%x\npriKey2:%x\n", v2Addr, v2PubKey, v2PriKey)

	validator3 := "5746036C781851B9eF234219e770E1591104561f"
	v3Addr, v3PubKey, v3PriKey := getAddrByBase64PriKey("Yqi/ZY6fxJ8udJCYdpZVa0VNDA/eSKVbFw9jeXRTtrkU5so5K4o+2hTwqero4AmjhHssugtGnmRVvPKPOgy5pw==")
	fmt.Printf("\naddr3:%x\npubKey3:%x\npriKey3:%x\n", v3Addr, v3PubKey, v3PriKey)

	validator4 := "7B6ba0Fe2610BF3c69Fc571C37aeC7bB87A281D2"
	v4Addr, v4PubKey, v4PriKey := getAddrByBase64PriKey("nBFy2Dyurr7ccYnuqkJz16082sH+jUJxAleG/i209myhEGVNVrSQ5j37paTe3H/wiawk+PAaR99TOmREQSYCaA==")
	fmt.Printf("\naddr4:%x\npubKey4:%x\npriKey4:%x\n", v4Addr, v4PubKey, v4PriKey)

	//calldata for apply for validator
	calldataF(validator1, "hub", "applyForValidator", v1Addr, v1Addr, "node1", "node1 apply for validator", []byte(v1PubKey))
	calldataF(validator2, "hub", "applyForValidator", v2Addr, v2Addr, "node2", "node2 apply for validator", []byte(v2PubKey))
	calldataF(validator3, "hub", "applyForValidator", v3Addr, v3Addr, "node3", "node3 apply for validator", []byte(v3PubKey))
	calldataF(validator4, "hub", "applyForValidator", v4Addr, v4Addr, "node4", "node4 apply for validator", []byte(v4PubKey))

	//admin(validator1) audit
	calldataF(validator1, "hub", "auditPass", getAddrByHexStr(validator1))
	calldataF(validator2, "hub", "auditPass", getAddrByHexStr(validator2))
	calldataF(validator3, "hub", "auditPass", getAddrByHexStr(validator3))
	calldataF(validator4, "hub", "auditPass", getAddrByHexStr(validator4))

	//admin(validator1) set init seed
	hexInitSeed := "ba1aa46438a7b446c0a6f1ca54d04eccda80fed5f1460be9e17cd6931eaef64c1f1cbe714c603521c2f06a4a39cd8d50015068890aaaf04d92d9ed997f9c0689"
	initSeed, _ := hex.DecodeString(hexInitSeed)
	calldataF(validator1, "vrf", "init", initSeed, big.NewInt(1))

	//validator1 send random
	valSignature1 := ed25519.Sign(v1PriKey, initSeed)
	calldataF(validator1, "vrf", "sendRandom", valSignature1, big.NewInt(2))

	//validator2 send random
	valSignature2 := ed25519.Sign(v2PriKey, initSeed)
	calldataF(validator2, "vrf", "sendRandom", valSignature2, big.NewInt(2))

	//validator3 send random
	valSignature3 := ed25519.Sign(v3PriKey, initSeed)
	calldataF(validator3, "vrf", "sendRandom", valSignature3, big.NewInt(2))

	//validator4 send random
	valSignature4 := ed25519.Sign(v4PriKey, initSeed)
	calldataF(validator4, "vrf", "sendRandom", valSignature4, big.NewInt(2))

	////test verifyEd25519Sign method
	//calldataF(validator1, "vrf", "verifyEd25519Sign", []byte(v1PubKey), valSignature1, initSeed)
}

func main() {
	//solo()
	multi()

}
