#!/bin/bash

# Usage: ./cast_demo.sh <contract_address>
# Example: ./cast_demo.sh 0x1234567890abcdef1234567890abcdef12345678

# Check if contract address is provided
if [ -z "$1" ]; then
  echo "❌ Error: Please provide the contract address as an argument"
  echo "Usage: $0 <contract_address>"
  exit 1
fi

# Contract address from argument
CONTRACT=$1

# RPC and account settings
RPC_URL="http://localhost:8545"
PRIVATE_KEY="0x57acb95d82739866a5c29e40b0aa2590742ae50425b7dd5b5d279a986370189e"
FROM_ADDR="0xF87A299e6bC7bEba58dbBe5a5Aa21d49bCD16D52"
TO_ADDR="0xA0C9A775c1803Bfebc6049F2f7f257df00413BdD"

echo "✅ Contract address: $CONTRACT"
echo "✅ Sender address: $FROM_ADDR"
echo "✅ Receiver address: $TO_ADDR"
echo "✅ RPC URL: $RPC_URL"
echo "---------------------------------------"

# 1. Send 1 wei ETH to a random address
echo "➡️ Sending 1 wei ETH to 0xc1bbFB1358bA0E54B5Eb6cA4c0020F7DA669E6d1"
cast send --rpc-url "$RPC_URL" \
  0xc1bbFB1358bA0E54B5Eb6cA4c0020F7DA669E6d1 \
  --value 1wei \
  --from "$FROM_ADDR" \
  --private-key "$PRIVATE_KEY"

echo "---------------------------------------"

# 2. Check ETH balance of sender
echo "📊 Checking ETH balance of account: $FROM_ADDR"
cast balance "$FROM_ADDR" --rpc-url "$RPC_URL"

echo "---------------------------------------"

# 3. Check token balance (ERC20 balanceOf)
echo "📊 Checking token balance in contract $CONTRACT"
cast call "$CONTRACT" "balanceOf(address)" "$FROM_ADDR" --rpc-url "$RPC_URL"

echo "---------------------------------------"

# 4. Transfer 100 tokens (assuming 18 decimals)
echo "🔁 Transferring 100 tokens to $TO_ADDR"
cast send "$CONTRACT" "transfer(address,uint256)" \
  "$TO_ADDR" \
  100000000000000000000 \
  --rpc-url "$RPC_URL" \
  --private-key "$PRIVATE_KEY"

echo "---------------------------------------"

# 5. Transfer 1 token unit (for testing raw uint256 values)
echo "🔁 Transferring 1 token unit to $TO_ADDR"
cast send "$CONTRACT" "transfer(address,uint256)" \
  "$TO_ADDR" \
  1 \
  --rpc-url "$RPC_URL" \
  --private-key "$PRIVATE_KEY"

echo "✅ All operations completed."

