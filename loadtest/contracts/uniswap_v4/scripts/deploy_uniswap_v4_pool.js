const hre = require("hardhat");
const { ethers } = require("hardhat");

async function main() {
  try {
    // Get the signer
    const [deployer] = await ethers.getSigners();
    console.log("Deploying Uniswap V4 Pool with the account:", deployer.address);
    console.log("Expected address: 0xF87A299e6bC7bEba58dbBe5a5Aa21d49bCD16D52");
    console.log("Account balance:", ethers.utils.formatEther(await deployer.getBalance()), "ETH");

    // Get token addresses from environment variables
    const token0 = process.env.TOKEN0_ADDRESS;
    const token1 = process.env.TOKEN1_ADDRESS;

    if (!token0 || !token1) {
      console.error("Please ensure TOKEN0_ADDRESS and TOKEN1_ADDRESS environment variables are set");
      process.exit(1);
    }

    console.log("\nToken addresses:");
    console.log("Token0:", token0);
    console.log("Token1:", token1);

    // Deploy Uniswap V4 Pool
    console.log("\nDeploying Uniswap V4 Pool...");
    const UniswapV4Pool = await ethers.getContractFactory("UniswapV4Pool");
    const pool = await UniswapV4Pool.deploy(token0, token1);
    await pool.deployed();
    console.log("Uniswap V4 Pool deployed to:", pool.address);

    // Verify pool state
    const reserves = await pool.getReserves();
    console.log("\nInitial Pool Reserves:");
    console.log("Reserve0:", ethers.utils.formatEther(reserves[0]), "Token0");
    console.log("Reserve1:", ethers.utils.formatEther(reserves[1]), "Token1");

    return pool.address;
  } catch (error) {
    console.error("Deployment failed:", error);
    process.exit(1);
  }
}

main()
  .then((address) => {
    console.log("\nDeployment successful. Pool address:", address);
    process.exit(0);
  })
  .catch((error) => {
    console.error("Deployment failed:", error);
    process.exit(1);
  }); 