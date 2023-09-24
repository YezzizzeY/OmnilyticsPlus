// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract DataStorage {

    // Define a struct to store the three key-value pairs for value
    struct Data {
        uint256[] grad;
        uint256[] beta;
        uint256[] sigma;
    }

    // Use a mapping to store the key and its corresponding Data, and set it to private
    mapping(uint256 => Data) private dataStore;

    // Manually create getter functions
    function getGrad(uint256 key) public view returns (uint256[] memory) {
        return dataStore[key].grad;
    }

    function getBeta(uint256 key) public view returns (uint256[] memory) {
        return dataStore[key].beta;
    }

    function getSigma(uint256 key) public view returns (uint256[] memory) {
        return dataStore[key].sigma;
    }

    // Store data
    function storeData(uint256 key, uint256[] memory grad, uint256[] memory beta, uint256[] memory sigma) public {
        // Check if all three arrays have the same length
        require(grad.length == beta.length && grad.length == sigma.length, "All arrays must have the same length");

        Data storage entry = dataStore[key];
        entry.grad = grad;
        entry.beta = beta;
        entry.sigma = sigma;
    }

    // For different keys, add the arrays of their corresponding key-value pairs at their respective positions
    function addArrays(uint256 key1, uint256 key2) public {
        Data storage data1 = dataStore[key1];
        Data storage data2 = dataStore[key2];

        require(data1.grad.length == data2.grad.length, "Grad arrays length mismatch");
        require(data1.beta.length == data2.beta.length, "Beta arrays length mismatch");
        require(data1.sigma.length == data2.sigma.length, "Sigma arrays length mismatch");

        for (uint256 i = 0; i < data1.grad.length; i++) {
            data1.grad[i] = data1.grad[i] + data2.grad[i];
            data1.beta[i] = data1.beta[i] + data2.beta[i];
            data1.sigma[i] = data1.sigma[i] + data2.sigma[i];
        }
    }
}
