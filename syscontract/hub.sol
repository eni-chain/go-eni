// SPDX-License-Identifier: GPL-3.0

pragma solidity >= 0.8.0;

import "./common.sol";

contract Hub {
    //todo: 为每个方法添加event

    //uint256 constant  base_ratio_denominator = 100000;
    //uint256 constant block_reward_base_numerator = 20000;
    uint256 constant BASE_RATIO_DENOMINATOR = 100000;
    uint256 constant BLOCK_REWARD_BASE_NUMERATOR = 20000;
    uint256 constant PER_COIN_INCREASE_NUMERATOR = 1;

    //管理员地址
    address _admin;

    //申请者结构
    struct applicant{
        address operator; //操作者地址，用于操作验证者的账户
        address node; //节点地址，用于共识
        address agent; //被operator授权后，可代理执行operator职能；
        bytes  pubKey; //验证者公钥
        uint256 amount;//验证者质押额
        string name; //验证者明文昵称
        string description; //验证者信息介绍
        uint256 enterTime; //申请时日期
    }

    //申请者列表
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

        //todo: 第二期开发，调用投票者管理合约获取指向验证者的所有投票者的质押额
        uint256 reward = calculateReward(pledgeAmount);

    }

    function calculateReward(uint256 pledge) internal returns (uint256){
        //todo: 奖励算法：出块奖励基数 * { 1 + (质押额 * 每币提升点)}
        //return (BLOCK_REWARD_BASE_NUMERATOR/BASE_RATIO_DENOMINATOR) * (1 + (pledge * PER_COIN_INCREASE_NUMERATOR/BASE_RATIO_DENOMINATOR));
        return BLOCK_REWARD_BASE_NUMERATOR * (1*BASE_RATIO_DENOMINATOR + (pledge*BASE_RATIO_DENOMINATOR * PER_COIN_INCREASE_NUMERATOR))/BASE_RATIO_DENOMINATOR;
    }
}

