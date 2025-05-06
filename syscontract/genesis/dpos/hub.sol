// SPDX-License-Identifier: GPL-3.0

pragma solidity >= 0.8.0;

import "./common.sol";

contract Hub is DelegateCallBase {

    uint256 constant ratioDeno = 100000;
    uint256 constant ratioNumer = 20000;
    uint256 constant increasePerCoin = 1;
    uint256 constant weiPerCoin = 1000000000000000000;

    //validator apply info
    struct applicant{
        address operator; //operator address, validator node's owner
        address node; //node address, for consensus
        address agent; //After being authorized by the operator, the agent can perform operator functions
        bytes  pubKey; //validator node' public key, ed25519 type
        uint256 amount;//validator pledge amount
        string name; //validator name
        string description; //validator description
        uint256 enterTime; //time of application
    }

    //List of applicants
    mapping (address=>applicant) _applicants;

    event ApplyForValidator(string indexed name, address indexed operator, address indexed node, bytes pubKey, uint256 pledge);

    event AuditPass(address indexed admin, string indexed name, address indexed operator, address node, bytes pubKey, uint256 pledge);

    event BlockReward(address indexed proposer, uint256 pledge, uint256 reward);

    function applyForValidator(
        address node,
        address agent,
        string calldata name,
        string calldata description,
        bytes  calldata pubKey
    ) payable external {
        require(msg.value >= MIN_PLEDGE_AMOUNT, "The transfer amount is less than the minimum pledge amount!");
        require(_applicants[msg.sender].amount == 0, "applicant already exsit");

        applicant storage a = _applicants[msg.sender];
        a.operator = msg.sender;
        a.node = node;
        a.agent = agent;
        a.pubKey = pubKey;
        a.amount = msg.value;
        a.name = name;
        a.description = description;
        a.enterTime = block.timestamp;

        llog(DEBUG, abi.encodePacked(name, " applyForValidator, operator: ", H(msg.sender), ", node:", H(node), ", plege amount: ", S(msg.value)));
        emit ApplyForValidator(name, msg.sender, node, pubKey, msg.value);
    }

    function auditPass(address operator) external onlyAdmin {
        applicant storage a = _applicants[operator];
        require(a.amount > 0, "applicant not exists");

        IValidatorManager(VALIDATOR_MANAGER_ADDR).addValidator(
            a.operator,
            a.node,
            a.agent,
            a.amount,
            a.enterTime,
            a.name,
            a.description,
            a.pubKey
        );

        llog(DEBUG, abi.encodePacked("auditPass, validator name:", a.name, ", operator:", H(a.operator), ", admin:", H(msg.sender),  ", pledge amount: ", S(a.amount)));
        emit AuditPass(msg.sender, a.name, a.operator, a.node, a.pubKey, a.amount);

        delete _applicants[operator];

    }

    function blockReward(address node) external returns (address, uint256) {
        address operator;
        uint256 pledgeAmount;
        (operator, pledgeAmount) = IValidatorManager(VALIDATOR_MANAGER_ADDR).getOperatorAndPledgeAmount(node);
        uint256 reward = calculateReward(pledgeAmount);

        llog(DEBUG, abi.encodePacked("blockReward, proposer:", H(operator), ", pledge amount: ", S(pledgeAmount), ", reward:", S(reward)));
        emit BlockReward(operator, pledgeAmount, reward);

        return (operator, reward);
    }

    //Reward algorithm: base * { 1 + (pledgeAmount * increasePerCoin)}
    function calculateReward(uint256 pledgeAmount) internal pure returns (uint256){
        require(pledgeAmount != 0, "Pledge amount is 0, maybe dpos not started");

        //convert wei to coin
        uint256 pledge = pledgeAmount/weiPerCoin;
        //llog(DEBUG, abi.encodePacked("calculateReward, pledge amount in wei:", S(pledgeAmount), ", in coin:", S(pledge)));

        //uint256 reward = (ratioNumer/ratioDeno) *(1 + (pledge*(increasePerCoin/ratioDeno)));
        //uint256 reward = (ratioNumer/ratioDeno) *(1*ratioDeno + (pledge*increasePerCoin))/ratioDeno;
        //uint256 reward = ((ratioNumer*(1*ratioDeno + (pledge*increasePerCoin)))*weiPerCoin)/(ratioDeno*ratioDeno);
        uint256 reward = (ratioNumer*(1*ratioDeno + (pledge*increasePerCoin)))*(weiPerCoin/(ratioDeno*ratioDeno));
        //llog(DEBUG, abi.encodePacked("calculateReward, reward amount in wei:", S(reward)));

        return reward;
    }
}
