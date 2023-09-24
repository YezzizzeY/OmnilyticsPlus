// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract CommitmentContract {
    struct Element {
        uint256 data1;
        uint256 data2;
        uint256 data3;
    }

    struct Commitment {
        Element[5] elements;
    }

    Commitment[] commitments;

    // Func1: 随机生成4个commitment
    function generateRandomCommitments() public {
        for (uint i = 0; i < 4; i++) {
            Commitment memory newCommitment;
            for (uint j = 0; j < 5; j++) {
                newCommitment.elements[j] = Element(1, 2, 3);
            }
            pushCommitment(newCommitment);
        }
    }

    // Helper function to handle memory to storage assignment
    function pushCommitment(Commitment memory newCommitment) internal {
        commitments.push();
        for (uint i = 0; i < 5; i++) {
            commitments[commitments.length - 1].elements[i] = newCommitment.elements[i];
        }
    }

    // Func2: 将这4个commitment相乘
    function multiplyCommitments() public {
        require(commitments.length >= 4, "At least 4 commitments are required");

        Commitment memory multipliedCommitment;

        for (uint j = 0; j < 5; j++) {
            multipliedCommitment.elements[j].data1 = 1;
            multipliedCommitment.elements[j].data2 = 1;
            multipliedCommitment.elements[j].data3 = 1;

            for (uint i = 0; i < 4; i++) {
                multipliedCommitment.elements[j].data1 *= commitments[i].elements[j].data1;
                multipliedCommitment.elements[j].data2 *= commitments[i].elements[j].data2;
                multipliedCommitment.elements[j].data3 *= commitments[i].elements[j].data3;
            }
        }

        pushCommitment(multipliedCommitment);
    }
}
