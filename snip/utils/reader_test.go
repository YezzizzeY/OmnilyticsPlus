package utils

import (
	"fmt"
	"math/big"
	"testing"
)

func scaleAndTruncate(nums []*big.Int) []*big.Int {
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(20), nil)
	var scaledNums []*big.Int

	for i, num := range nums {
		// Divide the number by 10^20
		scaledNum := new(big.Int).Div(num, divisor)
		scaledNums = append(scaledNums, scaledNum)
		if scaledNum.Cmp(big.NewInt(2)) == 0 {
			fmt.Println(i)
			fmt.Println(nums[i])
		}
	}
	return scaledNums
}

func TestName(t *testing.T) {
	_, numbers := ReadFile("../data/betafc1.txt", false)

	// Get the data from the fourth line
	data := numbers[0]
	fmt.Println(data[50:100])
	// Scale and truncate
	scaledData := scaleAndTruncate(data)

	fmt.Println(scaledData)

	fmt.Println("Total numbers:", len(scaledData))
	// Assuming you have the function Prime128Value() defined elsewhere

	sum := new(big.Int)
	for _, val := range scaledData {
		square := new(big.Int).Mul(val, val)
		sum.Add(sum, square)
	}

	fmt.Println("The sum of squares is:", sum)
}

func Prime128Value() *big.Int {
	p := big.NewInt(2)
	p.Exp(p, big.NewInt(127), nil)
	p.Sub(p, big.NewInt(1))
	return p
}
