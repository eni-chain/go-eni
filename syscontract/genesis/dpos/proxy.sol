// SPDX-License-Identifier: GPL-3.0

pragma solidity >= 0.8.0;

import "./common.sol";

contract ProxyContract is DelegateCallBase, LocalLog {
    //A structured storage slot parameter that defines logical contract address
    //value is hex("LogicContractAddressSlot")
    //uint256 private constant implSlotParam = 0x4c6f676963436f6e747261637441646472657373536c6f74;

    //system call, which can only be called once at deployment time
    function init(address impl) external {
        address addr = _getImpl();
        require(addr.code.length == 0, "init method can only be called once.");

        llog(DEBUG, abi.encodePacked("init implementation contract: ", H(impl)));
        _setImpl(impl);

        //bytes4 selector = bytes4(keccak256("init(address)"));
        //bytes memory data = abi.encodeWithSelector(selector, INIT_ADMIN_ADDR);
        //bytes memory data = abi.encodeCall(impl.init, (INIT_ADMIN_ADDR));

        bytes memory data = abi.encodeWithSignature("init(address)", INIT_ADMIN_ADDR);
        (bool success, bytes memory result) = impl.delegatecall(data);

        require(success, "DelegateCall failed");
        llog(DEBUG, abi.encodePacked("Delegate call implementation contract[", H(impl), "] succeed, return: ", H(result)));
    }

    function init(bytes memory bytecode) external {
        address addr = _getImpl();
        require(addr.code.length == 0, "Init method can only be called once.");

        assembly {
            addr := create(0,add(bytecode, 0x20), mload(bytecode))
        }
        require(addr.code.length != 0, "Create implementation contract failed fail");
        llog(DEBUG, abi.encodePacked("Deploy implementation contract[", H(addr), "] succeed."));

         _setImpl(addr);
        llog(DEBUG, abi.encodePacked("Set implementation contract[", H(addr), "] succeed."));

        bytes memory data = abi.encodeWithSignature("init(address)", INIT_ADMIN_ADDR);
        llog(DEBUG, abi.encodePacked("pack calldata[", H(data), "] succeed."));
        (bool success, bytes memory result) = addr.delegatecall(data);
        llog(DEBUG, abi.encodePacked("delegatecall finished, success:", S(success)));
        require(success, "DelegateCall failed");

        llog(DEBUG, abi.encodePacked("delegate call implementation contract[", H(addr), "] succeed, return: ", H(result)));
    }

    // //get implementation contract address
    // function getImpl() internal view returns (address impl) {
    //     bytes memory bs = new bytes(32);
    //     assembly {
    //         mstore(add(bs, 0x20), implSlotParam)
    //         let slot := sub(keccak256(add(bs, 0x20), 0x20), 1)
    //         impl := sload(slot)
    //     }
    // }

    // //set implementation contract address
    // function setImpl(address impl) internal  {
    //     require(impl.code.length > 0, "Invalid implementation address");

    //     bytes memory bs = new bytes(32);
    //     assembly {
    //         mstore(add(bs, 0x20), implSlotParam)
    //         let slot := sub(keccak256(add(bs, 0x20), 0x20), 1)
    //         sstore(slot, impl)
    //     }
    // }

    // captures all calls and forwards them to the logical contract
    fallback(bytes calldata data) external payable returns (bytes memory) {
        address impl = _getImpl();
        llog(DEBUG, abi.encodePacked(H(msg.sender), " delegate call address:", H(impl), ", calldata:", H(data)));
        require(impl.code.length > 0, "Invalid implementation address");

        (bool success, bytes memory result) = impl.delegatecall(data);
        //llog(DEBUG, abi.encodePacked("delegate call implementation contract[", H(impl), "] succeed:", S(success), ", result: ", H(result)));
        require(success, "DelegateCall failed");

        return result;
    }
}