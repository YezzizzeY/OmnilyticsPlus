package circuit

import (
	"fmt"
	"math/big"
	"snip/share"
	"testing"
)

func TestCircuitValue(t *testing.T) {

	// this is comparing circuit and normal calculating [(x1^2 + x2^2 + ...) -1] * (...-2) ... (...-p)
	mod := share.Prime128Value()

	args := []*big.Int{big.NewInt(5125), new(big.Int).Sub(mod, big.NewInt(41))}
	p := 3

	n := len(args)

	MulRe := big.NewInt(0)
	for i := 0; i < n; i++ {
		a := args[i]
		t1 := new(big.Int).Mul(a, a)
		t1 = t1.Mod(t1, mod)
		MulRe = MulRe.Add(MulRe, t1)
		MulRe = MulRe.Mod(MulRe, mod)
	}

	Re := big.NewInt(1)
	for i := 1; i <= p; i++ {
		t1 := new(big.Int).Sub(MulRe, big.NewInt(int64(i)))
		t1 = t1.Mod(t1, mod)
		Re = Re.Mul(Re, t1)
		Re = Re.Mod(Re, mod)
	}

	fmt.Println("total result is: ", Re)

	a := FormCircuitByInput(args, p, mod)

	fmt.Println("Prime128 value = :", mod)
	fmt.Println(a)

	mulGates := GatesOfType(a, Gate_Mul)
	fmt.Println(mulGates)

}
