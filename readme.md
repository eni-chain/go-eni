# Eni

![Banner!](assets/EniLogo.png)  
Website: 🌐 https://www.eniac.network  
Eni is the fastest general purpose L1 blockchain and the first parallelized EVM. This allows Eni to get the best of Solana and Ethereum - a hyper optimized execution layer that benefits from the tooling and mindshare around the EVM.

# Overview
**Eni** is a high-performance, low-fee, delegated proof-of-stake blockchain designed for developers. It supports optimistic parallel execution of both EVM and CosmWasm, opening up new design possibilities. With unique optimizations like twin turbo consensus and EniDB, Eni ensures consistent 400ms block times and a transaction throughput that’s orders of magnitude higher than Ethereum. This means faster, more cost-effective operations. Plus, Eni’s seamless interoperability between EVM and CosmWasm gives developers native access to the entire Cosmos ecosystem, including IBC tokens, multi-sig accounts, fee grants, and more.

# Documentation
For the most up-to-date documentation please visit:  
👉 http://docs.eniac.network/

# Eni Optimizations
Eni introduces four major innovations:

- Twin Turbo Consensus: This feature allows Eni to reach the fastest time to finality of any blockchain at 400ms, unlocking web2 like experiences for applications.
- Optimistic Parallelization: This feature allows developers to unlock parallel processing for their Ethereum applications, with no additional work.
- EniDB: This major upgrade allows Eni to handle the much higher rate of data storage, reads and writes which become extremely important for a high performance blockchain.
- Interoperable EVM: This allows existing developers in the Ethereum ecosystem to deploy their applications, tooling and infrastructure to Eni with no changes, while benefiting from the 100x performance improvements offered by Eni.

All these features combine to unlock a brand new, scalable design space for the Ethereum Ecosystem.

# Testnet
## Get started
**Validate on the Eni Testnet**
*Current Testnet: release/v0.2*

## Hardware Requirements
**Minimum**
* 64 GB RAM
* 1 TB NVME SSD
* 16 Cores (modern CPU's)

## Operating System

> Linux (x86_64) or Linux (amd64) Recommended Arch Linux

**Dependencies**
> Prerequisite:
> go1.24.2 + required.


## Deployment Guide

### 1. Clone git repository and build the node

```bash
git clone https://github.com/eni-chain/go-eni.git
cd go-eni
git checkout $VERSION
make build
```
### 2. Start a Single Node
```bash
 cd ${PROJECT_DIR}
 ./build/enid start --home=./eni-node

```
🔧 Config path for single node: ${PROJECT_DIR}/eni-node

### 3. Start a 4-Node Local Network
```bash
  make start4-node
```
🔧 Config path for multi-node setup: ${PROJECT_DIR}/eni-nodes

This will launch a 4-node validator testnet locally.


### 4. Stop the 4-Node Network
```bash
  make stop4-node
```
Gracefully stops all running processes from start4-node.

### 5. Clean Node Data
- Single Node Reset
```bash
  make reset-eni-node
```
- Multi-Node Reset
```bash
  make reset-multi-node
 ```
Removes all blockchain data and configuration from the respective node directories.

## Full node network access configuration
If you want to start a new full node, join the ENI mainnet, testnet, or other blockchain networks deployed by ENI code, you should follow these steps.

### 1. Clone git repository and build the node

```bash
git clone https://github.com/eni-chain/go-eni.git
cd go-eni
git checkout $VERSION
make build
```

### 2. Node configuration copy

The basic configuration of the new node must be consistent with the configuration of other nodes in the blockchain network to be joined. Therefore, a copy of the node configuration should be copied from any other node in the network first, and sensitive information such as the node's key should be deleted, and then sent to the new node user.

This step is applied by the new node user to the official blockchain network to be joined, and sent to the new node user after being operated by the official team.

#### 2.1 Configuration copy

Copy a configuration from any node to be added to the network.

```bash
cd go-eni/eni-nodes

# node5 should be replaced with the actual name of the new node
cp -r node4 node5
```
#### 2.2 Delete sensitive data

Delete sensitive and redundant data from the copy configuration.

```bash
# node5 should be replaced with the actual name of the new node
cd go-eni/eni-nodes/node5

# After deleting the sensitive data of the node and starting the new node, these two files will be automatically generated
rm -rf config/priv_validator_key.json
rm -rf config/node_key.json

# Delete redundant data
rm -rf data
rm -rf keyring-test
rm -rf config/gentx
```

After deleting sensitive data, the configuration copy is sent to the new node.

### 3 Configuration modification

#### 3.1 Configuration must be modified

The following configuration items are mandatory changes, otherwise network node name conflicts may occur.

```bash
# node5 should be replaced with the actual name of the new node
moniker = "node5"
```

#### 3.2 Optional configuration modification

The following content is optional to modify, but it is strongly recommended not to modify it. Generally, only when multiple nodes are started on one machine, in order to prevent port conflicts, the port of the configuration item will be modified.

**Config.toml configuration**

```bash
# TCP or UNIX socket address of the ABCI application,
# or the name of an ABCI application compiled in with the CometBFT binary
proxy_app = "tcp://127.0.0.1:26608"

[rpc]
# TCP or UNIX socket address for the RPC server to listen on
laddr = "tcp://127.0.0.1:26697"

[p2p]
# Address to listen for incoming connections
laddr = "tcp://0.0.0.0:26696"
```

**App.toml configuration**

```bash
[evm]
http_enabled = true
http_port = 8595
```

### 4 Node startup

After the configuration modification is completed, you can start the node to join the blockchain network.

```bash
#node5 should be replaced with the actual name of the new node
nohup ./build/enid start --home=./eni-nodes/node5 &> ./build/node.log &
```

# Build with Us!
If you are interested in building with Eni Network:  
👉 Visit our community hub: https://linktr.ee/ENI_OFFICIAL  
🐦 DM us on X: https://x.com/eni__official  
🌐 Website: https://www.eniac.network