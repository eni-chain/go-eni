#!/bin/bash

set -e  # Exit on any error

echo "Starting deployment process..."
echo "Using account: 0xF87A299e6bC7bEba58dbBe5a5Aa21d49bCD16D52"

# Store the original directory
ORIGINAL_DIR=$(pwd)

# Change to the uniswap_v4 directory
cd uniswap_v4 || {
    echo "Error: Could not change to uniswap_v4 directory"
    exit 1
}

# Install dependencies if node_modules doesn't exist
if [ ! -d "node_modules" ]; then
    echo "Installing dependencies..."
    npm install
fi

# Check if the local node is running
check_node() {
    if ! curl -s http://127.0.0.1:8545 > /dev/null; then
        echo "Error: Local node is not running at http://127.0.0.1:8545"
        exit 1
    fi
}

check_node

# Compile contracts
echo -e "\n Compiling contracts..."
npx hardhat compile

# Deploy ERC20 Token0
echo -e "\n Deploying Token0..."
TOKEN0_ADDRESS=$(npx hardhat run scripts/deploy_erc20.js --network localhost | grep -Eo "0x[a-fA-F0-9]{40}" | tail -n1)
echo "Token0 Address: $TOKEN0_ADDRESS"

# Deploy ERC20 Token1
echo -e "\n Deploying Token1..."
TOKEN1_ADDRESS=$(npx hardhat run scripts/deploy_erc20.js --network localhost | grep -Eo "0x[a-fA-F0-9]{40}" | tail -n1)
echo "Token1 Address: $TOKEN1_ADDRESS"

# Validate addresses
if [[ -z "$TOKEN0_ADDRESS" || -z "$TOKEN1_ADDRESS" ]]; then
    echo "Failed to deploy ERC20 tokens. Exiting."
    exit 1
fi

# Export token addresses for the JS deployment script
export TOKEN0_ADDRESS
export TOKEN1_ADDRESS

# Deploy Uniswap V4 Pool with retry
echo -e "\nDeploying Uniswap V4 Pool..."
MAX_RETRIES=3
RETRY_COUNT=0
POOL_ADDRESS=""

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    POOL_ADDRESS=$(npx hardhat run scripts/deploy_uniswap_v4_pool.js --network localhost | grep -Eo "0x[a-fA-F0-9]{40}" | tail -n1)
    if [[ -n "$POOL_ADDRESS" ]]; then
        echo "Uniswap V4 Pool deployed at: $POOL_ADDRESS"
        break
    fi

    RETRY_COUNT=$((RETRY_COUNT + 1))
    echo "Deployment failed. Retrying in 5 seconds... (Attempt $RETRY_COUNT of $MAX_RETRIES)"
    sleep 5
done

if [[ -z "$POOL_ADDRESS" ]]; then
    echo " Pool deployment failed after $MAX_RETRIES attempts. Exiting."
    exit 1
fi

# Output summary
echo -e "\n Deployment Summary:"
echo "Token0: $TOKEN0_ADDRESS"
echo "Token1: $TOKEN1_ADDRESS"
echo "Pool  : $POOL_ADDRESS"

# Write to .env
echo -e "\n Saving addresses to .env file..."
cat > .env <<EOL
# Auto-generated on $(date)
TOKEN0_ADDRESS=$TOKEN0_ADDRESS
TOKEN1_ADDRESS=$TOKEN1_ADDRESS
POOL_ADDRESS=$POOL_ADDRESS
PRIVATE_KEY=57acb95d82739866a5c29e40b0aa2590742ae50425b7dd5b5d279a986370189e
RPC_URL=http://localhost:8545
EOL

# Return to the original directory
cd "$ORIGINAL_DIR"

echo -e "\n Deployment completed successfully!"
