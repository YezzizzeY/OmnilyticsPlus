package main

import (
	"math/big"
	"snip/circuit"
	"snip/share"
)

func ShareToServer(a_array []*big.Int, t, k int) [][]*big.Int {
	// in this program, we share it with a polynomial of degree t and generate k shares
	x_num := len(a_array)
	shares_array := []share.Shares{}
	for i := 0; i < x_num; i++ {
		shares := share.ShareSecret(a_array[i], int64(t), int64(k))
		shares_array = append(shares_array, shares)
	}
	var servers [][]*big.Int
	for j := 0; j < k; j++ {
		tmp := []*big.Int{}
		for i := 0; i < x_num; i++ {
			tmp = append(tmp, shares_array[i].ShareValues[j].Y)
		}
		servers = append(servers, tmp)
	}

	return servers
}

// ShareXarray for every secret x, we share it using shamir secret sharing and of degree t and k shares
func ShareXarray(a_array []*big.Int, t, k int) []share.Shares {

	// in this program, we share it with a polynomial of degree t and generate k shares
	x_num := len(a_array)
	shares_array := []share.Shares{}
	for i := 0; i < x_num; i++ {
		shares := share.ShareSecret(a_array[i], int64(t), int64(k))
		shares_array = append(shares_array, shares)
	}

	return shares_array
}

// FormFGH calculate the polynomials of f, g, h in SNIP scheme of one client
func FormFGH(c circuit.Circuit) (share.Poly, share.Poly, share.Poly) {

	MulGates := circuit.GatesOfType(c, circuit.Gate_Mul)

	MulGatesNum := len(MulGates)

	var f_array []*big.Int
	var g_array []*big.Int

	for i := 0; i < MulGatesNum; i++ {
		temp := MulGates[i]

		f_array = append(f_array, temp.LChild.OutputValue)
		g_array = append(g_array, temp.RChild.OutputValue)

	}
	x_array := []*big.Int{}
	for i := 0; i < len(f_array); i++ {
		x_array = append(x_array, big.NewInt(int64(i)+1))
	}
	fPoly := share.RecoverPolyByNTL(x_array, f_array, share.Prime128Value())
	gPoly := share.RecoverPolyByNTL(x_array, g_array, share.Prime128Value())
	hPoly := share.MultPoly(fPoly, gPoly)

	return fPoly, gPoly, hPoly

}

// SharePoly shares a poly, with another polynomial of degree t and generate k polys: for example,
// when k=3, it generates 3 shares of every coefficient, thus another 3 polynomials f g h
func SharePoly(p share.Poly, t, k int) []share.Poly {

	pshares := ShareXarray(p.Cof, t, k)

	polys := []share.Poly{}

	// in this loop, i is the i'th poly share, j is j'th cof
	for i := 0; i < k; i++ {

		pt := share.Poly{}
		pt.Mod = p.Mod
		for j := 0; j < p.Degree()+1; j++ {
			pt.Cof = append(pt.Cof, pshares[j].Ys[i])
		}

		polys = append(polys, pt)
	}

	return polys

}
