package main

import (
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	poolAddress   = "0x..."                                                              // Replace with actual pool address
	token0Address = "0x..."                                                              // Replace with actual token0 address
	token1Address = "0x..."                                                              // Replace with actual token1 address
	privateKey    = "0x57acb95d82739866a5c29e40b0aa2590742ae50425b7dd5b5d279a986370189e" // Replace with actual private key
	rpcURL        = "http://localhost:8545"                                              // Replace with actual RPC URL
)

func main() {
	// Connect to Ethereum client
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
	}

	// Create authorized transactor
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(
		privateKeyECDSA,
		big.NewInt(1337), // Replace with actual chain ID
	)
	if err != nil {
		log.Fatal(err)
	}

	// Load Uniswap V4 pool contract
	pool, err := NewUniswapV4Pool(common.HexToAddress(poolAddress), client)
	if err != nil {
		log.Fatal(err)
	}

	// Load token contracts
	token0, err := NewIERC20(common.HexToAddress(token0Address), client)
	if err != nil {
		log.Fatal(err)
	}

	token1, err := NewIERC20(common.HexToAddress(token1Address), client)
	if err != nil {
		log.Fatal(err)
	}

	// Approve tokens
	amount := big.NewInt(1000000000000000000) // 1 token
	_, err = token0.Approve(auth, common.HexToAddress(poolAddress), amount)
	if err != nil {
		log.Fatal(err)
	}

	_, err = token1.Approve(auth, common.HexToAddress(poolAddress), amount)
	if err != nil {
		log.Fatal(err)
	}

	// Add liquidity
	_, err = pool.Mint(auth, auth.From)
	if err != nil {
		log.Fatal(err)
	}

	// Get initial reserves
	reserve0, reserve1, _, err := pool.GetReserves(nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Initial reserves: %s, %s\n", reserve0.String(), reserve1.String())

	// Perform swaps
	iterations := 100
	startTime := time.Now()

	for i := 0; i < iterations; i++ {
		// Swap token0 for token1
		amount0Out := big.NewInt(0)
		amount1Out := big.NewInt(1000000000000000) // 0.001 token
		_, err = pool.Swap(auth, amount0Out, amount1Out, auth.From)
		if err != nil {
			log.Printf("Swap %d failed: %v\n", i, err)
			continue
		}

		// Get updated reserves
		reserve0, reserve1, _, err = pool.GetReserves(nil)
		if err != nil {
			log.Printf("GetReserves failed: %v\n", err)
			continue
		}

		fmt.Printf("Swap %d completed. Reserves: %s, %s\n", i, reserve0.String(), reserve1.String())
	}

	duration := time.Since(startTime)
	fmt.Printf("Test completed in %v\n", duration)
	fmt.Printf("Average time per transaction: %v\n", duration/time.Duration(iterations))
}

// IERC20 represents the ERC20 token interface
type IERC20 struct {
	Address  common.Address
	Contract *bind.BoundContract
}

// NewIERC20 creates a new instance of IERC20
func NewIERC20(address common.Address, backend bind.ContractBackend) (*IERC20, error) {
	contract, err := bindERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IERC20{Address: address, Contract: contract}, nil
}

// Approve approves the spender to spend tokens
func (token *IERC20) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return token.Contract.Transact(opts, "approve", spender, amount)
}

// bindERC20 binds a generic wrapper to an already deployed contract
func bindERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// ERC20ABI is the input ABI used to generate the binding from
const ERC20ABI = `[{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`
