const { execSync } = require('child_process');
const { expect } = require("chai");
const fs = require('fs');
const path = require('path');

const { expectRevert } = require('@openzeppelin/test-helpers');
const { setupSigners, getAdmin, deployWasm, storeWasm, execute, isDocker, ABI, createTokenFactoryTokenAndMint, getEniBalance} = require("./lib");


describe("EVM Precompile Tester", function () {

    let accounts;
    let admin;

    before(async function () {
        accounts = await setupSigners(await hre.ethers.getSigners());
        admin = await getAdmin();
    })

    describe("EVM Addr Precompile Tester", function () {
        const AddrPrecompileContract = '0x0000000000000000000000000000000000001004';
        let addr;

        before(async function () {
            const signer = accounts[0].signer
            const contractABIPath = '../../precompiles/addr/abi.json';
            const contractABI = require(contractABIPath);
            // Get a contract instance
            addr = new ethers.Contract(AddrPrecompileContract, contractABI, signer);
        });

        it("Associates successfully", async function () {
            const unassociatedWallet = hre.ethers.Wallet.createRandom();
            try {
                await addr.getEniAddr(unassociatedWallet.address);
                expect.fail("Expected an error here since we look up an unassociated address");
            } catch (error) {
                expect(error).to.have.property('message').that.includes('execution reverted');
            }
            
            const message = `Please sign this message to link your EVM and Eni addresses. No ENI will be spent as a result of this signature.\n\n`;
            const messageLength = Buffer.from(message, 'utf8').length;
            const signatureHex = await unassociatedWallet.signMessage(message);

            const sig = hre.ethers.Signature.from(signatureHex);
            
            const appendedMessage = `\x19Ethereum Signed Message:\n${messageLength}${message}`;
            const associatedAddrs = await addr.associate(`0x${sig.v-27}`, sig.r, sig.s, appendedMessage)
            const addrs = await associatedAddrs.wait();
            expect(addrs).to.not.be.null;

            // Verify that addresses are now associated.
            const eniAddr = await addr.getEniAddr(unassociatedWallet.address);
            expect(eniAddr).to.not.be.null;
        });

        it("Associates with Public Key successfully", async function () {
            const unassociatedWallet = hre.ethers.Wallet.createRandom();
            try {
                await addr.getEniAddr(unassociatedWallet.address);
                expect.fail("Expected an error here since we look up an unassociated address");
            } catch (error) {
                expect(error).to.have.property('message').that.includes('execution reverted');
            }

            // Use the PublicKey without the '0x' prefix.
            const associatedAddrs = await addr.associatePubKey(unassociatedWallet.publicKey.slice(2))
            const addrs = await associatedAddrs.wait();
            expect(addrs).to.not.be.null;

            // Verify that addresses are now associated.
            const eniAddr = await addr.getEniAddr(unassociatedWallet.address);
            expect(eniAddr).to.not.be.null;
        });
    });

    describe("EVM Gov Precompile Tester", function () {
        const GovPrecompileContract = '0x0000000000000000000000000000000000001006';
        let gov;
        let govProposal;

        before(async function () {
            const govProposalResponse = JSON.parse(await execute(`enid tx gov submit-proposal param-change ../contracts/test/param_change_proposal.json --from admin --fees 20000ueni -b block -y -o json`))
            govProposal = govProposalResponse.logs[0].events[3].attributes[1].value;

            const signer = accounts[0].signer
            const contractABIPath = '../../precompiles/gov/abi.json';
            const contractABI = require(contractABIPath);
            // Get a contract instance
            gov = new ethers.Contract(GovPrecompileContract, contractABI, signer);
        });

        it("Gov deposit", async function () {
            const depositAmount = ethers.parseEther('0.01');
            const deposit = await gov.deposit(govProposal, {
                value: depositAmount,
            })
            const receipt = await deposit.wait();
            expect(receipt.status).to.equal(1);
        });
    });

    // TODO: Update when we add distribution query precompiles
    describe("EVM Distribution Precompile Tester", function () {
        const DistributionPrecompileContract = '0x0000000000000000000000000000000000001007';
        let distribution;
        before(async function () {
            const signer = accounts[0].signer;
            const contractABIPath = '../../precompiles/distribution/abi.json';
            const contractABI = require(contractABIPath);
            // Get a contract instance
            distribution = new ethers.Contract(DistributionPrecompileContract, contractABI, signer);
        });

        it("Distribution set withdraw address", async function () {
            const setWithdraw = await distribution.setWithdrawAddress(accounts[0].evmAddress)
            const receipt = await setWithdraw.wait();
            expect(receipt.status).to.equal(1);
        });
        it("Should query rewards and get non null response", async function () {
            const rewards = await distribution.rewards(accounts[0].evmAddress)
            expect(rewards).to.not.be.null;
        });
    });

    // TODO: Update when we add staking query precompiles
    describe("EVM Staking Precompile Tester", function () {
        const StakingPrecompileContract = '0x0000000000000000000000000000000000001005';
        let validatorAddr;
        let signer;
        let staking;

        before(async function () {
            validatorAddr = JSON.parse(await execute("enid q staking validators -o json")).validators[0].operator_address
            signer = accounts[0].signer;

            const contractABIPath = '../../precompiles/staking/abi.json';
            const contractABI = require(contractABIPath);

            staking = new ethers.Contract(StakingPrecompileContract, contractABI, signer);
        });

        it("Staking delegate", async function () {
            const delegateAmount = ethers.parseEther('0.01');
            const delegate = await staking.delegate(validatorAddr, {
                value: delegateAmount,
            });
            const receipt = await delegate.wait();
            expect(receipt.status).to.equal(1);

            const delegation = await staking.delegation(accounts[0].evmAddress, validatorAddr);
            expect(delegation).to.not.be.null;
            expect(delegation[0][0]).to.equal(10000n);

            const undelegate = await staking.undelegate(validatorAddr, delegation[0][0]);
            const undelegateReceipt = await undelegate.wait();
            expect(undelegateReceipt.status).to.equal(1);

            try {
                await staking.delegation(accounts[0].evmAddress, validatorAddr);
                expect.fail("Expected an error here since we undelegated the amount and delegation should not exist anymore.");
            } catch (error) {
                expect(error).to.have.property('message').that.includes('execution reverted');
            }
        });
    });

});

function parseHexToJSON(hexStr) {
    // Remove the 0x prefix
    hexStr = hexStr.slice(2);
    // Convert to bytes
    const bytes = Buffer.from(hexStr, 'hex');
    // Convert to JSON
    return JSON.parse(bytes.toString());
}