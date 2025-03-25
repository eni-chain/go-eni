#!/bin/bash


# Default values
DEFAULT_NODE_COUNT=4
DEFAULT_ACCOUNT_COUNT=4

# Get node count and account count from command line arguments
NODE_COUNT=$1
ACCOUNT_COUNT=$2
SINGLE_MACHINE=${3:-true}

# Validate if node count is a positive integer
if [[ -z $NODE_COUNT || ! $NODE_COUNT =~ ^[0-9]+$ || $NODE_COUNT -le 0 ]]; then
    echo "No valid node count provided, using default: $DEFAULT_NODE_COUNT"
    NODE_COUNT=$DEFAULT_NODE_COUNT
fi

# Validate if account count is a positive integer
if [[ -z $ACCOUNT_COUNT || ! $ACCOUNT_COUNT =~ ^[0-9]+$ || $ACCOUNT_COUNT -le 0 ]]; then
    echo "No valid account count provided, using default: $DEFAULT_ACCOUNT_COUNT"
    ACCOUNT_COUNT=$DEFAULT_ACCOUNT_COUNT
fi

# Generate configuration files
for ((i=1; i<=NODE_COUNT; i++))
do
    # Initialize each node
    echo "Initializing node$i..."
    ./build/enid init "node$i" --chain-id mychain --home "./eni-nodes/node$i"

    # Generate validator account keys
    echo "Generating validator account keys for node$i..."
    ./build/enid keys add "validator$i"  --keyring-backend test --home "./eni-nodes/node$i"
done

for ((i=1; i<=NODE_COUNT; i++))
do
    #Initializing 100w ueni token to validator account
    echo "Initializing 100w ueni token to validator$i account..."
    ./build/enid genesis add-genesis-account $(./build/enid keys show validator$i -a --keyring-backend test --home "./eni-nodes/node$i") 1000000000000000000000000ueni --home "./eni-nodes/node1"
done


# Generate test account keys
for ((i=1; i<=ACCOUNT_COUNT; i++))
do
    # Generate test account keys to node1
    echo "Generating test$i account keys for node1..."
    ./build/enid keys add "test$i"  --keyring-backend test --home "./eni-nodes/node1"


done

for ((i=1; i<=ACCOUNT_COUNT; i++))
do
    #Initializing 100w ueni token to test account
    echo "Initializing 100w ueni token to test$i account..."
    ./build/enid genesis add-genesis-account $(./build/enid keys show test$i -a --keyring-backend test --home "./eni-nodes/node1") 1000000000000000000000000ueni --home "./eni-nodes/node1"
done

#Replace stake to ueni in genesis.json
perl -pi -e  's/stake/ueni/g' ./eni-nodes/node1/config/genesis.json
perl -pi -e  's|"minimum_fee_per_gas": "1000000000.000000000000000000"|"minimum_fee_per_gas": "100.000000000000000000"|' ./eni-nodes/node1/config/genesis.json

#Copy genesis.json to other nodes
for ((i=2; i<=NODE_COUNT; i++))
do
    #Copy genesis.json to other nodes
    echo "Copying genesis.json to node$i..."                                                                                    
    cp ./eni-nodes/node1/config/genesis.json ./eni-nodes/node$i/config/genesis.json
done

#Generate gentx for each node
for ((i=1; i<=NODE_COUNT; i++))
do
    #Generate gentx for each node
    echo "Generating gentx for node$i..."                                                                                    
    ./build/enid genesis gentx validator$i 100000000000000000000ueni --chain-id mychain --keyring-backend test --home ./eni-nodes/node$i
done


#merge gentx for each node
for ((i=1; i<=NODE_COUNT; i++))
do
    #merge gentx for each node
    echo "Merging gentx for node$i..."                                                                                    
    cp ./eni-nodes/node$i/config/gentx/* ./eni-nodes/node1/config/gentx/
done

#collect gentx
echo "Collecting gentx..."                                                                                    
./build/enid genesis collect-gentxs --home ./eni-nodes/node1


#Copy genesis.json to other nodes
for ((i=2; i<=NODE_COUNT; i++))
do
    #Copy genesis.json to other nodes
    echo "Copying genesis.json to node$i..."                                                                                    
    cp ./eni-nodes/node1/config/genesis.json ./eni-nodes/node$i/config/
done

if [[ $SINGLE_MACHINE == "true" ]]; then
    for ((i=2; i<=NODE_COUNT; i++))
    do
        # Calculate port offsets based on node ID
        PORT_OFFSET=$((i * 10))
        P2P_PORT=$((26656 + PORT_OFFSET))
        RPC_PORT=$((26657 + PORT_OFFSET))
        APP_PORT=$((26658 + PORT_OFFSET))

        # Replace ports in config.toml
        echo "Updating ports for node$i (P2P: $P2P_PORT, RPC: $RPC_PORT, APP: $APP_PORT)..."
        perl -pi -e  "s/:26656/:$P2P_PORT/g; s/:26657/:$RPC_PORT/g; s/:26658/:$APP_PORT/g" ./eni-nodes/node$i/config/config.toml
    done
fi

# Generate peers list
peers=""
for ((i=1; i<=NODE_COUNT; i++))
do
    # Get node ID for each node
    node_id=$(./build/enid comet show-node-id --home ./eni-nodes/node$i)
    echo "nodeId $node_id"
    # Calculate P2P port based on node ID
    P2P_PORT=$((26656 + i * 10 - 10))
    
    # Append peer to the list
    peers+="$node_id@localhost:$P2P_PORT,"
done

echo "peers $peers"
peers_escaped=${peers//@/\\@}

#update config param...
echo "update config param..."
for ((i=1; i<=NODE_COUNT; i++))
do
    #update config param...
    echo "update config param node$i..."
    perl -pi -e  's|minimum-gas-prices = ""|minimum-gas-prices = "0ueni"|' ./eni-nodes/node$i/config/app.toml
    perl -pi -e  's|allow_duplicate_ip = false|allow_duplicate_ip = true|' ./eni-nodes/node$i/config/config.toml
    perl -pi -e  's|laddr = "tcp://127.0.0.1:26657"|laddr = "tcp://0.0.0.0:26657"|' ./eni-nodes/node$i/config/config.toml
    perl -pi -e  "s|persistent_peers = \".*\"|persistent_peers = \"$peers_escaped\"|" ./eni-nodes/node$i/config/config.toml
    perl -pi -e  's|keyring-backend = "os"|keyring-backend = "test"|' ./eni-nodes/node$i/config/client.toml
    perl -pi -e  's|chain-id = ""|chain-id = "mychain"|' ./eni-nodes/node$i/config/client.toml
done

echo "Configuration files generated successfully!"

