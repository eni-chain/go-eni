#!/bin/bash

# Deploy ERC20 tokens first
echo "Deploying ERC20 tokens..."
TOKEN0_ADDRESS=$(npx hardhat run scripts/deploy_erc20.js --network localhost)
TOKEN1_ADDRESS=$(npx hardhat run scripts/deploy_erc20.js --network localhost)

echo "Token0 deployed at: $TOKEN0_ADDRESS"
echo "Token1 deployed at: $TOKEN1_ADDRESS"

# Deploy Uniswap V4 Pool
echo "Deploying Uniswap V4 Pool..."
POOL_ADDRESS=$(npx hardhat run scripts/deploy_uniswap_v4_pool.js --network localhost --token0 $TOKEN0_ADDRESS --token1 $TOKEN1_ADDRESS)

echo "Uniswap V4 Pool deployed at: $POOL_ADDRESS"

# Save addresses to a file for later use
echo "TOKEN0_ADDRESS=$TOKEN0_ADDRESS" > .env
echo "TOKEN1_ADDRESS=$TOKEN1_ADDRESS" >> .env
echo "POOL_ADDRESS=$POOL_ADDRESS" >> .env

echo "Deployment completed!" 