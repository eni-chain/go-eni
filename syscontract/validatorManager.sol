// SPDX-License-Identifier: GPL-3.0

pragma solidity >= 0.8.0;

import "./common.sol";

contract ValidatorManager{
    //当前共识集合
    address[consensusSize] consensusNodes;

    //用于遍历和检索，因为mapping类型无法遍历
    address[] nodes;

    //验证者全量信息
    struct validator{
        address operator;   //操作者地址，用于操作验证者的账户，也是奖励接收地址
        address node;       //节点地址，用于共识出块
        address agent;      //被operator授权后，可代理执行operator职能，可不用
        bytes  pubKey;      //验证者节点公钥，用于校验验证者生成的随机数、恶意投票、重复提案等数据
        uint256 amount;      //验证者自质押额，用于计算区块奖励分配
        string name;        //验证者明文昵称
        string description; //验证者信息介绍
        uint256 enterTime;  //记录成为验证者时的区块号
        bool isJail;        //当前验证者是否被监禁
        uint256 expired;    //监禁截止日期
    }

    //操作者地址=>验证者信息，根据操作者，可查找到验证者全部信息
    mapping (address=>validator) infos;

    //共识节点地址=>操作者地址，eni系统发现共识节点作恶后，需要通过节点地址找到验证者
    mapping (address=>address) node2operator;

    //通过代理人地址可以找到操作者地址，用于代理人权限校验
    mapping (address=>address) agent2perator;

    //name=>operator, 通过验证者昵称可以找到验证者信息，用于浏览器等应用单查找验证者信息
    mapping (string=>address) names;

    modifier onlyHub() {
        require(msg.sender == HUB_ADDR, "the message sender must be hub contract");
        _;
    }

        modifier onlyVrf() {
        require(msg.sender == VRF_ADDR, "the message sender must be vrf contract");
        _;
    }

    function getPubkey(address node) external returns (bytes memory){
        address oper = node2operator[node];
        if(oper != address(0) ){
            return infos[oper].pubKey;
        }

        return bytes("");
    }

    function getValidatorSet() external  returns (address[] memory){
        return nodes;
    }

    function addValidator(
        address operator,
        address node,
        address agent,
        uint256 amount,
        uint256 enterTime,
        string calldata name,
        string calldata description,
        bytes  calldata pubKey
    ) external onlyHub {
        require(amount >= MIN_PLEDGE_AMOUNT, "The transfer amount is less than the minimum pledge amount!");

        validator storage v = infos[operator];
        v.operator = operator;
        v.node = node;
        v.agent = agent;
        v.pubKey = pubKey;
        v.amount = amount;
        v.enterTime = enterTime;
        v.name = name;
        v.description = description;
        v.isJail = false;
        v.expired = 0;

        nodes.push(node);
        names[name] = operator;
        node2operator[node] = operator;
        agent2perator[agent] = operator;
    }

    function undateConsensus(address[] calldata nodes)external onlyVrf {
        require(nodes.length <= consensusSize, "The number of consensuses exceeds the maximum limit");

        delete consensusNodes;
        for(uint i = 0; i < nodes.length; ++i){
            consensusNodes[i] = nodes[i];
        }
    }

    function getPledgeAmount(address node) external returns (uint256) {
        address oper = node2operator[node];
        if(oper != address(0) ){
            return infos[oper].amount;
        }

        return 0;
    }

}
