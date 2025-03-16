// SPDX-License-Identifier: GPL-3.0


pragma solidity >= 0.8.0;

uint constant consensusSize = 2;

uint constant MIN_PLEDGE_AMOUNT = 10000;

uint constant ED25519_VERIFY_PRECOMPILED = 0xa1;

address constant ADMIN_ADDR = 0x251604eBfD1ddeef1F4f40b8F9Fc425538BE1339;
address constant HUB_ADDR = 0x0000000000000000000000000000000000001001;
address constant VALIDATOR_MANAGER_ADDR = 0x0000000000000000000000000000000000001002;
address constant VRF_ADDR = 0x0000000000000000000000000000000000001003;
address constant VOTER_MANAGER_ADDR = 0x0000000000000000000000000000000000001004;
address constant SLASH_ADDR = 0x0000000000000000000000000000000000001005;


contract Common {
    bool public _alreadyInit = false;


    modifier onlyCoinbase() {
        require(msg.sender == block.coinbase, "the message sender must be the block producer");
        _;
    }

    modifier onlyZeroGasPrice() {
        require(tx.gasprice == 0, "gasprice is not zero");
        _;
    }

    modifier onlyNotInit() {
        require(!_alreadyInit, "the contract already init");
        _;
    }

    modifier onlyInit() {
        require(_alreadyInit, "the contract not init yet");
        _;
    }

    modifier onlySlash() {
        require(msg.sender == SLASH_ADDR, "the message sender must be slash contract");
        _;
    }

    modifier onlyHub() {
        require(msg.sender == HUB_ADDR, "the message sender must be hub contract");
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
}
