package utils

import (
	"fmt"
	"math/big"
	"testing"
)

func TestPolyRecover(t *testing.T) {
	p := big.NewInt(23)
	n := 3
	x := []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}
	y := []*big.Int{big.NewInt(7), big.NewInt(16), big.NewInt(6)}
	r := PolyRecover(n, p, x, y)
	fmt.Print("recovered polynomial: ", r)
}
