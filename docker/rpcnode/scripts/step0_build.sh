#!/usr/bin/env sh

# Input parameters
ARCH=$(uname -m)

# Build enid
echo "Building enid from local branch"
git config --global --add safe.directory /eni-chain/go-eni
LEDGER_ENABLED=false
make install
mkdir -p build/generated
echo "DONE" > build/generated/build.complete
