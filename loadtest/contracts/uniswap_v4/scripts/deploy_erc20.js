const hre = require("hardhat");

async function main() {
  const [deployer] = await hre.ethers.getSigners();
  console.log("Deploying contracts with the account:", deployer.address);

  // Deploy TestToken1
  const TestToken1 = await hre.ethers.getContractFactory("TestToken");
  const token1 = await TestToken1.deploy("Test Token 1", "TEST1");
  await token1.deployed();
  console.log("TestToken1 deployed to:", token1.address);

  // Deploy TestToken2
  const TestToken2 = await hre.ethers.getContractFactory("TestToken");
  const token2 = await TestToken2.deploy("Test Token 2", "TEST2");
  await token2.deployed();
  console.log("TestToken2 deployed to:", token2.address);

  // Return the address of the first token for the pool deployment
  return token1.address;
}

main()
  .then((address) => {
    console.log(address);
    process.exit(0);
  })
  .catch((error) => {
    console.error(error);
    process.exit(1);
  }); 