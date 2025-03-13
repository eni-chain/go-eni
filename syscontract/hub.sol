// SPDX-License-Identifier: GPL-3.0

pragma solidity >= 0.8.0;

import "./common.sol";

contract Hub {
    //todo: add event and emit

    //uint256 constant  base_ratio_denominator = 100000;
    //uint256 constant block_reward_base_numerator = 20000;
    uint256 constant BASE_RATIO_DENOMINATOR = 100000;
    uint256 constant BLOCK_REWARD_BASE_NUMERATOR = 20000;
    uint256 constant PER_COIN_INCREASE_NUMERATOR = 1;

    address _admin;

    struct applicant{
        address operator;
        address node;
        address agent;
        bytes  pubKey;
        uint256 amount;
        string name;
        string description;
        uint256 enterTime;
    }

    mapping (address=>applicant) _applicants;

    modifier onlyAdmin() {
        require(msg.sender == _admin, "the message sender must be administrator");
        _;
    }

    constructor(){
        _admin = ADMIN_ADDR;
    }

    function updateAdmin(address admin) external onlyAdmin {
        _admin = admin;
    }

    function applyForValidator(
        address node,
        address agent,
        string calldata name,
        string calldata description,
        bytes  calldata pubKey
    ) payable external {
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
        applicant storage a = _applicants[msg.sender];

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

        delete _applicants[msg.sender];
    }

    function blockReward(address node) external returns (uint256) {
        uint256 pledgeAmount = IValidatorManager(VALIDATOR_MANAGER_ADDR).getPledgeAmount(node);

        uint256 reward = calculateReward(pledgeAmount);

    }

    function calculateReward(uint256 pledge) internal returns (uint256){
        //return (BLOCK_REWARD_BASE_NUMERATOR/BASE_RATIO_DENOMINATOR) * (1 + (pledge * PER_COIN_INCREASE_NUMERATOR/BASE_RATIO_DENOMINATOR));
        return BLOCK_REWARD_BASE_NUMERATOR * (1*BASE_RATIO_DENOMINATOR + (pledge*BASE_RATIO_DENOMINATOR * PER_COIN_INCREASE_NUMERATOR))/BASE_RATIO_DENOMINATOR;
    }
}

