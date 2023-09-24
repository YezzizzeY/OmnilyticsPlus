package main

import (
	"fmt"
	"math/big"
	"snip/circuit"
	"snip/share"
	"testing"
)

// test when server number is 3 and can know h shares of all, test output wire result
func TestGodServer(t *testing.T) {

	// assume that first client construct and share f, g, h
	mod := share.Prime128Value()

	x_array := []*big.Int{big.NewInt(5125), new(big.Int).Sub(mod, big.NewInt(41)), big.NewInt(10), big.NewInt(10), big.NewInt(10)}
	p := 3
	fmt.Println("circuit input: ", x_array)

	a := circuit.FormCircuitByInput(x_array, p, mod)

	// fs, gs, hs are the original client polynomials
	fs, gs, hs := FormFGH(a)
	//hout := hs.ValueAt(big.NewInt(4))
	//fmt.Println("h polynomial at final wire output: ", hout)

	// test original f*g-h
	tmp := new(big.Int).Mul(fs.ValueAt(big.NewInt(11243)), gs.ValueAt(big.NewInt(11243)))
	tmp.Mod(tmp, mod)
	hnum := hs.ValueAt(big.NewInt(11243))
	fmt.Println("tmp: ", tmp, "h: ", hnum)
	tmp.Sub(tmp, hnum)
	tmp.Mod(tmp, mod)
	fmt.Println("final tmp: ", tmp)

	// then client share polynomials h

	hpolys := SharePoly(hs, 2, 3)

	//form and test the recovered output number of shared h polynomials
	xpoint := []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}
	yyarray := []*big.Int{hpolys[0].ValueAt(big.NewInt(4)), hpolys[1].ValueAt(big.NewInt(4)), hpolys[2].ValueAt(big.NewInt(4))}
	re0 := share.Interpolate(big.NewInt(0), x_array, yyarray, mod)
	fmt.Println("recovered output using [h] polynomial final wires: ", re0)

	// construct three servers and each have f, g, h share
	s1, s2, s3 := Server{}, Server{}, Server{}

	s1.Hpoly = hpolys[0]

	s1.CalOutputWire()

	s2.Hpoly = hpolys[1]

	s2.CalOutputWire()

	s3.Hpoly = hpolys[2]

	s3.CalOutputWire()

	// recover secret of identity test and output wire
	xpoint = []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}
	y_array := []*big.Int{s1.OutputWire, s2.OutputWire, s3.OutputWire}

	re1 := share.Interpolate(big.NewInt(0), xpoint, y_array, mod)

	fmt.Println("circuit result is: ", a.Output)
	fmt.Println("mod: ", mod)
	fmt.Println("recovered output wire re1: ", re1)
}

func TestServer_RecoverfgPoly(t *testing.T) {

	// assume that first client construct and share f, g, h
	mod := share.Prime128Value()

	x_array := []*big.Int{big.NewInt(5125), new(big.Int).Sub(mod, big.NewInt(41)), big.NewInt(10), big.NewInt(10), big.NewInt(10)}
	p := 3
	fmt.Println("circuit input: ", x_array)
	a := circuit.FormCircuitByInput(x_array, p, mod)
	fmt.Println("circuit output: ", a.Output)
	// form f, g, h
	f, g, h := FormFGH(a)
	hpolys := SharePoly(h, 2, 5)

	// share x_array
	tmp := ShareToServer(x_array, 2, 5)

	// construct three servers and each have f, g, h share
	s := make([]Server, 5)

	for i := 0; i < 5; i++ {
		s[i].Xshare = tmp[i]
		s[i].Hpoly = hpolys[i]
		s[i].CalOutputWire()
		s[i].Fpoly, s[i].Gpoly = s[i].RecoverfgPoly(p)
	}

	xpoint := []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3), big.NewInt(4), big.NewInt(5)}

	y2 := []*big.Int{s[0].Hpoly.ValueAt(big.NewInt(13)), s[1].Hpoly.ValueAt(big.NewInt(13)), s[2].Hpoly.ValueAt(big.NewInt(13)), s[3].Hpoly.ValueAt(big.NewInt(13)), s[4].Hpoly.ValueAt(big.NewInt(13))}
	y3 := []*big.Int{s[0].Fpoly.ValueAt(big.NewInt(13)), s[1].Fpoly.ValueAt(big.NewInt(13)), s[2].Fpoly.ValueAt(big.NewInt(13)), s[3].Fpoly.ValueAt(big.NewInt(13)), s[4].Fpoly.ValueAt(big.NewInt(13))}
	y4 := []*big.Int{s[0].Gpoly.ValueAt(big.NewInt(13)), s[1].Gpoly.ValueAt(big.NewInt(13)), s[2].Gpoly.ValueAt(big.NewInt(13)), s[3].Gpoly.ValueAt(big.NewInt(13)), s[4].Gpoly.ValueAt(big.NewInt(13))}

	re2 := share.Interpolate(big.NewInt(0), xpoint, y2, mod)
	fmt.Println("h(17): ", re2)

	re3 := share.Interpolate(big.NewInt(0), xpoint, y3, mod)
	fmt.Println("re f17: ", re3)
	fmt.Println("f(17): ", f.ValueAt(big.NewInt(17)))

	re4 := share.Interpolate(big.NewInt(0), xpoint, y4, mod)
	fmt.Println("re g17: ", re4)
	fmt.Println("g(17): ", g.ValueAt(big.NewInt(17)))

	temp := new(big.Int).Sub(new(big.Int).Mul(re3, re4), re2)
	temp = temp.Mod(temp, mod)
	fmt.Println("sub number: ", temp)
}

func TestOutputWire(t *testing.T) {
	// assume that first client construct and share f, g, h
	mod := share.Prime128Value()

	x_array := []*big.Int{big.NewInt(5125), new(big.Int).Sub(mod, big.NewInt(41)), big.NewInt(10), big.NewInt(10), big.NewInt(10)}
	p := 3
	fmt.Println("circuit input: ", x_array)
	a := circuit.FormCircuitByInput(x_array, p, mod)
	fmt.Println("circuit output: ", a.Output)

	// form f, g, h
	_, _, h := FormFGH(a)
	hpolys := SharePoly(h, 2, 3)

	// share x_array
	tmp := ShareToServer(x_array, 2, 3)

	// construct three servers and each have f, g, h share
	s := make([]Server, 3)
	MulGates := circuit.GatesOfType(a, circuit.Gate_Mul)

	for i := 0; i < 3; i++ {
		s[i].Xshare = tmp[i]
		s[i].Hpoly = hpolys[i]
		s[i].FinalWire = big.NewInt(int64(len(MulGates)))
		s[i].CalOutputWire()
		s[i].Fpoly, s[i].Gpoly = s[i].RecoverfgPoly(p)
	}
	fmt.Println("hs(): ", h.ValueAt(big.NewInt(int64(len(MulGates)))))

	xpoints := []*big.Int{}
	ypoints := []*big.Int{}
	for i := 1; i <= 3; i++ {
		xpoints = append(xpoints, big.NewInt(int64(i)))
		ypoints = append(ypoints, s[i-1].OutputWire)
	}
	re := share.Interpolate(big.NewInt(0), xpoints, ypoints, mod)
	fmt.Println("re: ", re)

}
