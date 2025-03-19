// SPDX-License-Identifier: GPL-3.0


pragma solidity >= 0.8.0;

uint constant consensusSize = 2;

uint constant MIN_PLEDGE_AMOUNT = 10000;

uint constant ED25519_VERIFY_PRECOMPILED = 0xa1;
uint constant LOCAL_NODE_LOG_PRECOMPILED = 0xa2;

address constant ADMIN_ADDR = 0x251604eBfD1ddeef1F4f40b8F9Fc425538BE1339;
address constant HUB_ADDR = 0x0000000000000000000000000000000000001001;
address constant VALIDATOR_MANAGER_ADDR = 0x0000000000000000000000000000000000001002;
address constant VRF_ADDR = 0x0000000000000000000000000000000000001003;
address constant VOTER_MANAGER_ADDR = 0x0000000000000000000000000000000000001004;
address constant SLASH_ADDR = 0x0000000000000000000000000000000000001005;

uint constant DEBUG = 1;
uint constant INFO = 2;
uint constant WARN = 3;
uint constant ERROR = 4;

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


contract LocalLog {

    function char2Hex(uint8 c) internal pure returns (bytes1) {
        if (c < 10) {
            return bytes1(uint8(c + 0x30)); // '0'-'9'
        } else {
            return bytes1(uint8(c - 10 + 0x61)); // 'a'-'f'
        }
    }

    //convert address to bytes
    function addr2Bytes(address a) internal pure returns (bytes memory b) {
        assembly {
            let m := mload(0x40)
            a := and(a, 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF)
            mstore(
                add(m, 20),
                xor(0x140000000000000000000000000000000000000000, a)
            )
            mstore(0x40, add(m, 52))
            b := m
        }
    }

    // convert bytes to hex string
    function H(bytes memory bs) public pure returns (bytes memory) {
        bytes memory hexStr = new bytes(bs.length * 2);
        for (uint i = 0; i < bs.length; i++) {
            uint8 v = uint8(bs[i]);
            hexStr[i*2] = char2Hex(v >> 4);
            hexStr[i*2+1] = char2Hex(v & 0xf);
        }
        return hexStr;
    }

    // convert address to hex string
    function H(address addr) public pure returns (bytes memory) {
        bytes memory addrBytes = addr2Bytes(addr);
        return H(addrBytes);
    }

    //convert uint256 to string
   function S(uint256 value) public pure returns (bytes memory) {
        if (value == 0) {
            return "0";
        }

        uint256 temp = value;
        uint256 digits;
        while (temp != 0) {
            digits++;
            temp /= 10;
        }

        bytes memory buffer = new bytes(digits);
        while (value != 0) {
            digits -= 1;
            buffer[digits] = bytes1(uint8(48 + uint256(value % 10)));
            value /= 10;
        }

        return buffer;
    }

    //convert bool to string
    function S(bool value) public pure returns (bytes memory) {
        if(value) {
            return bytes("true");
        } else {
            return bytes("false");
        }
    }

    function llog(uint level, bytes memory logs) public view {
        bytes memory input = bytes.concat(S(level), logs);

        bytes memory output = new bytes(32);
        uint256 len = input.length;
        assembly{
            if iszero(staticcall(not(0), LOCAL_NODE_LOG_PRECOMPILED, add(input, 0x20), len, add(output, 0x20), 0x20)){
                revert(0, 0)
            }
        }
    }
}