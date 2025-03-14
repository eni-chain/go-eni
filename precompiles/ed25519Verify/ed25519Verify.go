package ed25519Verify

import (
	"errors"
	cryptoencoding "github.com/cometbft/cometbft/crypto/encoding"
	crypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
)

const (
	EVM_WORD_LEN       = 32
	ED25519_PUBKEY_LEN = 32
	ED25519_SIGN_LEN   = 64
	ED25519_HASH_LEN   = 64
)

func init() {
	//todo: need to confirm the addr, and also modify addr in syscontract
	addr := common.BytesToAddress([]byte{0xa1})
	precompiled := ed25519Verify{}
	addEd25519VerifyToVM(addr, &precompiled)
}

type ed25519Verify struct{}

func (c *ed25519Verify) RequiredGas(input []byte) uint64 {
	return params.EcrecoverGas
}

func (c *ed25519Verify) Run(_ *vm.EVM, _ common.Address, _ common.Address, input []byte, _ *big.Int, _ bool, _ bool) ([]byte, error) {
	// | PubKey   | Signature  |  msgHash   |
	// | 32 bytes | 64 bytes   |  64 bytes  |
	if len(input) != (ED25519_PUBKEY_LEN + ED25519_SIGN_LEN + ED25519_HASH_LEN) {
		return nil, errors.New("invalid input")
	}
	pkBytes := input[:ED25519_PUBKEY_LEN]
	sigBytes := input[ED25519_PUBKEY_LEN : ED25519_PUBKEY_LEN+ED25519_SIGN_LEN]
	hashBytes := input[ED25519_PUBKEY_LEN+ED25519_SIGN_LEN:]

	pubKey := crypto.PublicKey{}
	err := pubKey.Unmarshal(pkBytes)
	if err != nil {
		return nil, err
	}

	pk, err := cryptoencoding.PubKeyFromProto(pubKey)
	if err != nil {
		return nil, err
	}

	res := make([]byte, EVM_WORD_LEN)
	ret := pk.VerifySignature(hashBytes, sigBytes)
	if !ret {
		return res, nil
	}

	res[EVM_WORD_LEN-1] = 1
	return res, nil
}

func addEd25519VerifyToVM(addr common.Address, p vm.PrecompiledContract) {
	vm.PrecompiledContractsHomestead[addr] = p
	vm.PrecompiledContractsByzantium[addr] = p
	vm.PrecompiledContractsIstanbul[addr] = p
	vm.PrecompiledContractsBerlin[addr] = p
	vm.PrecompiledContractsCancun[addr] = p
	vm.PrecompiledContractsBLS[addr] = p
	vm.PrecompiledAddressesHomestead = append(vm.PrecompiledAddressesHomestead, addr)
	vm.PrecompiledAddressesByzantium = append(vm.PrecompiledAddressesByzantium, addr)
	vm.PrecompiledAddressesIstanbul = append(vm.PrecompiledAddressesIstanbul, addr)
	vm.PrecompiledAddressesBerlin = append(vm.PrecompiledAddressesBerlin, addr)
	vm.PrecompiledAddressesCancun = append(vm.PrecompiledAddressesCancun, addr)
}
