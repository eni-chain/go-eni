// SPDX-License-Identifier: GPL-3.0

pragma solidity >= 0.8.0;

import "./common.sol";

contract ValidatorManager is DelegateCallBase, LocalLog, Common {
    //todo: add event and emit for every external method

    //current consensus node set
    address[CONSENSUS_SIZE] _consensusSet;

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
        llog(DEBUG, abi.encodePacked("addValidator ", name, " start."));

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
        //llog(DEBUG, abi.encodePacked("addValidator ", name, " succeed."));
    }

    function undateConsensus(address[] calldata nodes)external onlyVrf {
        require(nodes.length <= CONSENSUS_SIZE, "The number of consensuses exceeds the maximum limit");

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
        llog(DEBUG, abi.encodePacked("hub contract call getOperatorAndPledgeAmount, to node:", H(node)));
        address oper = _node2operator[node];
        if(oper != address(0) ){
            return (oper, _infos[oper].amount);
        }
        llog(DEBUG, abi.encodePacked("hub contract call getOperatorAndPledgeAmount, to operator:", H(oper)));
        return (address(0), 0);
    }
}
