#!/bin/bash

enidbin=$(which ~/go/bin/enid | tr -d '"')
keyname=$(printf "12345678\n" | $enidbin keys list --output json | jq ".[0].name" | tr -d '"')
keyaddress=$(printf "12345678\n" | $enidbin keys list --output json | jq ".[0].address" | tr -d '"')
chainid=$($enidbin status | jq ".NodeInfo.network" | tr -d '"')
enihome=$(git rev-parse --show-toplevel | tr -d '"')

cd $enihome || exit
echo "Deploying first set of contracts..."

beginning_block_height=$($enidbin status | jq -r '.SyncInfo.latest_block_height')
echo "$beginning_block_height" > $enihome/integration_test/contracts/wasm_beginning_block_height.txt
echo "$keyaddress"  > $enihome/integration_test/contracts/wasm_creator_id.txt

# store first set of contracts
for i in {1..10}
do
    echo "Storing first set contract #$i..."
    store_result=$(printf "12345678\n" | $enidbin tx wasm store integration_test/contracts/mars.wasm -y --from="$keyname" --chain-id="$chainid" --gas=5000000 --fees=1000000ueni --broadcast-mode=block --output=json)
    contract_id=$(echo "$store_result" | jq -r '.logs[].events[].attributes[] | select(.key == "code_id").value')
    printf "12345678\n" | $enidbin tx wasm instantiate "$contract_id" '{}' -y --no-admin --from="$keyname" --chain-id="$chainid" --gas=5000000 --fees=1000000ueni --broadcast-mode=block --label=dex --output=json
    echo "Got contract id $contract_id for iteration $i"
done

first_set_block_height=$($enidbin status | jq -r '.SyncInfo.latest_block_height')
echo "$first_set_block_height" > $enihome/integration_test/contracts/wasm_first_set_block_height.txt

sleep 5

# store second set of contracts
for i in {11..20}
do
    echo "Storing second set contract #$i..."
    store_result=$(printf "12345678\n" | $enidbin tx wasm store integration_test/contracts/saturn.wasm -y --from="$keyname" --chain-id="$chainid" --gas=5000000 --fees=1000000ueni --broadcast-mode=block --output=json)
    contract_id=$(echo "$store_result" | jq -r '.logs[].events[].attributes[] | select(.key == "code_id").value')
    printf "12345678\n" | $enidbin tx wasm instantiate "$contract_id" '{}' -y --no-admin --from="$keyname" --chain-id="$chainid" --gas=5000000 --fees=1000000ueni --broadcast-mode=block --label=dex --output=json
    echo "Got contract id $contract_id for iteration $i"
done

second_set_block_height=$($enidbin status | jq -r '.SyncInfo.latest_block_height')
echo "$second_set_block_height" > $enihome/integration_test/contracts/wasm_second_set_block_height.txt

sleep 5

# store third set of contracts
for i in {21..30}
do
    echo "Storing third set contract #$i..."
    store_result=$(printf "12345678\n" | $enidbin tx wasm store integration_test/contracts/venus.wasm -y --from="$keyname" --chain-id="$chainid" --gas=5000000 --fees=1000000ueni --broadcast-mode=block --output=json)
    contract_id=$(echo "$store_result" | jq -r '.logs[].events[].attributes[] | select(.key == "code_id").value')
    printf "12345678\n" | $enidbin tx wasm instantiate "$contract_id" '{}' -y --no-admin --from="$keyname" --chain-id="$chainid" --gas=5000000 --fees=1000000ueni --broadcast-mode=block --label=dex --output=json
    echo "Got contract id $contract_id for iteration $i"
done

third_set_block_height=$($enidbin status | jq -r '.SyncInfo.latest_block_height')
echo "$third_set_block_height" > $enihome/integration_test/contracts/wasm_third_set_block_height.txt

num_stored=$(enid q wasm list-code --count-total --limit 100 --output json | jq -r ".code_infos | length")
echo $num_stored

exit 0
