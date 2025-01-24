# Eni

![Banner!](assets/EniLogo.png)

Eni is the fastest general purpose L1 blockchain and the first parallelized EVM. This allows Eni to get the best of Solana and Ethereum - a hyper optimized execution layer that benefits from the tooling and mindshare around the EVM.

# Overview
**Eni** is a high-performance, low-fee, delegated proof-of-stake blockchain designed for developers. It supports optimistic parallel execution of both EVM and CosmWasm, opening up new design possibilities. With unique optimizations like twin turbo consensus and EniDB, Eni ensures consistent 400ms block times and a transaction throughput that’s orders of magnitude higher than Ethereum. This means faster, more cost-effective operations. Plus, Eni’s seamless interoperability between EVM and CosmWasm gives developers native access to the entire Cosmos ecosystem, including IBC tokens, multi-sig accounts, fee grants, and more.

# Documentation
For the most up to date documentation please visit https://www.docs.eni.io/

# Eni Optimizations
Eni introduces four major innovations:

- Twin Turbo Consensus: This feature allows Eni to reach the fastest time to finality of any blockchain at 400ms, unlocking web2 like experiences for applications.
- Optimistic Parallelization: This feature allows developers to unlock parallel processing for their Ethereum applications, with no additional work.
- EniDB: This major upgrade allows Eni to handle the much higher rate of data storage, reads and writes which become extremely important for a high performance blockchain.
- Interoperable EVM: This allows existing developers in the Ethereum ecosystem to deploy their applications, tooling and infrastructure to Eni with no changes, while benefiting from the 100x performance improvements offered by Eni.

All these features combine to unlock a brand new, scalable design space for the Ethereum Ecosystem.

# Testnet
## Get started
**How to validate on the Eni Testnet**
*This is the Eni Atlantic-2 Testnet ()*

> Genesis [Published](https://github.com/eni-protocol/testnet/blob/main/atlantic-2/genesis.json)

## Hardware Requirements
**Minimum**
* 64 GB RAM
* 1 TB NVME SSD
* 16 Cores (modern CPU's)

## Operating System 

> Linux (x86_64) or Linux (amd64) Recommended Arch Linux

**Dependencies**
> Prerequisite: go1.18+ required.
* Arch Linux: `pacman -S go`
* Ubuntu: `sudo snap install go --classic`

> Prerequisite: git. 
* Arch Linux: `pacman -S git`
* Ubuntu: `sudo apt-get install git`

> Optional requirement: GNU make. 
* Arch Linux: `pacman -S make`
* Ubuntu: `sudo apt-get install make`

## Enid Installation Steps

**Clone git repository**

```bash
git clone https://github.com/eni-protocol/eni-chain
cd eni-chain
git checkout $VERSION
make install
```
**Generate keys**

* `enid keys add [key_name]`

* `enid keys add [key_name] --recover` to regenerate keys with your mnemonic

* `enid keys add [key_name] --ledger` to generate keys with ledger device

## Validator setup instructions

* Install enid binary

* Initialize node: `enid init <moniker> --chain-id eni-testnet-1`

* Download the Genesis file: `wget https://github.com/eni-protocol/testnet/raw/main/eni-testnet-1/genesis.json -P $HOME/.eni/config/`
 
* Edit the minimum-gas-prices in ${HOME}/.eni/config/app.toml: `sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0.01ueni"/g' $HOME/.eni/config/app.toml`

* Start enid by creating a systemd service to run the node in the background
`nano /etc/systemd/system/enid.service`
> Copy and paste the following text into your service file. Be sure to edit as you see fit.

```bash
[Unit]
Description=Eni-Network Node
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/root/
ExecStart=/root/go/bin/enid start
Restart=on-failure
StartLimitInterval=0
RestartSec=3
LimitNOFILE=65535
LimitMEMLOCK=209715200

[Install]
WantedBy=multi-user.target
```
## Start the node

**Start enid on Linux**

* Reload the service files: `sudo systemctl daemon-reload` 
* Create the symlinlk: `sudo systemctl enable enid.service` 
* Start the node sudo: `systemctl start enid && journalctl -u enid -f`

**Start a chain on 4 node docker cluster**

* Start local 4 node cluster: `make docker-cluster-start`
* SSH into a docker container: `docker exec -it [container_name] /bin/bash`
* Stop local 4 node cluster: `make docker-cluster-stop`

### Create Validator Transaction
```bash
enid tx staking create-validator \
--from {{KEY_NAME}} \
--chain-id  \
--moniker="<VALIDATOR_NAME>" \
--commission-max-change-rate=0.01 \
--commission-max-rate=1.0 \
--commission-rate=0.05 \
--details="<description>" \
--security-contact="<contact_information>" \
--website="<your_website>" \
--pubkey $(enid tendermint show-validator) \
--min-self-delegation="1" \
--amount <token delegation>ueni \
--node localhost:26657
```
# Build with Us!
If you are interested in building with Eni Network: 
Email us at team@eninetwork.io 
DM us on Twitter https://twitter.com/EniNetwork
