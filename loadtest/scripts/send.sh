#!/bin/bash

# Check if file parameter exists
if [ $# -lt 1 ]; then
    echo "Usage: $0 <transactions_file.txt> [rpc-url]"
    echo "Default rpc-url: http://localhost:8545"
    exit 1
fi

TX_FILE="$1"
RPC_URL="${2:-http://localhost:8545}"

# Verify if the file exists
if [ ! -f "$TX_FILE" ]; then
    echo "Error: File $TX_FILE not found"
    exit 1
fi

# Set failure counter
fail_count=0
success_count=0
total=$(grep -cve '^\s*$' "$TX_FILE")

# Process transactions line by line
while IFS= read -r raw_tx || [[ -n "$raw_tx" ]]; do
    # Clean transaction data
    tx=$(echo "$raw_tx" | tr -d '[:space:]')

    # Basic format validation
    if [[ ! $tx =~ ^0x[0-9a-fA-F]+$ ]]; then
        echo "Skipping invalid transaction: ${tx:0:20}..."
        ((fail_count++))
        continue
    fi

    # Send raw transaction
    echo "Sending transaction (${success_count}/${total}): ${tx:0:20}..."

    if output=$(cast rpc eth_sendRawTransaction "$tx" --rpc-url "$RPC_URL" 2>&1); then
        echo "Success | Hash: $output"
        ((success_count++))
    else
        echo "Failure | Error: $output"
        ((fail_count++))
    fi

done < "$TX_FILE"

# Output statistics
echo "Sending completed"
echo "Success: $success_count"
echo "Failure: $fail_count"
echo "Total: $total"