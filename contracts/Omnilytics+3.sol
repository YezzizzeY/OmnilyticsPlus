// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract DataContract {
    
    // 定义数据结构Data
    struct Data {
        uint256[5] rshare;
        uint256[5] outputwire;
    }

    // 用于存储Data的映射
    mapping(uint256 => Data) dataMapping;

    // 当前的key
    uint256 public currentKey = 1;

    // 函数：存储10个数字，前5个存到outputwire，后5个存到rshare
    function storeData(uint256[10] memory numbers) public {
        require(msg.sender == address(this), "Only the contract itself can add data.");

        Data memory newData;

        // 分割并存储数据
        for (uint i = 0; i < 5; i++) {
            newData.outputwire[i] = numbers[i];
            newData.rshare[i] = numbers[i + 5];
        }

        // 存储到mapping中并更新key
        dataMapping[currentKey] = newData;
        currentKey++;
    }
}
