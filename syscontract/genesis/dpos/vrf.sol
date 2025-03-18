// SPDX-License-Identifier: GPL-3.0

pragma solidity >= 0.8.0;

import "./common.sol";

uint256 constant PUBKEY_LEN = 32; //ed25519 public key length
uint256 constant PRIKEY_LEN = 64; //ed25519 private key length
uint256 constant SEED_LEN = 64;   //random seed length
uint256 constant SIGN_LEN = 64;   //ed25519 signature length
uint256 constant HASH_LEN = 64;  //hash length

contract Vrf is nodeLog {
    //todo: add event and emit for every method

    //init rand seed, will be init by administrator
    bytes _initSeed;

    //administrator address
    address public _admin;

    //epoch => random seed
    mapping(uint256 => bytes)internal _seeds;

    //validator node address => ed25519 public key
    mapping (address => bytes) private _pubKeys;

    //map(epoch=>map(validator node address=>random))
    mapping (uint256=>mapping (address=>bytes)) private  _randoms;

    //invalid validator[], record the node not send random
    address[] private _unSendRandNodes;

    //valid validator[]
    address[] private _validNodes;

    modifier onlyAdmin() {
        require(msg.sender == _admin, "The message sender must be administrator");
        _;
    }

   modifier needInited() {
        require(_initSeed.length != 0, "The initial seed has not been initialized and dpos has not been started");
        _;
    }

    function initAdmin() external {
        require(msg.sender == ADMIN_ADDR, "The message sender must be administrator");
        _admin = msg.sender;
    }

    function updateAdmin(address admin) external onlyAdmin {
        //require(msg.sender == _admin, "Msg sender is not administrator");
        _admin = admin;
    }

    function init(bytes calldata rnd, uint256 epoch) external onlyAdmin {
        //require(msg.sender == _admin, "Msg sender is not administrator");
        require(_initSeed.length == 0, "vrf has been init!");

        _initSeed = rnd;
        _seeds[epoch] = rnd;
    }

    function getRandomSeed(uint256 epoch) external view needInited returns (bytes memory) {
        //require(_initSeed.length != 0, "vrf has not been init!");
        require(epoch > 1, "epoch number too small");

        //each random values is generated from the seeds of the previous epoch
        return _seeds[epoch-1];
    }

    // function setPubKey(address validator, bytes calldata pubkey) public {
    //     _pubKeys[validator] = pubkey;
    // }

    function verifyEd25519Sign(bytes memory pubKey, bytes memory signature, bytes memory msgHash) public view returns (bool) {
        require(pubKey.length == PUBKEY_LEN, "The public key length is not ed25519 public key size");
        require(signature.length == SIGN_LEN, "The signature length is not ed25519 signature size");
        require(msgHash.length == HASH_LEN, "The msg hash length is not SHA-512 hash size");

        // assemble input:
        // | PubKey   | Signature  |  msgHash   |
        // | 32 bytes | 64 bytes   |  64 bytes  |
        bytes memory input = bytes.concat(pubKey, signature, msgHash);
        //bytes32[2] memory output;
        bytes memory output = new bytes(32);

        assembly {
            let len := mload(input)
            if iszero(staticcall(not(0), ED25519_VERIFY_PRECOMPILED, add(input, 0x20), len, add(output, 0x20), 0x20)) {
                revert(0, 0)
            }
        }

        if(output[31] == 0){
            return  false;
        }

        printLog(ERROR, bytes.concat(bytes("verify rand succeed, user:"), toLog(msg.sender), bytes(", pubKey:"), toLog(pubKey), bytes(", signature:"), toLog(signature), bytes(", msgHash: "), toLog(msgHash)));
        return true;
    }

    function sendRandom(bytes calldata rnd, uint256 epoch) external needInited returns (bool success){
        //require(_initSeed.length != 0, "Needs init first!");
        require(epoch > 1, "Epoch number too small");
        require(_seeds[epoch-1].length == SEED_LEN, "Random values sent ahead of epoch!");
        require(rnd.length == SIGN_LEN, "Random length is not ed25519 signature size!");

        address nodeAddr;
        bytes memory pubKey;
        (nodeAddr, pubKey) = IValidatorManager(VALIDATOR_MANAGER_ADDR).getNodeAddrAndPubKey(msg.sender);
        require(pubKey.length != 0, "Msg sender is not validator operator");

        bool success = verifyEd25519Sign(pubKey, rnd, _seeds[epoch-1]);
        require(success == true, "Random is not signature that signed by validator");

        _randoms[epoch][nodeAddr] = rnd;
    }

    function updateConsensusSet(uint256 epoch) external needInited returns (address[] memory) {
        //todo: add permission- only the epoch module address can call the updateConsensusSet method
        //require(_randoms[epoch].length > 0, "Epoch has no random value!");
        require(keccak256(_seeds[epoch]) != keccak256(_initSeed), "Consensus set should be elected in next epoch");

        address[] memory validators = IValidatorManager(VALIDATOR_MANAGER_ADDR).getValidatorSet();
        require(validators.length > 0, "Validator set is empty");
        printLog(ERROR, bytes.concat(bytes("getValidatorSet:"), toLog(validators[0]), bytes(", "), toLog(validators[1])));

        //address[] memory validators = new address[](address(uint160(_randoms[epoch][keccak256("Vrf")])));
        for (uint i = 0; i < validators.length; ++i) {
            if(_randoms[epoch][validators[i]].length == 0){
                _unSendRandNodes.push(validators[i]);
                printLog(ERROR, bytes.concat(bytes("record unsend rand node:"), toLog(validators[i])));
            }else{
                _validNodes.push(validators[i]);
                printLog(ERROR, bytes.concat(bytes("record valid node:"), toLog(validators[i])));
            }
        }

        //todo: call slash contract to penalty evil node
        //ISlash(SLASH_ADDR).penaltyUnsendRandomValidator(_unSendRandNodes);

        address[] memory sorted = sortAddrs(_validNodes, epoch);
        printLog(ERROR, bytes.concat(bytes("sorted address:"), toLog(sorted[0]), bytes(", "),toLog(sorted[1])));
        address[] memory topN = getTopNAddresses(sorted, consensusSize);
        printLog(ERROR, bytes.concat(bytes("topN address:"), toLog(sorted[0]), bytes(", "),toLog(sorted[1])));

        //The seed of this epoch are generated for the next epoch to generate random values
        _seeds[epoch] = _seeds[epoch-1];
        for(uint i = 0; i < topN.length; ++i){
            _seeds[epoch] = addBytes(_seeds[epoch],  _randoms[epoch][topN[i]]);
        }
        printLog(ERROR, bytes.concat(bytes("epoch: "), toLog(epoch), bytes(", seed: "), toLog(_seeds[epoch])));

        //Empty the invalid node set and the valid node set for the next epoch
        delete _unSendRandNodes;
        delete _validNodes;

        return topN;
    }

    function compare(bytes memory a, bytes memory b) internal pure returns (bool) {
        return keccak256(a) < keccak256(b);
    }

    // Sort addresses in ascending order by random value
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

    // Adds two 64-byte bytes byte by byte
    function addBytes(bytes memory a, bytes memory b) internal pure returns (bytes memory) {
        require(a.length == SEED_LEN && b.length == SEED_LEN, "Invalid input length");

        bytes memory result = new bytes(SEED_LEN);
        for(uint256 i = 0; i < SEED_LEN; i++){
            result[i] = bytes1(uint8(a[i])+uint8(b[i]));
        }

        return result;
    }
}
