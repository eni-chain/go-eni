const { ethers } = require("ethers");
  const bip39 = require("bip39");
  const fs = require("fs");

  // Configuration parameters
  const CONFIG = {
    mnemonic: "party two quit over jaguar carry episode naive machine nothing borrow sell",
    senderPrivKey: "0x57acb95d82739866a5c29e40b0aa2590742ae50425b7dd5b5d279a986370189e",
    senderAddress: "0xF87A299e6bC7bEba58dbBe5a5Aa21d49bCD16D52",
    numAddresses: 100000,
    amountEth: 1,            // Initial ENI transfer amount
    transferAmountEth: 0.1,  // Secondary ENI transfer amount
    amountToken: 100,        // Initial ERC20 transfer amount
    transferAmountToken: 10, // Secondary ERC20 transfer amount
    chainId: 6912115,        // Chain ID parameterized
    gasPrice: 1000000000,    // Unified gasPrice
    outputFile: "init.txt",          // Initial ENI transfer
    transferFile: "transfer.txt"    // Secondary ENI transfer

  };

  async function main() {
    // ----------------------------------
    // Initialization (all offline)
    // ----------------------------------
    const { chainId, gasPrice } = CONFIG;
    const senderWallet = new ethers.Wallet(CONFIG.senderPrivKey);

    // Generate receiving addresses (HD wallet derivation)
    const seed = bip39.mnemonicToSeedSync(CONFIG.mnemonic);
    const rootNode = ethers.HDNodeWallet.fromSeed(seed);
    const receivers = Array.from({ length: CONFIG.numAddresses }, (_, i) =>
        rootNode.derivePath(`m/44'/60'/0'/0/${i}`).address
    );
    let mainNonce = 0;  // Main account nonce starts from 0
    // ----------------------------------
    // 1.1 Deploy ERC20 contract
    // ----------------------------------
    // console.log("\nGenerating ERC20 contract deployment...");
    // // a. Calculate contract address (using the next nonce of the main account)
    // const erc20Address = ethers.getCreateAddress({
    //   from: CONFIG.senderAddress,
    //   nonce: mainNonce
    // });
    // CONFIG.erc20Address = erc20Address;
    const initialEthTxs = [];
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
    // 1. Generate initial ENI transfer transactions (main account → receiving addresses)
    // ----------------------------------
    console.log("\nGenerating initial ENI transactions...");

    for (const receiver of receivers) {
      const txData = {
        to: receiver,
        value: ethers.parseEther(CONFIG.amountEth.toString()),
        nonce: mainNonce++,
        gasLimit: 21000,
        gasPrice,
        chainId
      };

      initialEthTxs.push(await senderWallet.signTransaction(txData));
    }
    fs.writeFileSync(CONFIG.outputFile, initialEthTxs.join("\n"));

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