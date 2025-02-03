#!/bin/bash

set -e

cd contracts
npm ci
npx hardhat test --network enilocal test/CW20toERC20PointerTest.js
npx hardhat test --network enilocal test/ERC20toCW20PointerTest.js
npx hardhat test --network enilocal test/ERC20toNativePointerTest.js
npx hardhat test --network enilocal test/CW721toERC721PointerTest.js
npx hardhat test --network enilocal test/ERC721toCW721PointerTest.js
npx hardhat test --network enilocal test/CW1155toERC1155PointerTest.js
npx hardhat test --network enilocal test/ERC1155toCW1155PointerTest.js
