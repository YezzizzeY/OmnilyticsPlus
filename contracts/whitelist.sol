// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Whitelist {
    mapping(uint256 => bool) public whitelist;

    function addToWhitelist(uint256 dataOwnerID) public {
        whitelist[dataOwnerID] = true;
    }

    function addMultipleToWhitelist(uint256[5] memory dataOwnerIDs) public {
        for (uint i = 0; i < 5; i++) {
            whitelist[dataOwnerIDs[i]] = true;
        }
    }

    function removeFromWhitelist(uint256 dataOwnerID) public {
        whitelist[dataOwnerID] = false;
    }

    function isWhitelisted(uint256 dataOwnerID) public view returns(bool) {
        return whitelist[dataOwnerID];
    }
}
