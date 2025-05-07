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
    function getPubKey(address validator) external returns (bytes memory);

    function getNodeAddrAndPubKey(address operator) external returns (address, bytes memory);

    function getPubKeysBySequence(address[] calldata nodes) external returns (bytes[] memory);

    function getValidatorSet() external  returns (address[] memory);

    function addValidator(address operator, address node, address agent, uint256 amount, uint256 enterTime, string calldata name, string calldata description, bytes  calldata pubKey) external;

    function undateConsensus(address[] calldata nodes)external;

    function getPledgeAmount(address node) external returns (uint256);

    function getOperatorAndPledgeAmount(address node) external returns (address, uint256);
}