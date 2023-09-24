pragma solidity ^0.8.0;

// this contract is for calculating gas consumption in generating g^alpha^m 
contract CyclicGroupSampler {

    // Using q = 2^256 - 2^32 - 977, a prime used in secp256k1 curve
    uint256 public constant q = 340282366920938463463374607431768211297;
    
    // This is a simplistic choice for demonstration and might not be a generator for our prime field.
    uint256 public constant g = 3;
    
    uint256[] public values;

    constructor() {}

    function sampleAndPublish(uint256 m) public {
        uint256 alpha = 5;

        for (uint256 i = 0; i < m; i++) {
            uint256 exponent = modExp(alpha, i, q); // Calculate alpha^i mod q
            uint256 result = modExp(g, exponent, q);  // Calculate g^(alpha^i) mod q
            values.push(result);
        }
    }


    // Modular exponentiation
    function modExp(uint256 base, uint256 exponent, uint256 modulus) internal pure returns (uint256 result) {
        base = base % modulus;
        exponent = exponent % (modulus - 1);  // 指数取模 modulus-1
        
        result = 1;
        while (exponent != 0) {
            if (exponent % 2 == 1) {  // 如果指数的最低位是1
                result = (result * base) % modulus;
            }
            exponent = exponent >> 1;  // 右移指数
            base = (base * base) % modulus;
        }
        return result;
    }

}
