const { expect } = require("chai");
const {isBigNumber} = require("hardhat/common");
const {uniq, shuffle} = require("lodash");
const { ethers, upgrades } = require('hardhat');
const { getImplementationAddress } = require('@openzeppelin/upgrades-core');
const { deployEvmContract, setupSigners, fundAddress, getCosmosTx, getEvmTx} = require("./lib")
const axios = require("axios");
const { default: BigNumber } = require("bignumber.js");

function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

async function delay() {
  // await sleep(3000)
}

function debug(msg) {
  // leaving commented out to make output readable (unless debugging)
  console.log(msg)
}

async function sendTransactionAndCheckGas(sender, recipient, amount) {
  // Get the balance of the sender before the transaction
  const balanceBefore = await ethers.provider.getBalance(sender.address);

  // Send the transaction
  const tx = await sender.sendTransaction({
    to: recipient.address,
    value: amount
  });

  // Wait for the transaction to be mined and get the receipt
  const receipt = await tx.wait();

  // Get the balance of the sender after the transaction
  const balanceAfter = await ethers.provider.getBalance(sender.address);

  // Calculate the total cost of the transaction (amount + gas fees)
  const gasPrice = receipt.gasPrice;
  const gasUsed = receipt.gasUsed;
  const totalCost = gasPrice * gasUsed + BigInt(amount);

  // Check that the sender's balance decreased by the total cost
  return balanceBefore - balanceAfter === totalCost
}

function generateWallet() {
  const wallet = ethers.Wallet.createRandom();
  return wallet.connect(ethers.provider);
}

function generateWallets(num) {
  const arr = []
  for(let i=0; i<num; i++) {
      const wallet = ethers.Wallet.createRandom();
      arr.push(wallet);
  }
  return arr;
}

async function sendTx(sender, txn, responses) {
  const txResponse = await sender.sendTransaction(txn);
  responses.push({nonce: txn.nonce, response: txResponse})
}

let debugLog = true;
if(debugLog) {
  // Intercept and log the JSON-RPC request and response
  const originalSend = ethers.provider.send;
  ethers.provider.send = async (method, params) => {
    console.log("JSON-RPC Request:", JSON.stringify({method, params}, null, 2));
    const result = await originalSend.call(ethers.provider, method, params);
    console.log("JSON-RPC Response:", JSON.stringify(result, null, 2));
    return result;
  };
}

describe("EVM Test", function () {

  describe("EVMCompatibilityTester", function () {
    let evmTester;
    let testToken;
    let owner;
    let evmAddr;
    let firstNonce;

    // The first contract address deployed from 0xF87A299e6bC7bEba58dbBe5a5Aa21d49bCD16D52
    // should always be 0xbD5d765B226CaEA8507EE030565618dAFFD806e2 when sent with nonce=0
    const firstContractAddress = "0xbD5d765B226CaEA8507EE030565618dAFFD806e2";
    // This function deploys a new instance of the contract before each test
    beforeEach(async function () {
      if(evmTester && testToken) {
        return
      }
      const accounts = await setupSigners(await ethers.getSigners())
      owner = accounts[0].signer;
      debug(`OWNER = ${owner.address}`)

      firstNonce = await ethers.provider.getTransactionCount(owner.address)

      const TestToken = await ethers.getContractFactory("TestToken")
      testToken = await TestToken.deploy("TestToken", "TTK");

      const EVMCompatibilityTester = await ethers.getContractFactory("EVMCompatibilityTester");
      evmTester = await EVMCompatibilityTester.deploy({ gasPrice: ethers.parseUnits('100', 'gwei') });

      await Promise.all([evmTester.waitForDeployment(), testToken.waitForDeployment()])

      let tokenAddr = await testToken.getAddress()
      evmAddr = await evmTester.getAddress()

      debug(`Token: ${tokenAddr}, EvmAddr: ${evmAddr}`);
    });


    describe("Contract Upgradeability", function() {
      it("Should allow for contract upgrades", async function() {
        // deploy BoxV1
        console.log("starting test")
        const Box = await ethers.getContractFactory("Box");
        console.log("Box Factory created")
        sleep(100)
        const val = 42;
        const box = await upgrades.deployProxy(Box, [val], { initializer: 'store' });
        console.log("Box deployed")
        sleep(100)
        const boxReceipt = await box.waitForDeployment()
        console.log("Box deployed to:", boxReceipt.contractAddress);
        sleep(100)
        const boxAddr = await box.getAddress();
        console.log("Proxy Contract Address:", boxAddr);
        const implementationAddress = await upgrades.erc1967.getImplementationAddress(boxAddr);
        console.log("Implementation Contract Address:", implementationAddress);

        // make sure you can retrieve the value
        const retrievedValue = await box.retrieve();
        expect(retrievedValue).to.equal(val);

        // increment value
        debug("Incrementing value...")
        const resp = await box.boxIncr({ gasPrice: ethers.parseUnits('100', 'gwei') });
        await resp.wait();

        // make sure value is incremented
        const retrievedValue1 = await box.retrieve();
        expect(retrievedValue1).to.equal(val+1);

        // upgrade to BoxV2
        const BoxV2 = await ethers.getContractFactory('BoxV2');
        console.log("BoxV2 Factory created");
        sleep(100)
        debug('Upgrading Box...');
        const box2 = await upgrades.upgradeProxy(boxAddr, BoxV2, [val+1], { initializer: 'store' });
        console.log("BoxV2 deployed");
        sleep(100)
        await box2.deployTransaction.wait();
        debug('Box upgraded');
        sleep(100)
        const boxV2Addr = await box2.getAddress();
        expect(boxV2Addr).to.equal(boxAddr); // should be same address as it should be the proxy
        console.log('BoxV2 deployed to:', boxV2Addr);
        const boxV2 = await BoxV2.attach(boxV2Addr);
        sleep(100)
        // check that value is still the same
        debug("Calling boxV2 retrieve()...")
        const retrievedValue2 = await boxV2.retrieve();
        console.log("retrievedValue2 = ", retrievedValue2)
        expect(retrievedValue2).to.equal(val+1);

        // use new function in boxV2 and increment value
        debug("Calling boxV2 boxV2Incr()...")
        const txResponse = await boxV2.boxV2Incr();
        await txResponse.wait();

        // make sure value is incremented
        expect(await boxV2.retrieve()).to.equal(val+2);

        // store something in value2 and check it(check value2)
        const store2Resp = await boxV2.store2(10);
        await store2Resp.wait();
        expect(await boxV2.retrieve2()).to.equal(10);

        // ensure value is still the same in boxV2 (checking for any storage corruption)
        expect(await boxV2.retrieve()).to.equal(val+2);
      });
    });

  });
});
