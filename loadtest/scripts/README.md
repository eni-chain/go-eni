## 1. How to Generate Offline Transfer Transactions

Run the following command to generate offline transfer transactions:

```bash
node genEniTx.js
```

This will generate offline transfer transactions, including a transaction from:

```
0xF87A299e6bC7bEba58dbBe5a5Aa21d49bCD16D52
```

to a new account `1ENI`, stored in `init.txt`. These transactions are serialized. These new accounts are derived from the mnemonic:

```
mnemonic: "party two quit over jaguar carry episode naive machine nothing borrow sell" // Replace with the actual mnemonic
```

You can also import this mnemonic into Metamask to view the corresponding accounts.

After generating `init.txt`, the script will generate transactions from the new accounts transferring `0.1ENI` to random addresses, stored in `transfer.txt`. These transactions are parallel.

To set the number of transactions to generate, modify the following in `gen_account.js`:

```javascript
numAddresses: 10000
```

Change it to the desired value.

## 2. How to Send Offline Transactions

In the current folder, there is a `send.sh` script. To send offline transactions, specify the offline transaction file name. For example:

```bash
./send.sh ./init.txt
```


```
./build/loadtest -tx ./loadtest/scripts/init.txt
```