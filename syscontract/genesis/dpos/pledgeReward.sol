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

    uint internal constant POW_MULTI_ERR = 0;
    uint internal constant POW_MULTI1 = 1;
    uint internal constant POW_MULTI2 = 2;
    uint internal constant POW_MULTI3 = 3;
    uint internal constant POW_MULTI4 = 4;
    uint internal constant POW_MULTI5 = 5;

    uint internal constant NO_CONTINUE = 1;
    uint internal constant CONTINUE_PLEDGE = 2;
    uint internal constant COMPOUND_INTEREST = 3;

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



    OrderId[][] internal OrderSequences_;


    mapping(address => Voter) internal Voters_;


    mapping (address => Validator) internal Validators_;


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


        uint256 rewardPerSharePhased = totalReward / (totalStaked / 1e18);


        rewardPerShare += rewardPerSharePhased;


        lastRewardBlock = block.number;
    }

    function claimRewardBasic(Order storage order) internal {

        order.unclaimedReward += (order.amount * rewardPerShare / 1e18) - order.rewardDebt;
        order.rewardDebt = order.amount * rewardPerShare / 1e18;
    }


    function claimReward(uint256 sequence) external {
        Voter storage v = Voters_[msg.sender];
        require(v.orders.length > 0, "There is no order for current user");
        require(sequence < v.orders.length, "Invalid order sequence");

        Order storage order = v.orders[sequence];
        require(order.amount != 0, "There is no order for the sequence");



        require(order.enterTime + order.lockPeriod < block.timestamp, "The lock-up period expires but the pledge is not renewed");


        require(order.currentValidatorIdx >= 0, "There is no validator for current order");
        address addr = order.validators[uint256(order.currentValidatorIdx)];
        Validator storage validator = Validators_[addr];
        require(validator.order.amount != 0, "The current order does not point to a valid validator");

        claimRewardBasic(order);
    }


    function claimReward() external {
        Validator storage val = Validators_[msg.sender];
        require(val.order.amount != 0, "msg.sender is not validator");
        require(val.frozen == false, "validator was frozen");

        Order storage order = val.order;
        require(order.enterTime + order.lockPeriod < block.timestamp, "The lock-up period expires but the pledge is not renewed");

        claimRewardBasic(order);
    }


    function multiNumber(uint256 lockPeriod) internal pure returns(uint256){
        if(lockPeriod == LOCK_TIME_90){
            return POW_MULTI1;
        }else if(lockPeriod == LOCK_TIME_180){
            return POW_MULTI2;
        }else if(lockPeriod == LOCK_TIME_360){
            return POW_MULTI3;
        }else if(lockPeriod == LOCK_TIME_720){
            return POW_MULTI4;
        }else if(lockPeriod == LOCK_TIME_1440){
            return POW_MULTI5;
        }else{
            return POW_MULTI_ERR;
        }
    }


    function continueFlagValid(uint256 continueFlag) internal pure returns(bool){
        if((continueFlag == NO_CONTINUE) ||(continueFlag == CONTINUE_PLEDGE)|| (continueFlag == COMPOUND_INTEREST)){
            return true;
        }
        return false;
    }


    function fillOrder(Order storage order, uint256 lockPeriod, uint256 continueFlag) internal{
        order.amount = msg.value;
        order.enterTime = block.timestamp;
        order.lockPeriod = lockPeriod;
        order.rewardDebt = order.amount * rewardPerShare / 1e18;
        order.unclaimedReward = 0;
        order.continueFlag = continueFlag;
        order.coolingExpired = block.timestamp + coolingPeriod;
    }


    function validatorStake(uint256 lockPeriod, uint256 continueFlag, address node, bytes calldata pubKey) payable external {
        uint256 multi = multiNumber(lockPeriod);
        require(multi != POW_MULTI_ERR, "Lock period error!");
        require(msg.value >= MIN_PLEDGE_AMOUNT, "The transfer amount is less than validator minimum pledge amount!");
        require(continueFlagValid(continueFlag), "Continue flag error!");

        updatePool();
        Validator storage vali = Validators_[msg.sender];
        require(vali.order.amount == 0, "Validator alread exist.");

        Order storage order = vali.order;
        fillOrder(order, lockPeriod, continueFlag);

        uint256 pledgePower = order.amount * multi;
        totalStaked += pledgePower;

        IValidatorManager(VALIDATOR_MANAGER_ADDR).addValidator(
            msg.sender,
            node,
            msg.sender,
            msg.value,
            block.number,
            "",
            "",
            pubKey
        );


        OrderId memory orderId;
        orderId.shareholder = msg.sender;
        OrderId[] storage orderIds = OrderSequences_[multi];
        orderIds.push(orderId);
    }


    function voteStake(uint256 lockPeriod, uint256 continueFlag, address[] calldata validator) payable external{
        uint256 multi = multiNumber(lockPeriod);
        require(multi != POW_MULTI_ERR, "Lock period error!");
        require(msg.value >= 1e18, "The transfer amount is less than 1 ENI!");
        require(validator.length > 0, "Validator address list is null!");


        Validator storage vali = Validators_[validator[0]];
        require(vali.order.amount != 0, "Validator not exist.");

        updatePool();

        Voter storage voter = Voters_[msg.sender];
        Order storage order = voter.orders[voter.orders.length];

        fillOrder(order, lockPeriod, continueFlag);
        order.validators = validator;
        order.currentValidatorIdx = 0;

        uint256 pledgePower = order.amount * multi;
        totalStaked += pledgePower;


        vali.poll += msg.value;


        OrderId memory orderId;
        orderId.shareholder = msg.sender;
        orderId.sequence = voter.orders.length;
        vali.voters.push(orderId);


        OrderId[] storage orderIds = OrderSequences_[multi];
        orderIds.push(orderId);
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