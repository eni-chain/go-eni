#!/bin/bash

set -e

cd contracts
npm ci
npx hardhat test --network enilocal test/EVMCompatabilityTest.js
npx hardhat test --network enilocal test/EVMPrecompileTest.js
npx hardhat test --network enilocal test/AssociateTest.js