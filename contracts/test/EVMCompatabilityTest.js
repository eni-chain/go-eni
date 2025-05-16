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

let debugLog = false;
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

    describe("Deployment", function () {
      it("Should deploy successfully", async function () {
        expect(await evmTester.getAddress()).to.be.properAddress;
        expect(await testToken.getAddress()).to.be.properAddress;
        expect(await evmTester.getAddress()).to.not.equal(await testToken.getAddress());
      });

      it("Should have correct address", async function () {
        if(firstNonce > 0) {
          this.skip()
        } else {
          expect(await testToken.getAddress()).to.equal(firstContractAddress);
        }
      });

      it("Should estimate gas for a contract deployment", async function () {
        const callData = evmTester.interface.encodeFunctionData("createToken", ["TestToken", "TTK"]);
        const estimatedGas = await ethers.provider.estimateGas({
          to: await evmTester.getAddress(),
          data: callData
        });
        expect(estimatedGas).to.greaterThan(0);
      });
    });
  });
});
