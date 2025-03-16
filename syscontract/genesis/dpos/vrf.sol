// SPDX-License-Identifier: GPL-3.0

pragma solidity >= 0.8.0;

import "./common.sol";

uint256 constant PUBKEY_LEN = 32;
uint256 constant PRIKEY_LEN = 64;
uint256 constant SEED_LEN = 64;
uint256 constant SIGN_LEN = 64;
uint256 constant HASH_LEN = 64;

contract Vrf {
    //todo: add event and emit

    address public _admin;

    bytes internal _initSeed;

    mapping(uint256 => bytes)internal _seeds;

    mapping (address => bytes) private _pubKeys;

    mapping (uint256=>mapping (address=>bytes)) private  _randoms;

    address[] private _unSendRandNodes;

    address[] private _validNodes;

    modifier onlyAdmin() {
        require(msg.sender == _admin, "the message sender must be administrator");
        _;
    }

    constructor(){
        _admin = ADMIN_ADDR;
    }

    function updateAdmin(address admin) external onlyAdmin {
        _admin = admin;
    }

    function init(bytes calldata rnd) external onlyAdmin {
        require(_initSeed.length == 0, "vrf has not been init!");

        _initSeed = rnd;
    }

    function getRandomSeed(uint256 epoch) external view returns (bytes memory) {
         require(epoch > 1, "epoch number too small");

        return _seeds[epoch-1];
    }


    function verifyEd25519Sign(bytes memory pubKey, bytes memory signature, bytes memory msgHash) internal view returns (bool) {
        require(pubKey.length == PUBKEY_LEN, "The public key length is not ed25519 public key size");
        require(pubKey.length == PRIKEY_LEN, "The private key length is not ed25519 public key size");
        require(msgHash.length == HASH_LEN, "The msg hash length is not SHA-512 hash size");

        // assemble input:
        // | PubKey   | Signature  |  msgHash   |
        // | 32 bytes | 64 bytes   |  64 bytes  |
        bytes memory input = bytes.concat(pubKey, signature, msgHash);
        //bytes32[2] memory output;
        bytes memory output = new bytes(1);

        assembly {
            let len := mload(input)
            if iszero(staticcall(not(0), ED25519_VERIFY_PRECOMPILED, add(input, 0x20), len, add(output, 0x20), 0x01)) {
                revert(0, 0)
            }
        }

        if(output[0] == 0){
            return  false;
        }

        return true;
    }

    function sendRandom(bytes calldata rnd, uint256 epoch) external returns (bool success){
        require(_initSeed.length > 0, "Send random needs init first!");

        require(epoch > 1, "Epoch number too small");
        require(_seeds[epoch-1].length == SEED_LEN, "Random values sent ahead of epoch!");
        require(rnd.length == SIGN_LEN, "Random length is not ed25519 signature size!");

        if(_pubKeys[msg.sender].length == 0){
            bytes memory pubkey = IValidatorManager(VALIDATOR_MANAGER_ADDR).getPubKey(msg.sender);
            if(pubkey.length == 0){
                revert("Msg sender is not validator");
            }

            _pubKeys[msg.sender] = pubkey;
        }

        bool success = verifyEd25519Sign(_pubKeys[msg.sender], rnd, _seeds[epoch-1]);
        require(success == true, "Random is not signature that signed by validator");

        _randoms[epoch][msg.sender] = rnd;
    }

    function updateConsensusSet(uint256 epoch) external returns (address[] memory) {
        //require(_randoms[epoch].length > 0, "Epoch has no random value!");

        address[] memory validators = IValidatorManager(VALIDATOR_MANAGER_ADDR).getValidatorSet();
        if(validators.length == 0){
            revert(" validator set is empty");
        }

        //address[] memory validators = new address[](address(uint160(_randoms[epoch][keccak256("Vrf")])));

        for (uint i = 0; i < validators.length; ++i) {
            if(_randoms[epoch][validators[i]].length == 0){
                _unSendRandNodes.push(validators[i]);
            }
        }

        //todo: call slash contract

        for (uint i = 0; i < validators.length; ++i) {
            bool found = false;

            for(uint ii = 0; ii < _unSendRandNodes.length; ++ii){
                if(_unSendRandNodes[ii] == validators[i]){
                    found = true;
                    break;
                }
            }

            if(!found){
                _validNodes.push(validators[i]);
            }
        }

        address[] memory sorted = sortAddrs(_validNodes, epoch);
        address[] memory topN = getTopNAddresses(sorted, consensusSize);

        _seeds[epoch] = _seeds[epoch-1];
        for(uint i = 0; i < topN.length; ++i){
            _seeds[epoch] = addBytes(_seeds[epoch],  _randoms[epoch][topN[i]]);
        }

        delete _unSendRandNodes;
        delete _validNodes;

        return topN;

    }

    function compare(bytes memory a, bytes memory b) internal pure returns (bool) {
        return keccak256(a) < keccak256(b);
    }

    function sortAddrs(address[] memory array, uint256 epoch) internal view returns (address[] memory) {
        uint256 n = array.length;
        if (n <= 1){
            return array;
        }

        //bytes[] memory sorted = array.clone();
        for (uint256 i = 0; i < n - 1; i++) {
            for (uint256 j = 0; j < n - 1 - i; j++) {
                bytes memory rndFront = _randoms[epoch][array[j]];
                bytes memory rndBack = _randoms[epoch][array[j+1]];
                if (compare(rndBack, rndFront)) {
                    address temp = array[j];
                    array[j] = array[j + 1];
                    array[j + 1] = temp;
                }
            }
        }

        return array;
    }

    function getTopNAddresses(address[] memory array, uint256 n) internal pure returns (address[] memory) {
        if(array.length < n){
            n = array.length;
        }

        require(n >= 1, "Invalid slice length");

        address[] memory result = new address[](n);
        for (uint256 i = 0; i < n; i++) {
            result[i] = array[i];
        }
        return result;
    }

    function addBytes(bytes memory a, bytes memory b) internal pure returns (bytes memory) {
        require(a.length == 64 && b.length == 64, "Invalid input length");

        uint256 numA0;
        uint256 numA1;
        uint256 numB0;
        uint256 numB1;

        assembly {
            numA0 := mload(add(a, 0x20))
            numA1 := mload(add(a, 0x40))
            numB0 := mload(add(b, 0x20))
            numB1 := mload(add(b, 0x40))
        }

        uint256 carry = numA0 + numB0;
        uint256 high = numA1 + numB1;

        bytes memory result = new bytes(64);
        assembly {
            mstore(add(result, 0x20), carry)
            mstore(add(result, 0x40), high)
        }

        return result;
    }
}