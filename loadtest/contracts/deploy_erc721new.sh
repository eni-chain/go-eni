#!/bin/bash
#set -x
# This script is used to deploy the NoopToken contract to the target network
# This avoids trying to predict what address it might be deployed to

evm_endpoint=$1

# first fund account if necessary
THRESHOLD=100000000 # 100 Eth
ACCOUNT="0x251604eBfD1ddeef1F4f40b8F9Fc425538BE1339"
BALANCE=$(cast balance $ACCOUNT --rpc-url "$evm_endpoint")
echo "balance $BALANCE"
if (( $(echo "$BALANCE < $THRESHOLD" | bc -l) )); then
  printf "12345678\n" | ~/go/bin/enid tx evm send $ACCOUNT 100000000000000000000 --from admin --evm-rpc "$evm_endpoint"
  sleep 3
fi

cd loadtest/contracts/evm || exit 1

./setup.sh > /dev/null
ownerOf
git submodule update --init --recursive > /dev/null

INPUT=$(forge create -r "$evm_endpoint" --private-key dc9bb398d00f7778a61dcbb7e90cfe527b7e7b69ce9d557a08d5e32ea8d3eac0 src/ERC721.sol:MyNFT --json | jq -r '.transaction.input')

deploy_output=$(cast send --rpc-url http://localhost:8545 --private-key  dc9bb398d00f7778a61dcbb7e90cfe527b7e7b69ce9d557a08d5e32ea8d3eac0 --gas-limit 700000000 --gas-price 2000000 --create $INPUT --json)

echo "deploy output: $deploy_output"

contract_address=$(echo "$deploy_output" | jq -r '.contractAddress')
echo "deploy contract address: $contract_address"


cast call $contract_address "getFixedString()" --rpc-url "$evm_endpoint"