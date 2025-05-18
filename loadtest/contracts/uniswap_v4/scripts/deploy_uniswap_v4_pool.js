const hre = require("hardhat");

async function main() {
  const [deployer] = await hre.ethers.getSigners();
  console.log("Deploying contracts with the account:", deployer.address);

  // Get token addresses from command line arguments
  const token0Address = process.argv[process.argv.indexOf("--token0") + 1];
  const token1Address = process.argv[process.argv.indexOf("--token1") + 1];

  if (!token0Address || !token1Address) {
    console.error("Please provide token addresses using --token0 and --token1 arguments");
    process.exit(1);
  }

  // Deploy Uniswap V4 Pool
  const UniswapV4Pool = await hre.ethers.getContractFactory("UniswapV4Pool");
  const pool = await UniswapV4Pool.deploy(token0Address, token1Address);
  await pool.deployed();

  console.log("Uniswap V4 Pool deployed to:", pool.address);
  return pool.address;
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  }); 