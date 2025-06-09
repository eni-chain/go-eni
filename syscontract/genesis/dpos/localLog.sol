// SPDX-License-Identifier: GPL-3.0

pragma solidity >= 0.8.0;

import "./common.sol";

contract LocalLog {
    //log level
    uint constant DEBUG = 1;
    uint constant INFO = 2;
    uint constant WARN = 3;
    uint constant ERROR = 4;

    function char2Hex(uint8 c) internal pure returns (bytes1) {
        if (c < 10) {
            return bytes1(uint8(c + 0x30)); // '0'-'9'
        } else {
            return bytes1(uint8(c - 10 + 0x61)); // 'a'-'f'
        }
    }

    //convert address to bytes
    function addr2Bytes(address a) internal pure returns (bytes memory b) {
        return abi.encodePacked(a);
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
  
        bool success;
        (success,) = address(uint160(LOCAL_NODE_LOG_PRECOMPILED)).staticcall(input);
        if(!success){
            revert("the call to the nodeLog precompiled contract failed");
        }

    }
}
