// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract UniswapV4Pool is Ownable {
    IERC20 public immutable token0;
    IERC20 public immutable token1;
    
    uint256 public reserve0;
    uint256 public reserve1;
    
    event Swap(address indexed sender, uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, address indexed to);
    event Mint(address indexed sender, uint256 amount0, uint256 amount1);
    event Burn(address indexed sender, uint256 amount0, uint256 amount1, address indexed to);
    
    constructor(address _token0, address _token1) {
        require(_token0 != _token1, "UniswapV4Pool: IDENTICAL_ADDRESSES");
        require(_token0 != address(0) && _token1 != address(0), "UniswapV4Pool: ZERO_ADDRESS");
        
        token0 = IERC20(_token0);
        token1 = IERC20(_token1);
    }
    
    function getReserves() public view returns (uint256 _reserve0, uint256 _reserve1) {
        _reserve0 = reserve0;
        _reserve1 = reserve1;
    }
    
    function mint(uint256 amount0, uint256 amount1) external {
        require(amount0 > 0 && amount1 > 0, "UniswapV4Pool: INSUFFICIENT_INPUT_AMOUNT");
        
        // Transfer tokens from sender
        require(token0.transferFrom(msg.sender, address(this), amount0), "UniswapV4Pool: TRANSFER_FAILED");
        require(token1.transferFrom(msg.sender, address(this), amount1), "UniswapV4Pool: TRANSFER_FAILED");
        
        // Update reserves
        reserve0 += amount0;
        reserve1 += amount1;
        
        emit Mint(msg.sender, amount0, amount1);
    }
    
    function burn(uint256 amount0, uint256 amount1, address to) external {
        require(amount0 > 0 || amount1 > 0, "UniswapV4Pool: INSUFFICIENT_LIQUIDITY_BURNED");
        require(to != address(0), "UniswapV4Pool: INVALID_TO");
        
        // Update reserves
        reserve0 -= amount0;
        reserve1 -= amount1;
        
        // Transfer tokens to recipient
        if (amount0 > 0) require(token0.transfer(to, amount0), "UniswapV4Pool: TRANSFER_FAILED");
        if (amount1 > 0) require(token1.transfer(to, amount1), "UniswapV4Pool: TRANSFER_FAILED");
        
        emit Burn(msg.sender, amount0, amount1, to);
    }
    
    function swap(uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, address to) external {
        require(amount0In > 0 || amount1In > 0, "UniswapV4Pool: INSUFFICIENT_INPUT_AMOUNT");
        require(amount0Out > 0 || amount1Out > 0, "UniswapV4Pool: INSUFFICIENT_OUTPUT_AMOUNT");
        require(to != address(0), "UniswapV4Pool: INVALID_TO");
        
        // Update reserves
        reserve0 = reserve0 + amount0In - amount0Out;
        reserve1 = reserve1 + amount1In - amount1Out;
        
        // Transfer input tokens from sender
        if (amount0In > 0) require(token0.transferFrom(msg.sender, address(this), amount0In), "UniswapV4Pool: TRANSFER_FAILED");
        if (amount1In > 0) require(token1.transferFrom(msg.sender, address(this), amount1In), "UniswapV4Pool: TRANSFER_FAILED");
        
        // Transfer output tokens to recipient
        if (amount0Out > 0) require(token0.transfer(to, amount0Out), "UniswapV4Pool: TRANSFER_FAILED");
        if (amount1Out > 0) require(token1.transfer(to, amount1Out), "UniswapV4Pool: TRANSFER_FAILED");
        
        emit Swap(msg.sender, amount0In, amount1In, amount0Out, amount1Out, to);
    }
} 