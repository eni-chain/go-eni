#!/bin/bash

set -e

cd contracts
npm ci
#TODO:Devin uncomment the following line
#npx hardhat test --network seilocal test/EVMCompatabilityTest.js
npx hardhat test --network seilocal test/EVMPrecompileTest.js
npx hardhat test --network seilocal test/AssociateTest.js