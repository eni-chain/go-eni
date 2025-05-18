# Uniswap V4 Load Testing

This directory contains scripts and contracts for load testing Uniswap V4 pools.

## Prerequisites

- Node.js and npm
- Go 1.16 or later
- Hardhat
- A local Ethereum node (e.g., Ganache)

## Setup

1. Install dependencies:
```bash
npm install
```

2. Create a `.env` file with the following variables:
```
PRIVATE_KEY=your_private_key_here
```

3. Start your local Ethereum node (e.g., Ganache)

## Deployment

1. Deploy the contracts:
```bash
chmod +x deploy.sh
./deploy.sh
```

This will:
- Deploy two ERC20 tokens
- Deploy a Uniswap V4 pool with the two tokens
- Save the contract addresses to `.env`

## Load Testing

1. Run the basic load test:
```bash
go run scripts/uniswap_v4_load_test.go
```

This will:
- Connect to the local Ethereum node
- Approve tokens for the pool
- Add initial liquidity
- Perform 100 swaps
- Measure and report performance metrics

2. Run the concurrent load test:
```bash
go run scripts/uniswap_v4_load_test.go -concurrent -workers 10 -iterations 1000
```

This will:
- Run the same test as above but with multiple concurrent workers
- Each worker will perform swaps in parallel
- Useful for testing pool performance under high load

## Contract Details

### UniswapV4Pool.sol

A simplified version of the Uniswap V4 pool contract that implements:
- Token swaps
- Liquidity provision
- Liquidity removal
- Reserve tracking

### Key Functions

- `mint(address to)`: Add liquidity to the pool
- `burn(address to)`: Remove liquidity from the pool
- `swap(uint256 amount0Out, uint256 amount1Out, address to)`: Perform a token swap
- `getReserves()`: Get current pool reserves

## Performance Metrics

The load test scripts measure:
- Total test duration
- Average transaction time
- Number of failed transactions
- Gas usage per transaction

## Notes

- The test uses a local Ethereum node for testing
- Make sure to have enough ETH in the test account for gas fees
- The test tokens are minted with a large supply for testing
- Adjust the number of iterations and workers based on your needs 