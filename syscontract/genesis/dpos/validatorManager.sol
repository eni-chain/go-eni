// SPDX-License-Identifier: GPL-3.0

pragma solidity >= 0.8.0;

import "./common.sol";

contract ValidatorManager{
    //todo: add event and emit

    address[consensusSize] _consensusSet;

    address[] _validatorNodes;

    struct validator{
        address operator;
        address node;
        address agent;
        bytes  pubKey;
        uint256 amount;
        string name;
        string description;
        uint256 enterTime;
        bool isJail;
        uint256 expired;
    }

    mapping (address=>validator) _infos;

    mapping (address=>address) _node2operator;

    mapping (address=>address) _agent2perator;

    mapping (string=>address) _names;

    modifier onlyHub() {
        require(msg.sender == HUB_ADDR, "the message sender must be hub contract");
        _;
    }

        modifier onlyVrf() {
        require(msg.sender == VRF_ADDR, "the message sender must be vrf contract");
        _;
    }

    function getPubKey(address node) external returns (bytes memory){
        address ope = _node2operator[node];
        if(ope != address(0) ){
            return _infos[ope].pubKey;
        }

        return bytes("");
    }

    function getPubKeysBySequence(address[] calldata nodes) external returns (bytes[] memory){
        bytes[] memory pubKeys = new bytes[](nodes.length);

        for(uint i = 0; i < nodes.length; i++){
            address ope = _node2operator[nodes[i]];
            if(ope != address(0) ){
                pubKeys[i] = _infos[ope].pubKey;
            }
        }

        return pubKeys;
    }

    function getValidatorSet() external  returns (address[] memory){
        return _validatorNodes;
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

        validator storage v = _infos[operator];
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

        _validatorNodes.push(node);
        _names[name] = operator;
        _node2operator[node] = operator;
        _agent2perator[agent] = operator;
    }

    function undateConsensus(address[] calldata nodes)external onlyVrf {
        require(nodes.length <= consensusSize, "The number of consensuses exceeds the maximum limit");

        delete _consensusSet;
        for(uint i = 0; i < nodes.length; ++i){
            _consensusSet[i] = nodes[i];
        }
    }

    function getPledgeAmount(address node) external returns (uint256) {
        address oper = _node2operator[node];
        if(oper != address(0) ){
            return _infos[oper].amount;
        }

        return 0;
    }

}
