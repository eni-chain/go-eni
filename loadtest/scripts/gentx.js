const { ethers } = require("ethers");
  const bip39 = require("bip39");
  const fs = require("fs");

  // Configuration parameters
  const CONFIG = {
    mnemonic: "party two quit over jaguar carry episode naive machine nothing borrow sell",
    senderPrivKey: "0x57acb95d82739866a5c29e40b0aa2590742ae50425b7dd5b5d279a986370189e",
    senderAddress: "0xF87A299e6bC7bEba58dbBe5a5Aa21d49bCD16D52",
    numAddresses: 10000,
    amountEth: 1,            // Initial ENI transfer amount
    transferAmountEth: 0.1,  // Secondary ENI transfer amount
    amountToken: 100,        // Initial ERC20 transfer amount
    transferAmountToken: 10, // Secondary ERC20 transfer amount
    chainId: 6912115,        // Chain ID parameterized
    gasPrice: 1000000000,    // Unified gasPrice
    outputFile: "init.txt",          // Initial ENI transfer
    transferFile: "transfer.txt",    // Secondary ENI transfer
    outputErc20File: "init_erc20.txt",    // Initial ERC20 deployment + transfer
    transferErc20File: "transfer_erc20.txt",  // Secondary ERC20 transfer
    erc20Address: "",
    erc20Bytecode: "0x60806040523480156200001157600080fd5b506040518060400160405280600d81526020016c2637b0b22a32b9ba2a37b5b2b760991b8152506040518060400160405280600381526020016213151560ea1b8152508160039081620000659190620002b4565b506004620000748282620002b4565b50505062000094336a52b7d2dcc80cd2e40000006200009a60201b60201c565b620003a8565b6001600160a01b038216620000ca5760405163ec442f0560e01b8152600060048201526024015b60405180910390fd5b620000d860008383620000dc565b5050565b6001600160a01b0383166200010b578060026000828254620000ff919062000380565b909155506200017f9050565b6001600160a01b03831660009081526020819052604090205481811015620001605760405163391434e360e21b81526001600160a01b03851660048201526024810182905260448101839052606401620000c1565b6001600160a01b03841660009081526020819052604090209082900390555b6001600160a01b0382166200019d57600280548290039055620001bc565b6001600160a01b03821660009081526020819052604090208054820190555b816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516200020291815260200190565b60405180910390a3505050565b634e487b7160e01b600052604160045260246000fd5b600181811c908216806200023a57607f821691505b6020821081036200025b57634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620002af57600081815260208120601f850160051c810160208610156200028a5750805b601f850160051c820191505b81811015620002ab5782815560010162000296565b5050505b505050565b81516001600160401b03811115620002d057620002d06200020f565b620002e881620002e1845462000225565b8462000261565b602080601f831160018114620003205760008415620003075750858301515b600019600386901b1c1916600185901b178555620002ab565b600085815260208120601f198616915b82811015620003515788860151825594840194600190910190840262000330565b5085821015620003705788850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b80820180821115620003a257634e487b7160e01b600052601160045260246000fd5b92915050565b61072180620003b86000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c8063313ce56711610066578063313ce567146100fe57806370a082311461010d57806395d89b4114610136578063a9059cbb1461013e578063dd62ed3e1461015157600080fd5b806306fdde0314610098578063095ea7b3146100b657806318160ddd146100d957806323b872dd146100eb575b600080fd5b6100a061018a565b6040516100ad919061056b565b60405180910390f35b6100c96100c43660046105d5565b61021c565b60405190151581526020016100ad565b6002545b6040519081526020016100ad565b6100c96100f93660046105ff565b610236565b604051601281526020016100ad565b6100dd61011b36600461063b565b6001600160a01b031660009081526020819052604090205490565b6100a061025a565b6100c961014c3660046105d5565b610269565b6100dd61015f36600461065d565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b60606003805461019990610690565b80601f01602080910402602001604051908101604052809291908181526020018280546101c590610690565b80156102125780601f106101e757610100808354040283529160200191610212565b820191906000526020600020905b8154815290600101906020018083116101f557829003601f168201915b5050505050905090565b60003361022a818585610277565b60019150505b92915050565b600033610244858285610289565b61024f85858561030d565b506001949350505050565b60606004805461019990610690565b60003361022a81858561030d565b610284838383600161036c565b505050565b6001600160a01b0383811660009081526001602090815260408083209386168352929052205460001981101561030757818110156102f857604051637dc7a0d960e11b81526001600160a01b038416600482015260248101829052604481018390526064015b60405180910390fd5b6103078484848403600061036c565b50505050565b6001600160a01b03831661033757604051634b637e8f60e11b8152600060048201526024016102ef565b6001600160a01b0382166103615760405163ec442f0560e01b8152600060048201526024016102ef565b610284838383610441565b6001600160a01b0384166103965760405163e602df0560e01b8152600060048201526024016102ef565b6001600160a01b0383166103c057604051634a1406b160e11b8152600060048201526024016102ef565b6001600160a01b038085166000908152600160209081526040808320938716835292905220829055801561030757826001600160a01b0316846001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9258460405161043391815260200190565b60405180910390a350505050565b6001600160a01b03831661046c5780600260008282546104619190620006ca565b909155506104de9050565b6001600160a01b038316600090815260208190526040902054818110156104bf5760405163391434e360e21b81526001600160a01b038516600482015260248101829052604481018390526064016102ef565b6001600160a01b03841660009081526020819052604090209082900390555b6001600160a01b0382166104fa57600280548290039055610519565b6001600160a01b03821660009081526020819052604090208054820190555b816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8360405161055e91815260200190565b60405180910390a3505050565b600060208083528351808285015260005b818110156105985785810183015185820160400152820161057c565b506000604082860101526040601f19601f8301168501019250505092915050565b80356001600160a01b03811681146105d057600080fd5b919050565b600080604083850312156105e857600080fd5b6105f1836105b9565b946020939093013593505050565b60008060006060848603121561061457600080fd5b61061d846105b9565b925061062b602085016105b9565b9150604084013590509250925092565b60006020828403121561064d57600080fd5b610656826105b9565b9392505050565b6000806040838503121561067057600080fd5b610679836105b9565b9150610687602084016105b9565b90509250929050565b600181811c908216806106a457607f821691505b6020821081036106c457634e487b7160e01b600052602260045260246000fd5b50919050565b8082018082111561023057634e487b7160e01b600052601160045260246000fdfea2646970667358221220852ff68806f8e96dbffd1e458d29bd065afa6103e9f7f362fcbfb4b8b2e1959264736f6c63430008140033",
  };

  const ERC20_ABI = [
    "function transfer(address to, uint256 amount) external returns (bool)"
  ];

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
    console.log("\nGenerating ERC20 contract deployment...");
    // a. Calculate contract address (using the next nonce of the main account)
    const erc20Address = ethers.getCreateAddress({
      from: CONFIG.senderAddress,
      nonce: mainNonce
    });
    CONFIG.erc20Address = erc20Address;
    const initialEthTxs = [];
    // b. Deployment transaction (using the current nonce of the main account)
    const deployTx = {
      data: CONFIG.erc20Bytecode,
      nonce: mainNonce++,
      gasLimit: 6000000,
      gasPrice,
      chainId
    };
    initialEthTxs.push(await senderWallet.signTransaction(deployTx));
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
    console.log("\nGenerating ERC20 init transfers...");
    const initErc20Txs = [];

    const erc20Interface = new ethers.Interface(ERC20_ABI);
    for (const receiver of receivers) {
      const data = erc20Interface.encodeFunctionData("transfer", [
        receiver,
        CONFIG.amountToken
      ]);

      const txData = {
        to: erc20Address,
        data,
        nonce: mainNonce++,
        gasLimit: 60000,
        gasPrice,
        chainId
      };

      initErc20Txs.push(await senderWallet.signTransaction(txData));
    }
    fs.writeFileSync(CONFIG.outputErc20File, initErc20Txs.join("\n"));

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
    console.log("\nGenerating secondary ERC20 transfers...");
    const secondaryErc20Txs = [];

    for (const receiver of receivers) {
      const randomReceiver = receivers[Math.floor(Math.random() * receivers.length)];
      const derived = rootNode.derivePath(`m/44'/60'/0'/0/${receivers.indexOf(receiver)}`);
      const senderWallet = new ethers.Wallet(derived.privateKey);

      const data = erc20Interface.encodeFunctionData("transfer", [
        randomReceiver,
        CONFIG.transferAmountToken
      ]);

      const txData = {
        to: erc20Address,
        data,
        nonce: 1, // Nonce of the receiving address starts from 1 (0 used for ENI transfer)
        gasLimit: 60000,
        gasPrice,
        chainId
      };

      secondaryErc20Txs.push(await senderWallet.signTransaction(txData));
    }
    fs.writeFileSync(CONFIG.transferErc20File, secondaryErc20Txs.join("\n"));

    // ----------------------------------
    // Output summary
    // ----------------------------------
    console.log("\nAll transactions generated successfully!");
    console.log(`- Initial ENI transactions: ${initialEthTxs.length} (saved to ${CONFIG.outputFile})`);
    console.log(`- ERC20 deployment + initial transfers: ${initErc20Txs.length} (saved to ${CONFIG.outputErc20File})`);
    console.log(`- Secondary ENI transfers: ${secondaryEthTxs.length} (saved to ${CONFIG.transferFile})`);
    console.log(`- Secondary ERC20 transfers: ${secondaryErc20Txs.length} (saved to ${CONFIG.transferErc20File})`);
  }

  main().catch(console.error);