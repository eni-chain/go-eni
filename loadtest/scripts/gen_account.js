const { ethers } = require("ethers");
  const bip39 = require("bip39");
  const fs = require("fs");

  // 配置参数
  const CONFIG = {
    mnemonic: "party two quit over jaguar carry episode naive machine nothing borrow sell", // 替换为实际助记词
    senderPrivKey: "0x57acb95d82739866a5c29e40b0aa2590742ae50425b7dd5b5d279a986370189e",
    senderAddress: "0xF87A299e6bC7bEba58dbBe5a5Aa21d49bCD16D52",
    numAddresses: 100,
    amountEth: 1,
    rpcUrl: "http://localhost:8545", // 推荐使用公共RPC
    outputFile: "init.txt",
    transferFile: "transfer.txt"
  };

  async function main() {
    // 初始化HD钱包
    const seed = bip39.mnemonicToSeedSync(CONFIG.mnemonic);
    const rootNode = ethers.HDNodeWallet.fromSeed(seed);

    // 生成接收地址
    const receivers = Array.from({ length: CONFIG.numAddresses }, (_, i) => {
      const derived = rootNode.derivePath(`m/44'/60'/0'/0/${i}`);
      return derived.address;
    });

    console.log(receivers);

    // 连接Provider获取实时nonce
    const provider = new ethers.JsonRpcProvider(CONFIG.rpcUrl);
    const wallet = new ethers.Wallet(CONFIG.senderPrivKey, provider);
    let nonce = 0;

    // 构造签名交易
    const rawTxs = [];
    for (const receiver of receivers) {
      const tx = {
        to: receiver,
        value: ethers.parseEther(CONFIG.amountEth.toString()),
        nonce: nonce++,
        gasLimit: 21000,
        chainId: 6912115
      };

      tx.gasPrice = 1000000000; // await provider.getGasPrice();

      const signedTx = await wallet.signTransaction(tx);
      rawTxs.push(signedTx);
    }

    fs.writeFileSync(CONFIG.outputFile, rawTxs.join("\n"));
    console.log(`生成 ${rawTxs.length} 笔初始账户余额交易至 ${CONFIG.outputFile}`);

    // 从接收地址再转账0.1个ETH到随机地址
    const transferTxs = [];
    for (const receiver of receivers) {
      const randomReceiver = receivers[Math.floor(Math.random() * receivers.length)];
      const derived = rootNode.derivePath(`m/44'/60'/0'/0/${receivers.indexOf(receiver)}`);
      const senderWallet = new ethers.Wallet(derived.privateKey, provider);

      const tx = {
        to: randomReceiver,
        value: ethers.parseEther("0.1"),
        nonce: 0,
        gasLimit: 21000,
        chainId: 6912115
      };

      tx.gasPrice = 1000000000; // await provider.getGasPrice();

      const signedTx = await senderWallet.signTransaction(tx);
      transferTxs.push(signedTx);
    }

    fs.writeFileSync(CONFIG.transferFile, transferTxs.join("\n"));
    console.log(`生成 ${transferTxs.length} 笔转账交易至 ${CONFIG.transferFile}`);
  }

  main().catch(console.error);