// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Test, console2} from "forge-std/Test.sol";
import {NativeEniTokensERC20} from "../src/NativeEniTokensERC20.sol";
import {IBank} from "../src/precompiles/IBank.sol";

address constant BANK_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000001001;

contract MockBank {
    mapping(address => uint256) balances;

    // mocking functions
    function setBalances(address[] memory addressesToFund) public {
        for (uint256 i = 0; i < addressesToFund.length; i++) {
            balances[addressesToFund[i]] = 1000;
        }
    }

    // subset of IBank functions
    function balance(address account, string memory denom) public view returns (uint256) {
        require(keccak256(abi.encodePacked(denom)) == keccak256(abi.encodePacked("ueni")), "MockBank: denom not supported");
        return balances[account];
    }

    function send(
        address fromAddress,
        address toAddress,
        string memory denom,
        uint256 amount
    ) external returns (bool success) {
        require(keccak256(abi.encodePacked(denom)) == keccak256(abi.encodePacked("ueni")), "MockBank: denom not supported");
        balances[fromAddress] -= amount;
        balances[toAddress] += amount;
        return true;
    }
}

contract NativeEniTokensERC20Test is Test {

    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);

    NativeEniTokensERC20 eniERC20;
    address alice;
    address bob;

    function setUp() public {
        alice = makeAddr("alice");
        bob = makeAddr("bob");
        eniERC20 = new NativeEniTokensERC20("ueni", "ENI", "ENISYMBOL", 6);

        MockBank mockBank = new MockBank();
        vm.etch(BANK_PRECOMPILE_ADDRESS, address(mockBank).code);
        address[] memory addressesToFund = new address[](2);
        addressesToFund[0] = alice;
        addressesToFund[1] = bob;
        MockBank(BANK_PRECOMPILE_ADDRESS).setBalances(addressesToFund);
    }

    function testName() public {
        assertEq(eniERC20.name(), "ENI");
    }

    function testSymbol() public {
        assertEq(eniERC20.symbol(), "ENISYMBOL");
    }

    function testBalanceOf() public {
        vm.mockCall(BANK_PRECOMPILE_ADDRESS, abi.encodeWithSelector(IBank.balance.selector, address(this), "ueni"), abi.encode(123));
        assertEq(eniERC20.balanceOf(address(this)), 123);
    }

    function testDecimals() public {
        assertEq(eniERC20.decimals(), 6);
    }

    function testTotalSupply() public {
        vm.mockCall(BANK_PRECOMPILE_ADDRESS, abi.encodeWithSelector(IBank.supply.selector, "ueni"), abi.encode(123));
        assertEq(eniERC20.totalSupply(), 123);
    }

    function testTransfer() public {
        vm.expectEmit();
        emit Transfer(alice, bob, 123);

        vm.startPrank(alice);
        bool success = eniERC20.transfer(bob, 123);
        vm.stopPrank();

        assertEq(success, true);
        assertEq(eniERC20.balanceOf(alice), 1000 - 123);
        assertEq(eniERC20.balanceOf(bob), 1000 + 123);
    }

    function testApprovals() public {
        // Alice approves Bob to spend 200 tokens on her behalf
        vm.expectEmit();
        emit Approval(alice, bob, 200);

        vm.startPrank(alice);
        bool approvalSuccess = eniERC20.approve(bob, 200);
        vm.stopPrank();

        assertEq(approvalSuccess, true);
        assertEq(eniERC20.allowance(alice, bob), 200);
    }

    function testTransferFrom() public {
        // expect fail because no approval was given
        vm.startPrank(bob);
        vm.expectRevert();
        eniERC20.transferFrom(alice, bob, 150);
        vm.stopPrank();

        // alice to approve bob to spend tokens on her behalf
        vm.startPrank(alice);
        eniERC20.approve(bob, 200);
        vm.stopPrank();

        vm.startPrank(bob);
        vm.expectEmit();
        emit Transfer(alice, bob, 150);
        bool transferFromSuccess = eniERC20.transferFrom(alice, bob, 150);
        vm.stopPrank();

        assertEq(transferFromSuccess, true);
        assertEq(eniERC20.balanceOf(alice), 1000 - 150);
        assertEq(eniERC20.balanceOf(bob), 1000 + 150);
        assertEq(eniERC20.allowance(alice, bob), 50); // Remaining allowance after the transfer
    }
}