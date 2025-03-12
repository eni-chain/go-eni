// SPDX-License-Identifier: GPL-3.0

pragma solidity >= 0.8.0;

import "./common.sol";

uint256 constant PUBKEY_LEN = 32; //ed25519公钥长度
uint256 constant PRIKEY_LEN = 64; //ed25519私钥长度
uint256 constant SEED_LEN = 64;   //随机值种子长度
uint256 constant SIGN_LEN = 64;   //ed25519签名长度
uint256 constant HASH_LEN = 64;  //ed25519签名前，需要先计算被签名数据的SHA-512哈希值，然后对哈希值签名，该哈希值长度为64字节

contract Vrf {
    //todo: 为每个方法添加event

    //管理员地址
    address public _admin;

    //初始化随机值种子,为外部（管理员）调用init方法传参进来的btc区块哈希
    bytes internal _initSeed;

    //epoch => 随机值种子
    mapping(uint256 => bytes)internal _seeds;

    //验证者节点地址的公钥，用于校验签名生成的随机值
    mapping (address => bytes) private _pubKeys;

    //map(epoch=>map(validator=>random)), 记录验证者节点在每个epoch提交的随机值
    mapping (uint256=>mapping (address=>bytes)) private  _randoms;

    //validator[],无效验证者，记录当前轮未发送随机值的验证者
    address[] private _unSendRandNodes;

    //validator[],有效验证者，记录本轮发送了随机值的验证者
    address[] private validNodes;

    modifier onlyAdmin() {
        require(msg.sender == _admin, "the message sender must be administrator");
        _;
    }

    constructor(){
        //todo: 对admin地址确认
        _admin = ADMIN_ADDR;
    }

    function updateAdmin(address admin) external onlyAdmin {
        //require(msg.sender == _admin, "Msg sender is not administrator");
        _admin = admin;
    }

    function init(bytes calldata rnd) external onlyAdmin {
        //require(msg.sender == _admin, "Msg sender is not administrator");
        require(_initSeed.length == 0, "vrf has not been init!");

        _initSeed = rnd;
    }

    function getRandomSeed(uint256 epoch) external view returns (bytes memory) {
         require(epoch > 1, "epoch number too small");

        //每轮随机值都是根据上一轮的随机值种子生成的
        return _seeds[epoch-1];
    }

    // function setPubKey(address validator, bytes calldata pubkey) public {
    //     //todo: 调用验证者管理合约检查validator是否为验证者，并返回验证者公钥
    //     _pubKeys[validator] = pubkey;
    // }

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
            //todo: 确认validatorManager合约实现后与IValidatorManager接口匹配
            bytes memory pubkey = IValidatorManager(VALIDATOR_MANAGER_ADDR).getPubkey(msg.sender);
            if(pubkey.length == 0){
                revert("Msg sender is not validator");
            }

            _pubKeys[msg.sender] = pubkey;
        }

        //校验随机值是否由验证者签名而来, 随机值种子和SHA-512哈希的长度都是64，因此不再对种子计算哈希
        bool success = verifyEd25519Sign(_pubKeys[msg.sender], rnd, _seeds[epoch-1]);
        require(success == true, "Random is not signature that signed by validator");

        _randoms[epoch][msg.sender] = rnd;
    }

    function updateConsensusSet(uint256 epoch) external returns (address[] memory) {
        //以epoch为索引，取出当前轮所有的验证者与其随机值
        //require(_randoms[epoch].length > 0, "Epoch has no random value!");

        address[] memory validators = IValidatorManager(VALIDATOR_MANAGER_ADDR).getValidatorSet();
        if(validators.length == 0){
            revert(" validator set is empty");
        }

        //address[] memory validators = new address[](address(uint160(_randoms[epoch][keccak256("Vrf")])));
        
        for (uint i = 0; i < validators.length; ++i) {
            //遍历validator集合，如果有validator在随机值记录表中未查到，则将其记录到未发送随机值列表
            if(_randoms[epoch][validators[i]].length == 0){
                _unSendRandNodes.push(validators[i]);
            }
        }

        //todo: 第二期开发，以未发送随机值列表未参数调用slash合约

        //将验证者集合中未发送随机值的节点剔除
        for (uint i = 0; i < validators.length; ++i) {
            bool found = false;

            for(uint ii = 0; ii < _unSendRandNodes.length; ++ii){
                if(_unSendRandNodes[ii] == validators[i]){
                    found = true;
                    break;
                }
            }
            
            //如果未找到，则将当前validator插入remain集合中
            if(!found){
                validNodes.push(validators[i]);
            }
        }
        
        address[] memory sorted = sortAddrs(validNodes);
        address[] memory topN = getTopNAddresses(sorted, consensusSize);

        //生成本轮的随机值种子，供下一轮生成随机值使用
        _seeds[epoch] = _seeds[epoch-1];
        for(uint i = 0; i < topN.length; ++i){
            _seeds[epoch] = addBytes(_seeds[epoch],  _randoms[epoch][topN[i]]);
        }

        //清空无效节点集合和有效节点集合，供下一轮使用
        delete _unSendRandNodes;
        delete validNodes;

        return topN;

    }

    // 比较两个 bytes 元素（通过哈希值比较）
    // function compare(bytes memory a, bytes memory b) internal pure returns (bool) {
    //     return keccak256(a) < keccak256(b);
    // }

    // 对 bytes[] 进行升序排序
    function sortAddrs(address[] memory array) internal pure returns (address[] memory) {
        uint256 n = array.length;
        if (n <= 1){
            return array;
        }

        //bytes[] memory sorted = array.clone();
        for (uint256 i = 0; i < n - 1; i++) {
            for (uint256 j = 0; j < n - 1 - i; j++) {
                //if (compare(array[j + 1], array[j])) {
                if (array[j + 1] <= array[j]) {
                    address temp = array[j];
                    array[j] = array[j + 1];
                    array[j + 1] = temp;
                }
            }
        }

        return array;
    }

    function getTopNAddresses(address[] memory array, uint256 n) internal pure returns (address[] memory) {
        require(n <= array.length, "Invalid slice length");
        
        address[] memory result = new address[](n); // 创建新数组（长度为 n）
        for (uint256 i = 0; i < n; i++) {
            result[i] = array[i]; // 逐个复制元素
        }
        return result;
    }

    // 将两个 64 字节的 bytes 数组逐字节相加
    function addBytes(bytes memory a, bytes memory b) internal pure returns (bytes memory) {
        require(a.length == 64 && b.length == 64, "Invalid input length");

        // 将 bytes 转换为 uint256
        uint256 numA0;
        uint256 numA1;
        uint256 numB0;
        uint256 numB1;

        assembly {
            // 读取第一个 bytes 数组的前 256 位
            numA0 := mload(add(a, 0x20)) // mload 读取 32 字节（256 位）
            // 读取第一个 bytes 数组的后 256 位
            numA1 := mload(add(a, 0x40))
            // 读取第二个 bytes 数组的前 256 位
            numB0 := mload(add(b, 0x20))
            // 读取第二个 bytes 数组的后 256 位
            numB1 := mload(add(b, 0x40))
        }

        // 逐位相加并处理进位，不考虑溢出问题，因为溢出也无更多的字节位存储
        uint256 carry = numA0 + numB0; // 低 256 位相加
        uint256 high = numA1 + numB1; // 高 256 位相加

        // 合并结果为 bytes 数组
        bytes memory result = new bytes(64);
        assembly {
            // 写入低 256 位
            mstore(add(result, 0x20), carry)
            // 写入高 256 位
            mstore(add(result, 0x40), high)
        }

        return result;
    }

}
