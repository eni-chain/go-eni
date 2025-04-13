const {ethers} = require("ethers");
const bip39 = require("bip39");
const fs = require("fs");

// Configuration parameters
const CONFIG = {
    senderMnemonic: "party two quit over jaguar carry episode naive machine nothing borrow sell",
    receiverMnemonic: "word sea trophy enhance rain glad skill drastic proof guitar lemon decline",
    senderPrivKey: "0x57acb95d82739866a5c29e40b0aa2590742ae50425b7dd5b5d279a986370189e",
    senderAddress: "0xF87A299e6bC7bEba58dbBe5a5Aa21d49bCD16D52",
    numFromAccounts: 1000,    // Number of "from" accounts
    numToAccounts: 1000,      // Number of "to" accounts
    txsPerAccount: 100,       // Transactions per "from" account
    amountEth: 0.1,           // ENI transfer amount
    chainId: 6912115,         // Chain ID
    gasPrice: 1000000000,     // Unified gas price
    outputFilePrefix: "tx_batch_", // Output file prefix
    txsPerFile: 10000,        // Transactions per file
};

async function main() {
    // Initialize (all offline)
    const {chainId, gasPrice} = CONFIG;
    const senderWallet = new ethers.Wallet(CONFIG.senderPrivKey);

    // Generate "from" accounts (assumed funded) using HD wallet
    const senderSeed = bip39.mnemonicToSeedSync(CONFIG.senderMnemonic);
    const senderRootNode = ethers.HDNodeWallet.fromSeed(senderSeed);
    const fromAccounts = Array.from({length: CONFIG.numFromAccounts}, (_, i) => {
        const wallet = senderRootNode.derivePath(`m/44'/60'/0'/0/${i}`);
        return {
            address: wallet.address,
            privateKey: wallet.privateKey
        };
    });

    // Generate fixed "to" accounts (using different mnemonic to avoid duplicates)
    const receiverSeed = bip39.mnemonicToSeedSync(CONFIG.receiverMnemonic);
    const receiverRootNode = ethers.HDNodeWallet.fromSeed(receiverSeed);
    const toAccounts = Array.from({length: CONFIG.numToAccounts}, (_, i) => {
        const wallet = receiverRootNode.derivePath(`m/44'/60'/0'/0/${i}`);
        return wallet.address;
    });

    // Verify no duplicates between "from" and "to" accounts
    const fromAddresses = new Set(fromAccounts.map(acc => acc.address));
    const toAddresses = new Set(toAccounts);
    const intersection = [...fromAddresses].filter(addr => toAddresses.has(addr));
    if (intersection.length > 0) {
        throw new Error(`Found ${intersection.length} duplicate accounts`);
    }
    console.log("\nVerification passed: No duplicates between from and to accounts");

    // Assign a fixed "to" address to each "from" account
    const fromToMapping = fromAccounts.map((fromAccount, index) => ({
        from: fromAccount,
        to: toAccounts[index % CONFIG.numToAccounts] // Cycle through "to" addresses
    }));

    // Generate 100,000 transactions
    console.log("\nGenerating 100,000 transactions...");
    const allTransactions = [];

    for (const {from, to} of fromToMapping) {
        const fromWallet = new ethers.Wallet(from.privateKey);

        for (let i = 0; i < CONFIG.txsPerAccount; i++) {
            const txData = {
                to: to, // Use fixed "to" address
                value: ethers.parseEther(CONFIG.amountEth.toString()),
                nonce: i, // Incremental nonce per account
                gasLimit: 21000,
                gasPrice,
                chainId
            };

            const signedTx = await fromWallet.signTransaction(txData);
            allTransactions.push(signedTx);
        }
    }

    // Split transactions into 10 files (10,000 transactions each)
    console.log("\nSplitting transactions into 10 files...");
    const numFiles = 10;
    for (let i = 0; i < numFiles; i++) {
        // Each file uses a distinct set of "from" accounts (100 accounts × 100 txs = 10,000 txs)
        const startAccountIdx = i * 100; // 100 accounts per file
        const fileTxs = [];

        for (let j = 0; j < 100; j++) { // 100 accounts
            const accountIdx = startAccountIdx + j;
            if (accountIdx >= CONFIG.numFromAccounts) break;

            const {from, to} = fromToMapping[accountIdx];
            const fromWallet = new ethers.Wallet(from.privateKey);

            // Generate 100 transactions for this account
            for (let k = 0; k < CONFIG.txsPerAccount; k++) {
                const txData = {
                    to: to,
                    value: ethers.parseEther(CONFIG.amountEth.toString()),
                    nonce: k,
                    gasLimit: 21000,
                    gasPrice,
                    chainId
                };

                const signedTx = await fromWallet.signTransaction(txData);
                fileTxs.push(signedTx);
            }
        }

        const outputFile = `${CONFIG.outputFilePrefix}${i + 1}.txt`;
        fs.writeFileSync(outputFile, fileTxs.join("\n"));
        console.log(`- Wrote ${fileTxs.length} transactions to ${outputFile}`);
    }

    // Output summary
    console.log("\nTransaction generation completed!");
    console.log(`- Total "from" accounts: ${fromAccounts.length}`);
    console.log(`- Total "to" accounts: ${toAccounts.length}`);
    console.log(`- Total transactions: ${allTransactions.length}`);
    console.log(`- Files created: ${numFiles} (each with ${CONFIG.txsPerFile} transactions)`);
}

main().catch(console.error);