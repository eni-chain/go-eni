package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"time"

	uniswap "github.com/eni-chain/go-eni/loadtest/contracts/uniswap_v4/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func loadEnvFile() error {
	// Try to load .env file from multiple locations
	envFiles := []string{
		".env",                                   // Current directory
		"loadtest/contracts/uniswap_v4/.env",     // Parent directory
		"contracts/uniswap_v4/.env",              // Two levels up
		"uniswap_v4/.env",                        // Three levels up
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

	// Get pool reserves
	reserves, err := pool.GetReserves(nil)
	if err != nil {
		log.Fatal("Failed to get pool reserves:", err)
	}
	fmt.Printf("Pool reserves - Token0: %s, Token1: %s\n", reserves.Reserve0.String(), reserves.Reserve1.String())

	// Check if pool has liquidity
	if reserves.Reserve0.Cmp(big.NewInt(0)) == 0 || reserves.Reserve1.Cmp(big.NewInt(0)) == 0 {
		fmt.Println("Pool has no liquidity, adding initial liquidity...")

		// Calculate initial liquidity amounts
		initialAmount0 := big.NewInt(1000000000000000000) // 1 token
		initialAmount1 := big.NewInt(1000000000000000000) // 1 token

		// Approve tokens for minting if needed
		if allowance0.Cmp(initialAmount0) < 0 {
			fmt.Println("Approving token0 for minting...")
			auth.GasLimit = uint64(300000)
			auth.GasPrice, err = client.SuggestGasPrice(context.Background())
			if err != nil {
				log.Fatal("Failed to get gas price:", err)
			}
			tx, err := token0.Approve(auth, common.HexToAddress(poolAddress), initialAmount0)
			if err != nil {
				log.Fatal("Failed to approve token0 for minting:", err)
			}
			receipt, err := bind.WaitMined(context.Background(), client, tx)
			if err != nil {
				log.Fatal("Failed to wait for token0 approval:", err)
			}
			if receipt.Status == 0 {
				log.Fatal("Token0 approval failed")
			}
			fmt.Println("Token0 approved for minting")
		}

		if allowance1.Cmp(initialAmount1) < 0 {
			fmt.Println("Approving token1 for minting...")
			auth.GasLimit = uint64(300000)
			auth.GasPrice, err = client.SuggestGasPrice(context.Background())
			if err != nil {
				log.Fatal("Failed to get gas price:", err)
			}
			tx, err := token1.Approve(auth, common.HexToAddress(poolAddress), initialAmount1)
			if err != nil {
				log.Fatal("Failed to approve token1 for minting:", err)
			}
			receipt, err := bind.WaitMined(context.Background(), client, tx)
			if err != nil {
				log.Fatal("Failed to wait for token1 approval:", err)
			}
			if receipt.Status == 0 {
				log.Fatal("Token1 approval failed")
			}
			fmt.Println("Token1 approved for minting")
		}

		// Get current nonce
		nonce, err := client.PendingNonceAt(context.Background(), auth.From)
		if err != nil {
			log.Fatal("Failed to get nonce:", err)
		}

		auth.Nonce = big.NewInt(int64(nonce))
		auth.GasLimit = uint64(500000) // Higher gas limit for mint
		auth.GasPrice, err = client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatal("Failed to get gas price:", err)
		}

		// Mint initial liquidity
		fmt.Printf("Adding initial liquidity - Token0: %s, Token1: %s\n",
			initialAmount0.String(), initialAmount1.String())

		tx, err := pool.Mint(auth, initialAmount0, initialAmount1)
		if err != nil {
			log.Fatal("Failed to mint initial liquidity:", err)
		}

		fmt.Printf("Mint transaction sent, tx hash: %s. Waiting for confirmation...\n", tx.Hash().Hex())

		// Wait for the transaction to be mined
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		receipt, err := bind.WaitMined(ctx, client, tx)
		if err != nil {
			log.Fatal("Failed to wait for mint transaction:", err)
		}

		if receipt.Status == 1 {
			fmt.Println("Initial liquidity added successfully")

			// Get updated reserves
			reserves, err = pool.GetReserves(nil)
			if err != nil {
				log.Fatal("Failed to get updated reserves:", err)
			}
			fmt.Printf("Updated pool reserves - Token0: %s, Token1: %s\n",
				reserves.Reserve0.String(), reserves.Reserve1.String())
		} else {
			log.Fatal("Failed to add initial liquidity")
		}
	}

	iterations := 10 // Number of swaps
	for i := 0; i < iterations; i++ {
		// Get current reserves
		reserves, err = pool.GetReserves(nil)
		if err != nil {
			log.Fatal("Failed to get reserves:", err)
		}

		// Calculate swap amounts based on current reserves
		amount0In := big.NewInt(1000000000000000) // 0.001 token
		amount1In := big.NewInt(0)

		// Calculate expected output using constant product formula: x * y = k
		// amount0In * reserve1 / (reserve0 + amount0In)
		amount1Out := new(big.Int).Mul(amount0In, reserves.Reserve1)
		amount1Out = amount1Out.Div(amount1Out, new(big.Int).Add(reserves.Reserve0, amount0In))

		// Apply 0.3% fee
		fee := new(big.Int).Mul(amount1Out, big.NewInt(3))
		fee = fee.Div(fee, big.NewInt(1000))
		amount1Out = amount1Out.Sub(amount1Out, fee)

		amount0Out := big.NewInt(0)

		fmt.Printf("\nSwap %d parameters:\n", i)
		fmt.Printf("  Current reserves - Token0: %s, Token1: %s\n",
			reserves.Reserve0.String(), reserves.Reserve1.String())
		fmt.Printf("  Input: %s token0\n", amount0In.String())
		fmt.Printf("  Expected output: %s token1\n", amount1Out.String())
		fmt.Printf("  Fee: %s token1\n", fee.String())

		// Get the current nonce for the account
		nonce, err := client.PendingNonceAt(context.Background(), auth.From)
		if err != nil {
			log.Fatalf("Failed to get nonce: %v", err)
		}

		// Approve tokens for swapping if needed
		if allowance0.Cmp(amount0In) < 0 {
			fmt.Println("Approving token0 for swapping...")
			auth.Nonce = big.NewInt(int64(nonce))
			auth.GasLimit = uint64(300000)
			auth.GasPrice, err = client.SuggestGasPrice(context.Background())
			if err != nil {
				log.Fatal("Failed to get gas price:", err)
			}
			tx, err := token0.Approve(auth, common.HexToAddress(poolAddress), amount0In)
			if err != nil {
				log.Fatal("Failed to approve token0 for swapping:", err)
			}
			receipt, err := bind.WaitMined(context.Background(), client, tx)
			if err != nil {
				log.Fatal("Failed to wait for token0 approval:", err)
			}
			if receipt.Status == 0 {
				log.Fatal("Token0 approval failed")
			}
			fmt.Println("Token0 approved for swapping")

			// Update nonce after approval
			nonce++
		}

		// Create the swap transaction
		auth.Nonce = big.NewInt(int64(nonce))
		auth.GasLimit = uint64(300000)
		auth.GasPrice, err = client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatalf("Failed to get gas price: %v", err)
		}

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

			// Update reserves after successful swap
			newReserves, err := pool.GetReserves(nil)
			if err != nil {
				log.Printf("Failed to get updated reserves: %v\n", err)
			} else {
				fmt.Printf("Updated pool reserves - Token0: %s, Token1: %s\n",
					newReserves.Reserve0.String(), newReserves.Reserve1.String())
			}
		} else {
			log.Printf("Swap %d transaction failed, block number: %d\n", i, receipt.BlockNumber.Uint64())
		}

		// Add a small delay between transactions
		time.Sleep(1 * time.Second)
	}

	fmt.Println("All transactions completed")
}
