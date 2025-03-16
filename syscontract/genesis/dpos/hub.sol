// SPDX-License-Identifier: GPL-3.0

pragma solidity >= 0.8.0;

import "./common.sol";

contract Hub {
    //todo: add event and commit for methods

    //uint256 constant  base_ratio_denominator = 100000;
    //uint256 constant block_reward_base_numerator = 20000;
    uint256 constant BASE_RATIO_DENOMINATOR = 100000;
    uint256 constant BLOCK_REWARD_BASE_NUMERATOR = 20000;
    uint256 constant PER_COIN_INCREASE_NUMERATOR = 1;

    //administrator address
    address _admin;

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

    modifier onlyAdmin() {
        require(msg.sender == _admin, "The message sender must be administrator");
        _;
    }

    function initAdmin() external {
        require(msg.sender == ADMIN_ADDR, "The message sender must be administrator");
        _admin = msg.sender;
    }

    function updateAdmin(address admin) external onlyAdmin {
        _admin = admin;
    }

    function getAdmin() external  returns (address){
        return _admin;
    }

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

        delete _applicants[operator];
    }

    function blockReward(address node) external returns (uint256) {
        uint256 pledgeAmount = IValidatorManager(VALIDATOR_MANAGER_ADDR).getPledgeAmount(node);
        uint256 reward = calculateReward(pledgeAmount);
        return reward;
    }

    function calculateReward(uint256 pledge) internal returns (uint256){
        //Reward algorithm: base * { 1 + (pledgeAmount * increasePerCoin)}
        //return (BLOCK_REWARD_BASE_NUMERATOR/BASE_RATIO_DENOMINATOR) * (1 + (pledge * PER_COIN_INCREASE_NUMERATOR/BASE_RATIO_DENOMINATOR));
        //return BLOCK_REWARD_BASE_NUMERATOR * (1*BASE_RATIO_DENOMINATOR + (pledge*BASE_RATIO_DENOMINATOR * PER_COIN_INCREASE_NUMERATOR))/BASE_RATIO_DENOMINATOR;
        return 100000000;
    }
}
