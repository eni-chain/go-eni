// SPDX-License-Identifier: GPL-3.0

pragma solidity >= 0.8.0;

import "./common.sol";
import "./localLog.sol";

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
