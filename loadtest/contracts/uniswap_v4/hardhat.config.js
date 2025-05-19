require("@nomicfoundation/hardhat-toolbox");
require("@nomiclabs/hardhat-ethers");

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: "0.8.19",
  networks: {
    localhost: {
      url: "http://127.0.0.1:8545",
      chainId: 6912115,
      accounts: ["0x57acb95d82739866a5c29e40b0aa2590742ae50425b7dd5b5d279a986370189e"],
      timeout: 60000,
      gas: 2100000,
      gasPrice: 8000000000,
      allowUnlimitedContractSize: true,
      loggingEnabled: true
    }
  },
  paths: {
    sources: "./contracts",
    tests: "./test",
    cache: "./cache",
    artifacts: "./artifacts"
  },
  mocha: {
    timeout: 40000
  },
  ethers: {
    version: "5.7.2"
  },
  // Disable ENS resolution for local network
  namedAccounts: {
    deployer: 0
  }
}; 