package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	uniswap "github.com/liuyunlong/go-eni/loadtest/contracts/uniswap_v4/bindings"
)

func loadEnvFile() error {
	// Try to load .env file from multiple locations
	envFiles := []string{
		".env",                                   // Current directory
		"../.env",                                // Parent directory
		"../../.env",                             // Two levels up
		"../../../.env",                          // Three levels up
		filepath.Join(os.Getenv("HOME"), ".env"), // Home directory
	}

	for _, envFile := range envFiles {
		if err := godotenv.Load(envFile); err == nil {
			fmt.Printf("Loaded environment from %s\n", envFile)
			return nil
		}
	}

	return fmt.Errorf("could not find .env file in any of the following locations: %v", envFiles)
}

func main() {
	// Load .env file
	if err := loadEnvFile(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Read configuration from .env file
	poolAddress := os.Getenv("POOL_ADDRESS")
	token0Address := os.Getenv("TOKEN0_ADDRESS")
	token1Address := os.Getenv("TOKEN1_ADDRESS")
	privateKey := os.Getenv("PRIVATE_KEY")
	rpcURL := os.Getenv("RPC_URL")

	// Validate required environment variables
	if poolAddress == "" || token0Address == "" || token1Address == "" || privateKey == "" || rpcURL == "" {
		log.Fatal("Missing required environment variables in .env file")
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
	}

	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKeyECDSA, chainID)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize pool contract
	pool, err := uniswap.NewUniswapV4Pool(common.HexToAddress(poolAddress), client)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize token contracts
	token0, err := uniswap.NewIERC20(common.HexToAddress(token0Address), client)
	if err != nil {
		log.Fatal(err)
	}

	token1, err := uniswap.NewIERC20(common.HexToAddress(token1Address), client)
	if err != nil {
		log.Fatal(err)
	}

	// Get account balance
	balance, err := client.BalanceAt(context.Background(), auth.From, nil)
	if err != nil {
		log.Fatal("Failed to get account balance:", err)
	}
	fmt.Printf("Account balance: %s ETH\n", balance.String())

	// Get token balances
	token0Balance, err := token0.BalanceOf(nil, auth.From)
	if err != nil {
		log.Fatal("Failed to get token0 balance:", err)
	}
	fmt.Printf("Token0 balance: %s\n", token0Balance.String())

	token1Balance, err := token1.BalanceOf(nil, auth.From)
	if err != nil {
		log.Fatal("Failed to get token1 balance:", err)
	}
	fmt.Printf("Token1 balance: %s\n", token1Balance.String())

	// Get current allowances
	allowance0, err := token0.Allowance(nil, auth.From, common.HexToAddress(poolAddress))
	if err != nil {
		log.Fatal("Failed to get token0 allowance:", err)
	}
	fmt.Printf("Token0 allowance: %s\n", allowance0.String())

	allowance1, err := token1.Allowance(nil, auth.From, common.HexToAddress(poolAddress))
	if err != nil {
		log.Fatal("Failed to get token1 allowance:", err)
	}
	fmt.Printf("Token1 allowance: %s\n", allowance1.String())

	// Approve tokens for the pool if needed
	approveAmount := big.NewInt(0).Mul(big.NewInt(1000000), big.NewInt(1e18)) // 1 million tokens
	if allowance0.Cmp(approveAmount) < 0 {
		fmt.Println("Approving token0...")
		auth.GasLimit = uint64(300000)
		auth.GasPrice, err = client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatal("Failed to get gas price:", err)
		}
		tx, err := token0.Approve(auth, common.HexToAddress(poolAddress), approveAmount)
		if err != nil {
			log.Fatal("Failed to approve token0:", err)
		}
		receipt, err := bind.WaitMined(context.Background(), client, tx)
		if err != nil {
			log.Fatal("Failed to wait for token0 approval:", err)
		}
		if receipt.Status == 0 {
			log.Fatal("Token0 approval failed")
		}
		fmt.Println("Token0 approved successfully")
	}

	if allowance1.Cmp(approveAmount) < 0 {
		fmt.Println("Approving token1...")
		auth.GasLimit = uint64(300000)
		auth.GasPrice, err = client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatal("Failed to get gas price:", err)
		}
		tx, err := token1.Approve(auth, common.HexToAddress(poolAddress), approveAmount)
		if err != nil {
			log.Fatal("Failed to approve token1:", err)
		}
		receipt, err := bind.WaitMined(context.Background(), client, tx)
		if err != nil {
			log.Fatal("Failed to wait for token1 approval:", err)
		}
		if receipt.Status == 0 {
			log.Fatal("Token1 approval failed")
		}
		fmt.Println("Token1 approved successfully")
	}

	iterations := 10 // Number of swaps
	for i := 0; i < iterations; i++ {
		// Get the current nonce for the account
		nonce, err := client.PendingNonceAt(context.Background(), auth.From)
		if err != nil {
			log.Fatalf("Failed to get nonce: %v", err)
		}

		auth.Nonce = big.NewInt(int64(nonce))
		auth.GasLimit = uint64(300000)
		auth.GasPrice, err = client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatalf("Failed to get gas price: %v", err)
		}

		// Set swap parameters
		amount0In := big.NewInt(1000000000000000) // 0.001 token
		amount1In := big.NewInt(0)
		amount0Out := big.NewInt(0)
		amount1Out := big.NewInt(990000000000000) // 0.00099 token (0.1% fee)

		// Create the swap transaction
		tx, err := pool.Swap(auth, amount0In, amount1In, amount0Out, amount1Out, auth.From)
		if err != nil {
			log.Printf("Swap %d transaction failed to send: %v\n", i, err)
			continue
		}

		fmt.Printf("Swap %d transaction sent, tx hash: %s. Waiting for confirmation...\n", i, tx.Hash().Hex())

		// Wait for the transaction to be mined
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		receipt, err := bind.WaitMined(ctx, client, tx)
		if err != nil {
			log.Printf("Swap %d transaction mining failed or timed out: %v\n", i, err)
			continue
		}

		if receipt.Status == 1 {
			fmt.Printf("Swap %d transaction confirmed successfully, block number: %d\n", i, receipt.BlockNumber.Uint64())
		} else {
			log.Printf("Swap %d transaction failed, block number: %d\n", i, receipt.BlockNumber.Uint64())
		}

		// Add a small delay between transactions
		time.Sleep(1 * time.Second)
	}

	fmt.Println("All transactions completed")
}
