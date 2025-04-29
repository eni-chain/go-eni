// SPDX-License-Identifier: GPL-3.0

pragma solidity >= 0.8.0;

import "./common.sol";

contract ProxyContract is DelegateCallBase {
    //A structured storage slot parameter that defines logical contract address
    //value is hex("LogicContractAddressSlot")
    //uint256 private constant implSlotParam = 0x4c6f676963436f6e747261637441646472657373536c6f74;

    // //get delegate call's implementation contract address
    // function _getImpl() internal view returns (address impl) {
    //     bytes memory key = new bytes(32);
    //     assembly {
    //         //let offset := add(key, 0x20)
    //         mstore(add(key, 0x20), IMPL_SLOT_BASE)

    //         let slot := sub(keccak256(add(key, 0x20), 0x20), 1)
    //         let addr := sload(slot)
    //         impl := and(addr, 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF)
    //     }

    //     return impl;
    // }

    // //set delegate call's implementation contract address
    // function _setImpl(address impl) internal  {
    //     require(impl.code.length > 0, "Invalid implementation address");

    //     bytes memory key = new bytes(32);
    //     bytes memory addr = new bytes(32);
    //     assembly {
    //         //let offset := add(bs, 0x20)
    //         mstore(add(key, 0x20), IMPL_SLOT_BASE)
    //         mstore(add(addr, 0x20), impl)

    //         let slot := sub(keccak256(add(key, 0x20), 0x20), 1)
    //         sstore(slot, addr)
    //     }
    // }

    //system call, which can only be called once at deployment time
    // function init(address impl) external {
    //     address addr = _getImpl();
    //     require(addr.code.length == 0, "init method can only be called once.");

    //     llog(DEBUG, abi.encodePacked("init implementation contract: ", H(impl)));
    //     _setImpl(impl);

    //     bytes4 selector = bytes4(keccak256("init(address)"));
    //     bytes memory data = abi.encodeWithSelector(selector, INIT_ADMIN_ADDR);
    //     bytes memory data = abi.encodeCall(impl.init, (INIT_ADMIN_ADDR));

    //     bytes memory data = abi.encodeWithSignature("init(address)", INIT_ADMIN_ADDR);
    //     (bool success, bytes memory result) = impl.delegatecall(data);

    //     require(success, "DelegateCall failed");
    //     llog(DEBUG, abi.encodePacked("Delegate call implementation contract[", H(impl), "] succeed, return: ", H(result)));
    // }

    event Init(address indexed self, address indexed admin, address indexed impl);

    function init(bytes memory bytecode) external {
        address addr = _getImpl();
        require(addr.code.length == 0, "Init method can only be called once.");

        assembly {
            addr := create(0,add(bytecode, 0x20), mload(bytecode))
        }
        require(addr.code.length != 0, "Create implementation contract failed fail");
        llog(DEBUG, abi.encodePacked("init, deploy implementation contract: ", H(addr)));

         _setImpl(addr);
        _setAdmin(INIT_ADMIN_ADDR);
        llog(DEBUG, abi.encodePacked("init, set impl:", H(addr), ", set admin:", H(_admin)));

        emit Init(address(this), _admin, addr);
    }

    // captures all calls and forwards them to the logical contract
    fallback(bytes calldata data) external payable returns (bytes memory) {
        address impl = _getImpl();
        require(impl.code.length > 0, "Invalid implementation address");

        (bool success, bytes memory result) = impl.delegatecall(data);
        require(success, "DelegateCall failed");
        //llog(DEBUG, abi.encodePacked("fallback, delegate call impl:", H(impl), ", succeed:", S(success), ", result: ", H(result)));

        return result;
    }
}