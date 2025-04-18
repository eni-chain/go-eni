package evmrpc_test

import (
	"fmt"
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestAssocation(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	// Sign empty payload prepended with Ethereum Signed Message
	customMessageHash := crypto.Keccak256Hash([]byte("\x19Ethereum Signed Message:\n0"))
	signature, err := crypto.Sign(customMessageHash[:], privateKey)
	if err != nil {
		log.Fatalf("Failed to sign payload: %v", err)
	}

	txArgs := map[string]interface{}{
		"r":              fmt.Sprintf("0x%v", new(big.Int).SetBytes(signature[:32]).Text(16)),
		"s":              fmt.Sprintf("0x%v", new(big.Int).SetBytes(signature[32:64]).Text(16)),
		"v":              fmt.Sprintf("0x%v", new(big.Int).SetBytes([]byte{signature[64]}).Text(16)),
		"custom_message": "\x19Ethereum Signed Message:\n0",
	}

	body := sendRequestGoodWithNamespace(t, "eni", "associate", txArgs)
	require.Equal(t, nil, body["result"])
}

func TestGetEniAddress(t *testing.T) {
	body := sendRequestGoodWithNamespace(t, "eni", "getEniAddress", "0x1df809C639027b465B931BD63Ce71c8E5834D9d6")
	require.Equal(t, "eni1mf0llhmqane5w2y8uynmghmk2w4mh0xlz8e38u", body["result"])
}

func TestGetEvmAddress(t *testing.T) {
	body := sendRequestGoodWithNamespace(t, "eni", "getEVMAddress", "eni1mf0llhmqane5w2y8uynmghmk2w4mh0xlz8e38u")
	require.Equal(t, "0x1df809C639027b465B931BD63Ce71c8E5834D9d6", body["result"])
}

func TestGetCosmosTx(t *testing.T) {
	body := sendRequestGoodWithNamespace(t, "eni", "getCosmosTx", "0xf02362077ac075a397344172496b28e913ce5294879d811bb0269b3be20a872e")
	fmt.Println(body)
	require.Equal(t, "443733A22FFA3D1D7F99769F64FC8D815F6CC13E0036137A0C13E43DEBC587A5", body["result"])
}

func TestGetEvmTx(t *testing.T) {
	body := sendRequestGoodWithNamespace(t, "eni", "getEvmTx", "690D39ADF56D4C811B766DFCD729A415C36C4BFFE80D63E305373B9518EBFB14")
	fmt.Println(body)
	require.Equal(t, "0xf02362077ac075a397344172496b28e913ce5294879d811bb0269b3be20a872e", body["result"])
}
