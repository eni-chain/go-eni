#!/bin/bash
enidbin=$(which ~/go/bin/enid | tr -d '"')
keyname=$(printf "12345678\n" | $enidbin keys list --output json | jq ".[0].name" | tr -d '"')
chainid=$($enidbin status | jq ".NodeInfo.network" | tr -d '"')
enihome=$(git rev-parse --show-toplevel | tr -d '"')

echo $keyname
echo $enidbin
echo $chainid
echo $enihome

# Deploy all contracts
echo "Deploying sei tester contract"

cd $enihome/loadtest/contracts
# store
echo "Storing..."

eni_tester_res=$(printf "12345678\n" | $enidbin tx wasm store eni_tester.wasm -y --from=$keyname --chain-id=$chainid --gas=5000000 --fees=1000000ueni --broadcast-mode=block --output=json)
eni_tester_id=$(python3 parser.py code_id $eni_tester_res)

# instantiate
echo "Instantiating..."
tester_in_res=$(printf "12345678\n" | $enidbin tx wasm instantiate $eni_tester_id '{}' -y --no-admin --from=$keyname --chain-id=$chainid --gas=5000000 --fees=1000000ueni --broadcast-mode=block  --label=dex --output=json)
tester_addr=$(python3 parser.py contract_address $tester_in_res)

# TODO fix once implemented in loadtest config
jq '.eni_tester_address = "'$tester_addr'"' $enihome/loadtest/config.json > $enihome/loadtest/config_temp.json && mv $enihome/loadtest/config_temp.json $enihome/loadtest/config.json


echo "Deployed contracts:"
echo $tester_addr

exit 0