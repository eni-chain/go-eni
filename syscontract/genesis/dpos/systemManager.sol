// SPDX-License-Identifier: GPL-3.0

pragma solidity >= 0.8.0;

import "./common.sol";
import "./localLog.sol";

contract SystemManager is LocalLog {
    //eni chain system address
    address internal _sys;

    event UpdateSysAddr(address indexed oldSysAddr, address indexed newSysAddr);

    modifier onlySystem() {
        require(msg.sender == _sys, "The message sender must be system address");
        _;
    }

    function updateSysAddr(address addr) external onlySystem {
        llog(DEBUG, abi.encodePacked("updateSysAddr, old system address:", H(_sys), ", new system address:", H(addr)));
        emit UpdateSysAddr(_sys, addr);

        return _setSysAddr(addr);
    }

    function _setSysAddr(address addr) internal {
        _sys = addr;
    }

    function getSysAddr() internal view returns (address){
        return _sys;
    }
}

