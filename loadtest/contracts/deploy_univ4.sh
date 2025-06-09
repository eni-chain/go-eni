#!/bin/bash

set -e  # Exit on any error

echo "Starting deployment process..."
echo "Using account: 0xF87A299e6bC7bEba58dbBe5a5Aa21d49bCD16D52"

# Store the original directory
ORIGINAL_DIR=$(pwd)

# Change to the uniswap_v4 directory
cd "$(dirname "$0")/uniswap_v4" || {
    echo "Error: Could not change to uniswap_v4 directory, pwd is $ORIGINAL_DIR"
    exit 1
}

echo "Current directory: $(pwd)"

# Clean install dependencies
echo "Cleaning and reinstalling dependencies..."
rm -rf node_modules
rm -f package-lock.json
npm install
if [ $? -ne 0 ]; then
    echo "Failed to install dependencies"
    exit 1
fi

# Install OpenZeppelin contracts explicitly
echo "Installing OpenZeppelin contracts..."
npm install @openzeppelin/contracts@4.9.0
if [ $? -ne 0 ]; then
    echo "Failed to install OpenZeppelin contracts"
    exit 1
fi

# Check if the local node is running
check_node() {
    echo "Checking if local node is running..."
    if ! curl -s http://127.0.0.1:8545 > /dev/null; then
        echo "Error: Local node is not running at http://127.0.0.1:8545"
        echo "Please start a local node first"
        exit 1
    fi
    echo "Local node is running"
}

check_node

# Clean and compile contracts
echo -e "\nCleaning and compiling contracts..."
rm -rf cache artifacts
npx hardhat clean
npx hardhat compile --verbose
if [ $? -ne 0 ]; then
    echo "Failed to compile contracts"
    exit 1
fi

# Deploy ERC20 Token0
echo -e "\nDeploying Token0..."
TOKEN0_ADDRESS=$(npx hardhat run scripts/deploy_erc20.js --network localhost | grep -Eo "0x[a-fA-F0-9]{40}" | tail -n1)
if [ $? -ne 0 ]; then
    echo "Failed to deploy Token0"
    exit 1
fi
echo "Token0 Address: $TOKEN0_ADDRESS"

# Deploy ERC20 Token1
echo -e "\nDeploying Token1..."
TOKEN1_ADDRESS=$(npx hardhat run scripts/deploy_erc20.js --network localhost | grep -Eo "0x[a-fA-F0-9]{40}" | tail -n1)
if [ $? -ne 0 ]; then
    echo "Failed to deploy Token1"
    exit 1
fi
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
    if [ $? -eq 0 ] && [[ -n "$POOL_ADDRESS" ]]; then
        echo "Uniswap V4 Pool deployed at: $POOL_ADDRESS"
        break
    fi

    RETRY_COUNT=$((RETRY_COUNT + 1))
    echo "Deployment failed. Retrying in 5 seconds... (Attempt $RETRY_COUNT of $MAX_RETRIES)"
    sleep 5
done

if [[ -z "$POOL_ADDRESS" ]]; then
    echo "Pool deployment failed after $MAX_RETRIES attempts. Exiting."
    exit 1
fi

# Output summary
echo -e "\n Deployment Summary:"
echo "Token0: $TOKEN0_ADDRESS"
echo "Token1: $TOKEN1_ADDRESS"
echo "Pool  : $POOL_ADDRESS"

# Write to .env
echo -e "\nSaving addresses to .env file..."
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

echo -e "\nDeployment completed successfully!"
