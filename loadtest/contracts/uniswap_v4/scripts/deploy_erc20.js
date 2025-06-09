const hre = require("hardhat");

async function main() {
  try {
    const [deployer] = await hre.ethers.getSigners();
    console.log("Deploying contracts with the account:", deployer.address);
    console.log("Expected address: 0xF87A299e6bC7bEba58dbBe5a5Aa21d49bCD16D52");
    console.log("Account balance:", hre.ethers.utils.formatEther(await deployer.getBalance()), "ETH");

    // Deploy TestToken1
    console.log("\nDeploying TestToken1...");
    const TestToken1 = await hre.ethers.getContractFactory("TestToken");
    const token1 = await TestToken1.deploy("Test Token 1", "TEST1");
    await token1.deployed();
    console.log("TestToken1 deployed to:", token1.address);

    // Deploy TestToken2
    console.log("\nDeploying TestToken2...");
    const TestToken2 = await hre.ethers.getContractFactory("TestToken");
    const token2 = await TestToken2.deploy("Test Token 2", "TEST2");
    await token2.deployed();
    console.log("TestToken2 deployed to:", token2.address);

    // Verify token balances
    const balance1 = await token1.balanceOf(deployer.address);
    const balance2 = await token2.balanceOf(deployer.address);
    console.log("\nToken Balances:");
    console.log("TestToken1 balance:", hre.ethers.utils.formatEther(balance1), "TEST1");
    console.log("TestToken2 balance:", hre.ethers.utils.formatEther(balance2), "TEST2");

    // Return the address of the first token for the pool deployment
    return token1.address;
  } catch (error) {
    console.error("Deployment failed:", error);
    process.exit(1);
  }
}

main()
  .then((address) => {
    console.log("\nDeployment successful. Token1 address:", address);
    process.exit(0);
  })
  .catch((error) => {
    console.error("Deployment failed:", error);
    process.exit(1);
  }); 