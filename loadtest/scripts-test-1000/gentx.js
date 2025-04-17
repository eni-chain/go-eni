const { ethers } = require("ethers");
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
    txsPerAccountPerBatch: 10, // Transactions per account per batch
    numBatches: 10,           // Number of batches
    amountEth: 0.0001,           // ENI transfer amount
    chainId: 6912115,         // Chain ID
    gasPrice: 0,     // Unified gas price
    outputFilePrefix: "newtx_batch1000_", // Output file prefix
    txsPerFile: 10000,        // Transactions per file (1,000 accounts × 10 txs)
};

async function main() {
    // Initialize (all offline)
    const { chainId, gasPrice } = CONFIG;
    const senderWallet = new ethers.Wallet(CONFIG.senderPrivKey);

    // Generate "from" accounts (assumed funded) using HD wallet
    const senderSeed = bip39.mnemonicToSeedSync(CONFIG.senderMnemonic);
    const senderRootNode = ethers.HDNodeWallet.fromSeed(senderSeed);
    const fromAccounts = Array.from({ length: CONFIG.numFromAccounts }, (_, i) => {
        const wallet = senderRootNode.derivePath(`m/44'/60'/0'/0/${i}`);
        return {
            address: wallet.address,
            privateKey: wallet.privateKey
        };
    });

    // Generate fixed "to" accounts (using different mnemonic to avoid duplicates)
    const receiverSeed = bip39.mnemonicToSeedSync(CONFIG.receiverMnemonic);
    const receiverRootNode = ethers.HDNodeWallet.fromSeed(receiverSeed);
    const toAccounts = Array.from({ length: CONFIG.numToAccounts }, (_, i) => {
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

    // Generate 10 batches, each with 10,000 transactions (1,000 accounts × 10 txs)
    console.log("\nGenerating 10 batches of 10,000 transactions each...");
    for (let batch = 0; batch < CONFIG.numBatches; batch++) {
        const fileTxs = [];
        const nonceOffset = batch * CONFIG.txsPerAccountPerBatch; // Nonce start for this batch

        // Generate transactions for all 1,000 accounts in this batch
        for (const { from, to } of fromToMapping) {
            const fromWallet = new ethers.Wallet(from.privateKey);

            // Generate 10 transactions per account for this batch
            for (let i = 0; i < CONFIG.txsPerAccountPerBatch; i++) {
                const nonce = nonceOffset + i; // Nonce for this transaction
                const txData = {
                    to: to, // Use fixed "to" address
                    value: ethers.parseEther(CONFIG.amountEth.toString()),
                    nonce: nonce, // Incremental nonce across batches
                    gasLimit: 21000,
                    gasPrice,
                    chainId
                };

                const signedTx = await fromWallet.signTransaction(txData);
                fileTxs.push(signedTx);
            }
        }

        // Write batch to file
        const outputFile = `${CONFIG.outputFilePrefix}${batch + 1}.txt`;
        fs.writeFileSync(outputFile, fileTxs.join("\n"));
        console.log(`- Wrote ${fileTxs.length} transactions to ${outputFile} (batch ${batch + 1})`);
    }

    // Output summary
    console.log("\nTransaction generation completed!");
    console.log(`- Total "from" accounts: ${fromAccounts.length}`);
    console.log(`- Total "to" accounts: ${toAccounts.length}`);
    console.log(`- Total transactions: ${CONFIG.numFromAccounts * CONFIG.txsPerAccountPerBatch * CONFIG.numBatches}`);
    console.log(`- Files created: ${CONFIG.numBatches} (each with ${CONFIG.txsPerFile} transactions)`);
}

main().catch(console.error);