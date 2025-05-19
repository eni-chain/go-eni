#!/bin/bash

# Create bindings directory if it doesn't exist
mkdir -p bindings

# Extract ABI from UniswapV4Pool.json
jq '.abi' artifacts/contracts/UniswapV4Pool.sol/UniswapV4Pool.json > artifacts/contracts/UniswapV4Pool.sol/UniswapV4Pool.abi

# Extract ABI from IERC20.json
jq '.abi' artifacts/contracts/IERC20.sol/IERC20.json > artifacts/contracts/IERC20.sol/IERC20.abi

# Generate bindings for UniswapV4Pool
abigen --abi=artifacts/contracts/UniswapV4Pool.sol/UniswapV4Pool.abi \
       --pkg=uniswap \
       --out=bindings/UniswapV4Pool.go \
       --type=UniswapV4Pool

# Generate bindings for IERC20
abigen --abi=artifacts/contracts/IERC20.sol/IERC20.abi \
       --pkg=uniswap \
       --out=bindings/IERC20.go \
       --type=IERC20 