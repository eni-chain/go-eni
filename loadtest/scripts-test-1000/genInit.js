const { ethers } = require("ethers");
const bip39 = require("bip39");
const fs = require("fs");

// Configuration parameters
const CONFIG = {
    mnemonic: "party two quit over jaguar carry episode naive machine nothing borrow sell",
    senderPrivKey: "0x57acb95d82739866a5c29e40b0aa2590742ae50425b7dd5b5d279a986370189e",
    senderAddress: "0xF87A299e6bC7bEba58dbBe5a5Aa21d49bCD16D52",
    numAddresses: 4400,
    fromAddress: 400,
    amountEth: 1,            // Initial ENI transfer amount
    transferAmountEth: 0.1,  // Secondary ENI transfer amount
    amountToken: 100,        // Initial ERC20 transfer amount
    transferAmountToken: 10, // Secondary ERC20 transfer amount
    chainId: 6912115,        // Chain ID parameterized
    gasPrice: 0,    // Unified gasPrice
    outputFile: "init4000.txt",          // Initial ENI transfer
    transferFile: "transfer1000.txt"    // Secondary ENI transfer
};

async function main() {
    // ----------------------------------
    // Initialization (all offline)
    // ----------------------------------

    const { chainId, gasPrice } = CONFIG;
    const senderWallet = new ethers.Wallet(CONFIG.senderPrivKey);

    // Generate "from" accounts (assumed funded) using HD wallet
    const senderSeed = bip39.mnemonicToSeedSync(CONFIG.mnemonic);
    const senderRootNode = ethers.HDNodeWallet.fromSeed(senderSeed);
    const Accounts = Array.from({ length: CONFIG.numAddresses }, (_, i) => {
        const wallet = senderRootNode.derivePath(`m/44'/60'/0'/0/${i}`);
        return {
            address: wallet.address,
            privateKey: wallet.privateKey
        };
    });


    fromAccounts = Accounts.splice(0,400)
    console.log("receiverAddress length==============",fromAccounts.length)
    console.log("receivers length==============",Accounts.length)


    let mainNonce = 0;  // Main account nonce starts from 0
    // ----------------------------------
    // 1. Generate initial ENI transfer transactions (main account → receiving addresses)
    // ----------------------------------
    console.log("\nGenerating initial ENI transactions...");
    const initialEthTxs = [];
    toIndex = 0
    for (let i = 0; i < 10; i++) {
        for (const from of fromAccounts) {
            const fromWallet = new ethers.Wallet(from.privateKey);
            const txData = {
                to: Accounts[toIndex].address,
                value: ethers.parseEther(CONFIG.amountEth.toString()),
                nonce: mainNonce,
                gasLimit: 21000,
                gasPrice,
                chainId
            };
            const signedTx = await fromWallet.signTransaction(txData);
            // console.log("toindex=====",toIndex)
            // console.log("signedTx=====",signedTx)
            initialEthTxs.push(signedTx);

            toIndex++
        }
        mainNonce++;
    }
    fs.writeFileSync(CONFIG.outputFile, initialEthTxs.join("\n"));
    return

    // // b. Deployment transaction (using the current nonce of the main account)
    // const deployTx = {
    //   data: CONFIG.erc20Bytecode,
    //   nonce: mainNonce++,
    //   gasLimit: 6000000,
    //   gasPrice,
    //   chainId
    // };
    // initialEthTxs.push(await senderWallet.signTransaction(deployTx));



    // ----------------------------------
    // 2. generate initial ERC20 transactions
    // ----------------------------------
    // console.log("\nGenerating ERC20 init transfers...");
    // const initErc20Txs = [];
    //
    // const erc20Interface = new ethers.Interface(ERC20_ABI);
    // for (const receiver of receivers) {
    //   const data = erc20Interface.encodeFunctionData("transfer", [
    //     receiver,
    //     CONFIG.amountToken
    //   ]);
    //
    //   const txData = {
    //     to: erc20Address,
    //     data,
    //     nonce: mainNonce++,
    //     gasLimit: 60000,
    //     gasPrice,
    //     chainId
    //   };
    //
    //   initErc20Txs.push(await senderWallet.signTransaction(txData));
    // }
    // fs.writeFileSync(CONFIG.outputErc20File, initErc20Txs.join("\n"));

    // ----------------------------------
    // 3. Generate secondary ENI transfers (receiving addresses → random addresses)
    // ----------------------------------
    console.log("\nGenerating secondary ENI transfers...");
    const secondaryEthTxs = [];

    for (const receiver of receivers) {
        const randomReceiver = receivers[Math.floor(Math.random() * receivers.length)];
        const derived = rootNode.derivePath(`m/44'/60'/0'/0/${receivers.indexOf(receiver)}`);
        const senderWallet = new ethers.Wallet(derived.privateKey);

        const txData = {
            to: randomReceiver,
            value: ethers.parseEther(CONFIG.transferAmountEth.toString()),
            nonce: 0, // Nonce of the receiving address starts from 0
            gasLimit: 21000,
            gasPrice,
            chainId
        };

        secondaryEthTxs.push(await senderWallet.signTransaction(txData));
    }
    fs.writeFileSync(CONFIG.transferFile, secondaryEthTxs.join("\n"));

    // ----------------------------------
    // 4. Generate secondary ERC20 transfers (receiving addresses → random addresses)
    // ----------------------------------
    // console.log("\nGenerating secondary ERC20 transfers...");
    // const secondaryErc20Txs = [];
    //
    // for (const receiver of receivers) {
    //   const randomReceiver = receivers[Math.floor(Math.random() * receivers.length)];
    //   const derived = rootNode.derivePath(`m/44'/60'/0'/0/${receivers.indexOf(receiver)}`);
    //   const senderWallet = new ethers.Wallet(derived.privateKey);
    //
    //   const data = erc20Interface.encodeFunctionData("transfer", [
    //     randomReceiver,
    //     CONFIG.transferAmountToken
    //   ]);
    //
    //   const txData = {
    //     to: erc20Address,
    //     data,
    //     nonce: 1, // Nonce of the receiving address starts from 1 (0 used for ENI transfer)
    //     gasLimit: 60000,
    //     gasPrice,
    //     chainId
    //   };
    //
    //   secondaryErc20Txs.push(await senderWallet.signTransaction(txData));
    // }
    // fs.writeFileSync(CONFIG.transferErc20File, secondaryErc20Txs.join("\n"));

    // ----------------------------------
    // Output summary
    // ----------------------------------
    console.log("\nAll transactions generated successfully!");
    console.log(`- Initial ENI transactions: ${initialEthTxs.length} (saved to ${CONFIG.outputFile})`);
    // console.log(`- ERC20 deployment + initial transfers: ${initErc20Txs.length} (saved to ${CONFIG.outputErc20File})`);
    console.log(`- Secondary ENI transfers: ${secondaryEthTxs.length} (saved to ${CONFIG.transferFile})`);
    // console.log(`- Secondary ERC20 transfers: ${secondaryErc20Txs.length} (saved to ${CONFIG.transferErc20File})`);
}

main().catch(console.error);