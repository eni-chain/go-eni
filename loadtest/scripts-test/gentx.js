const { ethers } = require("ethers");
const bip39 = require("bip39");
const fs = require("fs");

// Configuration parameters
const CONFIG = {
    mnemonic: "party two quit over jaguar carry episode naive machine nothing borrow sell",
    senderPrivKey: "0x57acb95d82739866a5c29e40b0aa2590742ae50425b7dd5b5d279a986370189e",
    senderAddress: "0xF87A299e6bC7bEba58dbBe5a5Aa21d49bCD16D52",
    numFromAccounts: 1000,    // Number of "from" accounts
    numToAccounts: 1000,      // Number of "to" accounts
    txsPerAccount: 100,       // Transactions per "from" account
    amountEth: 0.1,           // ENI transfer amount
    chainId: 6912115,         // Chain ID
    gasPrice: 1000000000,     // Unified gasPrice
    outputFilePrefix: "tx_batch_", // Prefix for output files
};

async function main() {
    // ----------------------------------
    // Initialization (all offline)
    // ----------------------------------
    const { chainId, gasPrice } = CONFIG;
    const senderWallet = new ethers.Wallet(CONFIG.senderPrivKey);

    // Generate "from" accounts (assumed to have funds) using HD wallet
    const seed = bip39.mnemonicToSeedSync(CONFIG.mnemonic);
    const rootNode = ethers.HDNodeWallet.fromSeed(seed);
    const fromAccounts = Array.from({ length: CONFIG.numFromAccounts }, (_, i) => {
        const wallet = rootNode.derivePath(`m/44'/60'/0'/0/${i}`);
        return {
            address: wallet.address,
            privateKey: wallet.privateKey
        };
    });

    // Generate random "to" accounts
    const toAccounts = Array.from({ length: CONFIG.numToAccounts }, () => {
        const wallet = ethers.Wallet.createRandom();
        return wallet.address;
    });

    // ----------------------------------
    // Generate 100 transactions per "from" account
    // ----------------------------------
    console.log("\nGenerating 100,000 transactions...");
    const allTransactions = [];

    for (const fromAccount of fromAccounts) {
        const fromWallet = new ethers.Wallet(fromAccount.privateKey);

        for (let i = 0; i < CONFIG.txsPerAccount; i++) {
            // Pick a random "to" address
            const toAddress = toAccounts[Math.floor(Math.random() * toAccounts.length)];

            const txData = {
                to: toAddress,
                value: ethers.parseEther(CONFIG.amountEth.toString()),
                nonce: i, // Incremental nonce for each account
                gasLimit: 21000,
                gasPrice,
                chainId
            };

            const signedTx = await fromWallet.signTransaction(txData);
            allTransactions.push(signedTx);
        }
    }

    // ----------------------------------
    // Split transactions into 10 files (10,000 txs each)
    // ----------------------------------
    console.log("\nSplitting transactions into 10 files...");
    const txsPerFile = 10000;
    for (let i = 0; i < 10; i++) {
        const startIdx = i * txsPerFile;
        const endIdx = startIdx + txsPerFile;
        const fileTxs = allTransactions.slice(startIdx, endIdx);

        const outputFile = `${CONFIG.outputFilePrefix}${i + 1}.txt`;
        fs.writeFileSync(outputFile, fileTxs.join("\n"));
        console.log(`- Wrote ${fileTxs.length} transactions to ${outputFile}`);
    }

    // ----------------------------------
    // Output summary
    // ----------------------------------
    console.log("\nTransaction generation complete!");
    console.log(`- Total "from" accounts: ${fromAccounts.length}`);
    console.log(`- Total "to" accounts: ${toAccounts.length}`);
    console.log(`- Total transactions: ${allTransactions.length}`);
    console.log(`- Files created: 10 (${txsPerFile} transactions each)`);
}

main().catch(console.error);