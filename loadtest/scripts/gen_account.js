const { ethers } = require("ethers");
const bip39 = require("bip39");

// 配置参数
const CONFIG = {
  mnemonic: "party two quit over jaguar carry episode naive machine nothing borrow sell", // 替换为实际助记词
  senderPrivKey: "0x57acb95d82739866a5c29e40b0aa2590742ae50425b7dd5b5d279a986370189e",
  senderAddress: "0xF87A299e6bC7bEba58dbBe5a5Aa21d49bCD16D52",
  numAddresses: 100,
  amountEth: 1,
  rpcUrl: "http://localhost:8545", // 推荐使用公共RPC
  outputFile: "transactions.txt"
};

async function main() {
  // 初始化HD钱包
  const hdNode = ethers.utils.HDNode.fromMnemonic(bip39.mnemonicToSeedSync(CONFIG.mnemonic));
  
  // 生成接收地址
  const receivers = Array.from({ length: CONFIG.numAddresses }, (_, i) => {
    const derived = hdNode.derivePath(`m/44'/60'/0'/0/${i}`);
    return derived.address;
  });

  // 连接Provider获取实时nonce
  const provider = new ethers.providers.JsonRpcProvider(CONFIG.rpcUrl);
  const wallet = new ethers.Wallet(CONFIG.senderPrivKey, provider);
  let nonce = await wallet.getTransactionCount();

  // 构造签名交易
  const rawTxs = [];
  for (const receiver of receivers) {
    const tx = {
      to: receiver,
      value: ethers.utils.parseEther(CONFIG.amountEth.toString()),
      nonce: nonce++,
      gasLimit: 21000,
      chainId: 6912115
    };
    
    // 动态获取gasPrice
    tx.gasPrice =  1000000000; // await provider.getGasPrice();
    
    // 签名交易
    const signedTx = await wallet.signTransaction(tx);
    rawTxs.push(signedTx);
  }

  // 写入文件
  const fs = require("fs");
  fs.writeFileSync(CONFIG.outputFile, rawTxs.join("\n"));
  console.log(`生成 ${rawTxs.length} 笔交易至 ${CONFIG.outputFile}`);
}

main().catch(console.error);