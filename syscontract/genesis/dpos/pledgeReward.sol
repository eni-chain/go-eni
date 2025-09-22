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
