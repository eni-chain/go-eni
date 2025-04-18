// SPDX-License-Identifier: GPL-3.0

pragma solidity >= 0.8.0;

import "./common.sol";

contract ValidatorManager{
    //todo: add event and emit for every external method

    //administrator address
    address _admin;

    //current consensus node set
    address[consensusSize] _consensusSet;

    //For traversal and retrieval, because the mapping type cannot be traversed
    address[] _validatorNodes;

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

    //validator info
    struct validator{
        address operator;   //operator address, validator node's owner, also the address to receive the block reward
        address node;       //node address, for consensus
        address agent;      //After being authorized by the operator, the agent can perform operator functions
        bytes  pubKey;      //validator node' public key, ed25519 type, used to verify data such as random, malicious votes, and duplicate proposals
        uint256 amount;      //validator pledge amount
        string name;        //validator name
        string description; //validator description
        uint256 enterTime;  //time to be a validator
        bool isJail;        //current validator is jailed
        uint256 expired;    //expired time of jail
    }

    //operator addr=>validator info
    mapping (address=>validator) _infos;

    //node addr=>operator addr
    mapping (address=>address) _node2operator;

    //agent addr => operator
    mapping (address=>address) _agent2operator;

    //validator name=>operator addr
    mapping (string=>address) _names;

    modifier onlyHub() {
        require(msg.sender == HUB_ADDR, "the message sender must be hub contract");
        _;
    }

        modifier onlyVrf() {
        require(msg.sender == VRF_ADDR, "the message sender must be vrf contract");
        _;
    }

    modifier onlyAdmin() {
        require(msg.sender == _admin, "The message sender must be administrator");
        _;
    }

    function init() external {
        _admin = ADMIN_ADDR;
    }

    function updateAdmin(address admin) external onlyAdmin {
        _admin = admin;
    }

    function getAdmin() external  returns (address){
        return _admin;
    }

    function getPubKey(address node) external returns (bytes memory){
        address ope = _node2operator[node];
        if(ope != address(0) ){
            return _infos[ope].pubKey;
        }

        return bytes("");
    }

    function getNodeAddrAndPubKey(address operator) external returns (address, bytes memory){
        validator storage a = _infos[operator];
        require(a.amount > 0, "Operator and validator not exist");
        return (a.node, a.pubKey);
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

    function getValidatorSet() external returns (address[] memory){
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
        require(_infos[operator].amount == 0, "validator already exist");

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
        _agent2operator[agent] = operator;
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

    function getOperatorAndPledgeAmount(address node) external returns (address, uint256) {
        address oper = _node2operator[node];
        if(oper != address(0) ){
            return (oper, _infos[oper].amount);
        }

        return (address(0), 0);
    }
}
