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
PRIVATE_KEY="0x56087e01ed75db8413066cd3a92e6c01e7558005bd688f0719a97a50f213b702"
FROM_ADDR="0x3140aedbf686A3150060Cb946893b0598b266f5C"
TO_ADDR="0xA0C9A775c1803Bfebc6049F2f7f257df00413BdD"

echo "✅ Contract address: $CONTRACT"
echo "✅ Sender address: $FROM_ADDR"
echo "✅ Receiver address: $TO_ADDR"
echo "✅ RPC URL: $RPC_URL"
echo "---------------------------------------"

# 1. Send 1 wei ENI to a random address
echo "➡️ Sending 1 wei ENI to 0xc1bbFB1358bA0E54B5Eb6cA4c0020F7DA669E6d1"
/home/ubuntu/.foundry/bin/cast send --rpc-url "$RPC_URL" \
  0xc1bbFB1358bA0E54B5Eb6cA4c0020F7DA669E6d1 \
  --value 1wei \
  --from "$FROM_ADDR" \
  --private-key "$PRIVATE_KEY"

echo "---------------------------------------"

# 2. Check ENI balance of sender
echo "📊 Checking ENI balance of account: $FROM_ADDR"
/home/ubuntu/.foundry/bin/cast balance "$FROM_ADDR" --rpc-url "$RPC_URL"

echo "---------------------------------------"

# 3. Check token balance (ERC20 balanceOf)
echo "📊 Checking token balance in contract $CONTRACT"
/home/ubuntu/.foundry/bin/cast call "$CONTRACT" "balanceOf(address)" "$FROM_ADDR" --rpc-url "$RPC_URL"

echo "---------------------------------------"

# 4. Transfer 100 tokens (assuming 18 decimals)
echo "🔁 Transferring 100 tokens to $TO_ADDR"
/home/ubuntu/.foundry/bin/cast send "$CONTRACT" "transfer(address,uint256)" \
  "$TO_ADDR" \
  100000000000000000000 \
  --rpc-url "$RPC_URL" \
  --private-key "$PRIVATE_KEY"

echo "---------------------------------------"

# 5. Transfer 1 token unit (for testing raw uint256 values)
echo "🔁 Transferring 1 token unit to $TO_ADDR"
/home/ubuntu/.foundry/bin/cast send "$CONTRACT" "transfer(address,uint256)" \
  "$TO_ADDR" \
  1 \
  --rpc-url "$RPC_URL" \
  --private-key "$PRIVATE_KEY"

echo "✅ All operations completed."

