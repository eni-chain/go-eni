// SPDX-License-Identifier: GPL-3.0


pragma solidity >= 0.8.0;

uint constant consensusSize = 21; 

uint constant ED25519_VERIFY_PRECOMPILED = 0x03ef;

address constant ADMIN_ADDR = 0x0000000000000000000000000000000000001000;
address constant HUB_ADDR = 0x0000000000000000000000000000000000001001;
address constant VALIDATOR_MANAGER_ADDR = 0x0000000000000000000000000000000000001002;
address constant VRF_ADDR = 0x0000000000000000000000000000000000001003;
address constant VOTER_MANAGER_ADDR = 0x0000000000000000000000000000000000001004;
address constant SLASH_ADDR = 0x0000000000000000000000000000000000001005;


contract Common {
    bool public alreadyInit = false;


    modifier onlyCoinbase() {
        require(msg.sender == block.coinbase, "the message sender must be the block producer");
        _;
    }

    modifier onlyZeroGasPrice() {
        require(tx.gasprice == 0, "gasprice is not zero");
        _;
    }

    modifier onlyNotInit() {
        require(!alreadyInit, "the contract already init");
        _;
    }

    modifier onlyInit() {
        require(alreadyInit, "the contract not init yet");
        _;
    }

    modifier onlySlash() {
        require(msg.sender == SLASH_ADDR, "the message sender must be slash contract");
        _;
    }

}

interface IValidatorManager {
    function getPubkey(address validator) external returns (bytes memory);
    function getValidatorSet() external  returns (address[] memory);
}


