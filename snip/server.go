package main

import (
	"math/big"
	"snip/mpc"
	"snip/share"
)

type Server struct {
	Xshare     []*big.Int
	Fpoly      share.Poly
	Gpoly      share.Poly
	Hpoly      share.Poly
	SubNum     *big.Int
	OutputWire *big.Int
	BeaverS    mpc.Beaver
	FinalWire  *big.Int
}

// CalculateSub server calculate sub num of shares f, g and h
func (s *Server) CalculateSub(point *big.Int) {

	a := mpc.CalculateMul(s.BeaverS).Y
	s.SubNum = new(big.Int).Sub(a, s.Hpoly.ValueAt(point))
	s.SubNum = s.SubNum.Mod(s.SubNum, s.Fpoly.Mod)

}

// CalOutputWire server calculate outputwire
func (s *Server) CalOutputWire() {
	s.OutputWire = s.Hpoly.ValueAt(s.FinalWire)
}

// RecoverfgPoly uses h share of h polynomial and share x to recover f and g polynomial, pnum is
// the num of |valid|
// we have len(xshare_array) mul gates of pow and p-1 mul gates of sub
func (s *Server) RecoverfgPoly(pnum int) (share.Poly, share.Poly) {

	// first calculate first l  mul gates' inputs
	fshare_array := []*big.Int{}
	gshare_array := []*big.Int{}

	l := len(s.Xshare)

	for i := 0; i < l; i++ {
		fshare_array = append(fshare_array, s.Xshare[i])
		gshare_array = append(gshare_array, s.Xshare[i])
	}

	// addPow is the num of x1^2+x2^2+...
	addPow := big.NewInt(0)
	for i := 0; i < l; i++ {
		addPow = addPow.Add(addPow, s.Hpoly.ValueAt(big.NewInt(int64(i+1))))
	}

	// next rest g and h points
	// with c is pownum
	// first calculate c-1 and c-2, next loop fill in
	fshare_array = append(fshare_array, new(big.Int).Sub(addPow, big.NewInt(1)))
	gshare_array = append(gshare_array, new(big.Int).Sub(addPow, big.NewInt(2)))

	for i := 1; i <= pnum-2; i++ {
		fshare_array = append(fshare_array, s.Hpoly.ValueAt(big.NewInt(int64(l+i))))
		tmp := new(big.Int).Sub(addPow, big.NewInt(int64(i+2)))
		tmp.Mod(tmp, s.Hpoly.Mod)
		gshare_array = append(gshare_array, tmp)
	}

	// recover f and g using interpolation
	xpoint := []*big.Int{}
	for i := 1; i <= (l + pnum - 1); i++ {
		xpoint = append(xpoint, big.NewInt(int64(i)))
	}
	fpoly := share.RecoverPolyByNTL(xpoint, fshare_array, s.Hpoly.Mod)
	gpoly := share.RecoverPolyByNTL(xpoint, gshare_array, s.Hpoly.Mod)

	return fpoly, gpoly
}
