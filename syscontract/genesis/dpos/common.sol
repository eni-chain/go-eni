// SPDX-License-Identifier: GPL-3.0

pragma solidity >= 0.8.0;

//system contract parameters
uint constant CONSENSUS_SIZE = 40;
uint constant MIN_PLEDGE_AMOUNT = 10000000000000000000000; //wei

//precompiled contract address
uint constant ED25519_VERIFY_PRECOMPILED = 0xa1;
uint constant LOCAL_NODE_LOG_PRECOMPILED = 0xa2;

//initial address of the system administrator
address constant INIT_ADMIN_ADDR = 0x3140aedbf686A3150060Cb946893b0598b266f5C;

//system contract address
address constant HUB_ADDR = 0x0000000000000000000000000000000000001001;
address constant VALIDATOR_MANAGER_ADDR = 0x0000000000000000000000000000000000001002;
address constant VRF_ADDR = 0x0000000000000000000000000000000000001003;
address constant VOTER_MANAGER_ADDR = 0x0000000000000000000000000000000000001004;
address constant SLASH_ADDR = 0x0000000000000000000000000000000000001005;

contract Common {
    modifier onlyCoinbase() {
        require(msg.sender == block.coinbase, "the message sender must be the block producer");
        _;
    }

    modifier onlyZeroGasPrice() {
        require(tx.gasprice == 0, "gasprice is not zero");
        _;
    }

    modifier onlyHub() {
        require(msg.sender == HUB_ADDR, "the message sender must be hub contract");
        _;
    }

    modifier onlyValidatorManager() {
        require(msg.sender == VALIDATOR_MANAGER_ADDR, "the message sender must be validator manager contract");
        _;
    }

    modifier onlyVrf() {
        require(msg.sender == VRF_ADDR, "the message sender must be vrf contract");
        _;
    }

    modifier onlyVoteManager() {
        require(msg.sender == VRF_ADDR, "the message sender must be vote manager contract");
        _;
    }

    modifier onlySlash() {
        require(msg.sender == SLASH_ADDR, "the message sender must be slash contract");
        _;
    }
}

interface IValidatorManager {
    //validator info
    struct validator{
        address operator;   //operator address, validator node's owner, also the address to receive the block reward
        address node;       //node address, for consensus
        address agent;      //After being authorized by the operator, the agent can perform operator functions
        bytes  pubKey;      //validator node' public key, ed25519 type, used to verify data such as random, malicious votes, and duplicate proposals
        uint256 amount;     //validator pledge amount
        string name;        //validator name
        string description; //validator description
        uint256 applyBlockNumber;  //bock number when applied to be validator
        uint256 passBlockNumber;   //bock number when approved to be validator
        bool isJail;        //current validator is jailed
        uint256 expired;    //expired time of jail
    }

    function getPubKey(address validator) external view returns (bytes memory);

    function getNodeAddrAndPubKey(address operator) external view returns (address, bytes memory);

    function getPubKeysBySequence(address[] calldata nodes) external view returns (bytes[] memory);

    function getDefaultValidatorSet() external view returns (address[] memory);

    function getJoinedValidatorSet() external view returns (address[] memory);

    function getValidatorSet() external view returns (address[] memory);

    function addDefaultValidator(address operator, address node, address agent, uint256 amount, string calldata name, string calldata description, bytes calldata pubKey ) external;

    function addValidator(address operator, address node, address agent, uint256 amount, uint256 applyBlockNumber, string calldata name, string calldata description, bytes calldata pubKey) external;

    function undateConsensus(address[] calldata nodes)external;

    function getPledgeAmount(address node) external view returns (uint256);

    function getOperatorAndPledgeAmount(address node) external view returns (address, uint256);

    function getValidatorInfo(address operator) external view returns (validator memory);
}