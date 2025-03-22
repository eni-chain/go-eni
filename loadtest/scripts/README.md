## 1.如何生成离线转账交易

node gen_account.js

即可生成离线转账交易。其中包括

```
0xF87A299e6bC7bEba58dbBe5a5Aa21d49bCD16D52
```

转账给新账号1ENI的交易，存放在init.txt中，这些交易都是串行的。这些新账号都是通过助记词：

```
    mnemonic: "party two quit over jaguar carry episode naive machine nothing borrow sell", // 替换为实际助记词
```

派生出来的。所以你也可以将这个助记词导入Metamask查看对应账号。

生成完成init.txt后，脚本又生成了从新账号转账0.1ENI给随机地址的交易，存放在transfer.txt中，这些交易都是并行的。

如果要设定生成交易的数量，请修改gen_account.js中的

```
numAddresses: 10000
```

改为想要的值即可。

## 2.如何发送离线交易

在当前文件夹有send.sh脚本。要发送离线交易，请指定离线交易文件名即可。比如：

```bash
./send.sh ./init.txt
```
