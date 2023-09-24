package main

import (
	"fmt"
	"math/big"
	"snip/circuit"
	"snip/share"
	"testing"
)

func TestName(t *testing.T) {

	mod := share.Prime128Value()

	// this is comparing circuit and normal calculating [(x1^2 + x2^2 + ...) -1] * (...-2) ... (...-p)
	args := []*big.Int{big.NewInt(3), big.NewInt(2), big.NewInt(5)}
	p := 7

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

	a := circuit.FormCircuitByInput(args, p, mod)
	fmt.Println("circuit result is: ", a.Output)

	fs, gs, hs := FormFGH(a)

	fmt.Println("fs: ", fs)
	fmt.Println("gs: ", gs)
	fmt.Println("hs: ", hs)

	MulGates := circuit.GatesOfType(a, circuit.Gate_Mul)
	fmt.Println("hs(): ", hs.ValueAt(big.NewInt(int64(len(MulGates)))))

}
