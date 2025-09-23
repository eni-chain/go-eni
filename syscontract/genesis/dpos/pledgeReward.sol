// SPDX-License-Identifier: GPL-3.0

pragma solidity >= 0.8.0;

import "./common.sol";
import "./localLog.sol";
import "./delegateCallBase.sol";

contract StakeManager is DelegateCallBase, Common {
    uint internal constant LOCK_TIME_90 = 90 days;
    uint internal constant LOCK_TIME_180 = 180 days;
    uint internal constant LOCK_TIME_360 = 360 days;
    uint internal constant LOCK_TIME_720 = 720 days;
    uint internal constant LOCK_TIME_1440 = 1440 days;

    uint internal constant powerMultiple1 = 1;
    uint internal constant powerMultiple2 = 2;
    uint internal constant powerMultiple3 = 3;
    uint internal constant powerMultiple4 = 4;
    uint internal constant powerMultiple5 = 5;
    
    uint internal constant noContinue = 1;
    uint internal constant continueStake = 2;
    uint internal constant compoundInterest = 3;

    uint256 internal coolingPeriod = 14 days;

    uint256 internal totalRewardPerBlock = 600000000000000000;
   
    uint256 internal rewardPerShare;
    uint256 internal lastRewardBlock;
    uint256 internal totalStaked;
    
    struct OrderId {
        address shareholder;
        uint256 sequence;
    }

    struct Order {
      uint256 amount;
      uint256 enterTime;
      uint256 lockPeriod;
      uint256 rewardDebt;
      uint256 unclaimedReward;
      uint256 continueFlag;
      uint256 coolingExpired;
      address[] validators;
      int256 currentValidatorIdx;
    }

    struct Voter {
        Order[] orders;
    }

    struct Validator {
        Order order;
        uint256 poll;
        OrderId[] voters;
        uint256 exitExpired;
        bool frozen;
    }



    OrderId[][] internal OrderSequences; 
    

    mapping(address => Voter) internal voters;   
    

    mapping (address => Validator) internal Validators;


    function updatePool() internal {

        if (block.number <= lastRewardBlock){
            return;
        }


        if (totalStaked == 0) {
            lastRewardBlock = block.number;
            return;
        }


        uint256 blocksPassed = block.number - lastRewardBlock;


        uint256 totalReward = blocksPassed * totalRewardPerBlock;


        uint256 rewardPerShare = totalReward / totalStaked;


        rewardPerShare += rewardPerShare;


        lastRewardBlock = block.number;
    }

    //for voter
    function claimReward(uint256 sequence) external returns(string memory){
        Voter storage v = voters[msg.sender];
        require(v.orders.length > 0, "There is no order for current user");
        require(sequence < v.orders.length, "Invalid order sequence");

        Order storage order = v.orders[sequence];
        require(order.amount != 0, "There is no order for the sequence");

        //require(block.timestamp <= order.coolingExpired, "The cooling-off period has not expired");
        //When order.enterTime+order.lockPeriod>=block.timestamp, rewards will be automatically calculated
        //and reinvested according to the renewed pledge flag, and enterTime will be reset. Therefore,
        //if the pledge is not renewed after the expiration date, rewards will no longer be calculated.
        require(order.enterTime + order.lockPeriod < block.timestamp, "The lock-up period expires but the pledge is not renewed");

        //No validator is specified, no reward is calculated
        require(order.currentValidatorIdx >= 0, "There is no validator for current order");
        address addr = order.validators[uint256(order.currentValidatorIdx)];
        Validator storage validator = Validators[addr];
        require(validator.order.amount != 0, "The current order does not point to a valid validator");

        order.unclaimedReward += (order.amount * rewardPerShare / 1e18) - order.rewardDebt;
        order.rewardDebt = order.amount * rewardPerShare / 1e18;
    }

    //for validator
    function claimReward() external {
        Validator storage val = Validators[msg.sender];
        require(val.order.amount != 0, "msg.sender is not validator");
        require(val.frozen == false, "validator was frozen");

        Order storage order = val.order;
        require(order.enterTime + order.lockPeriod < block.timestamp, "The lock-up period expires but the pledge is not renewed");

        order.unclaimedReward += (order.amount * rewardPerShare / 1e18) - order.rewardDebt;
        order.rewardDebt = order.amount * rewardPerShare / 1e18;
    }


    function validatorStake(uint256 amount, uint256 lockPeriod, bool continueStake, bool compoundInterest, address node, bytes calldata pubKey) external {
        
    }


    function voteStake(uint256 amount, uint256 lockPeriod, bool continueStake, bool compoundInterest, address[] calldata validator) external{

    }


    function transferStake() external{

    }


    function reStake(uint256 continueFlag) external{

    }
    

    function changeOrderValidator(uint256 sequence, address validator) external{

    }


    function withdrawProfit() external{

    }


    function redemption() external{

    }
}
