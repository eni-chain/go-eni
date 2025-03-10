# Eni

![Banner!](assets/EniLogo.png)

Eni is the fastest general purpose L1 blockchain and the first parallelized EVM. This allows Eni to get the best of Solana and Ethereum - a hyper optimized execution layer that benefits from the tooling and mindshare around the EVM.

# Overview
**Eni** is a high-performance, low-fee, delegated proof-of-stake blockchain designed for developers. It supports optimistic parallel execution of both EVM and CosmWasm, opening up new design possibilities. With unique optimizations like twin turbo consensus and EniDB, Eni ensures consistent 400ms block times and a transaction throughput that’s orders of magnitude higher than Ethereum. This means faster, more cost-effective operations. Plus, Eni’s seamless interoperability between EVM and CosmWasm gives developers native access to the entire Cosmos ecosystem, including IBC tokens, multi-sig accounts, fee grants, and more.

# Documentation
For the most up to date documentation please visit http://doc.eniac.network/

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
*This is the Eni newbase2 Testnet ()*

> Genesis [Published](https://github.com/eni-chain/go-eni/blob/newbase2/eni-node/config/genesis.json)

## Hardware Requirements
**Minimum**
* 64 GB RAM
* 1 TB NVME SSD
* 16 Cores (modern CPU's)

## Operating System

> Linux (x86_64) or Linux (amd64) Recommended Arch Linux

**Dependencies**
> Prerequisite: 
> go1.23+ required. 
> Ignite v28.7.0 required


## Enid Installation Steps

**Clone git repository**

```bash
git clone https://github.com/eni-chain/go-eni.git
cd go-eni
git checkout $VERSION
make install
```
**Generate keys**

* `enid keys add [key_name]`

* `enid keys add [key_name] --recover` to regenerate keys with your mnemonic

* `enid keys add [key_name] --ledger` to generate keys with ledger device

## Start the node

**Start enid on Linux/macOS**
```bash
enid start --home=${PROJECT_DIR}/eni-node
```


**Start by Ignite**

```bash 
ignite chain serve
```
# Build with Us!
If you are interested in building with Eni Network:
Email us at team@eninetwork.io
DM us on X https://x.com/eni__official/
