package mpc

import (
	"fmt"
	"math/big"
	"snip/share"
	"testing"
)

func TestGenerateBeaver(t *testing.T) {

	a5 := big.NewInt(4)
	b6 := big.NewInt(3)

	beavers := GenerateBeaver(a5, b6, share.Prime128Value(), 1, 3)

	mul := []share.Share{}
	for i := 0; i < 3; i++ {
		mul = append(mul, CalculateMul(beavers[i]))
	}

	v := share.Recovershares(mul)

	fmt.Println("recovered share value: ", v)
}
