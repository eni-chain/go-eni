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

contract DelegateCallBase is LocalLog {
    //administrator address
    address internal _admin;
    address internal _impl;

    event UpdateAdmin(address indexed oldAdmin, address indexed newAdmin);

    event UpdateImpl(address indexed oldImpl, address indexed newImpl);

    modifier onlyAdmin() {
        require(msg.sender == _admin, "The message sender must be administrator");
        _;
    }

    modifier onlyValidContract(address impl){
        require(impl.code.length != 0, "implement contracts has no code");
        _;
    }

    function updateAdmin(address admin) external onlyAdmin {
        llog(DEBUG, abi.encodePacked("updateAdmin, old admin:", H(_admin), ", new admin:", H(admin)));
        emit UpdateAdmin(_admin, admin);

        return _setAdmin(admin);
    }

    function updateImpl(address impl) external onlyAdmin  {
        llog(DEBUG, abi.encodePacked("updateImpl, old impl:", H(_impl), ", new impl:", H(impl)));
        emit UpdateImpl(_impl, impl);

        return _setImpl(impl);
    }

    function _setAdmin(address admin) internal {
        _admin = admin;
    }

    function getAdmin() external view returns (address){
        return _admin;
    }

    //set delegate call's implementation contract address
    function _setImpl(address impl) internal onlyValidContract(impl){
        _impl = impl;
    }

    //get delegate call's implementation contract address
    function _getImpl() internal view returns (address) {
        return _impl;
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