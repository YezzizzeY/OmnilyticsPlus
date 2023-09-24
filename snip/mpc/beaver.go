package mpc

import (
	"crypto/rand"
	"log"
	"math/big"
	"snip/share"
)

type Beaver struct {
	D      *big.Int
	E      *big.Int
	Ashare share.Share
	Bshare share.Share
	Cshare share.Share
}

// GenerateBeaver given share of original number alpha and beta
func GenerateBeaver(alpha, beta *big.Int, mod *big.Int, t, k int64) []Beaver {

	// first select random a an b
	a, err1 := rand.Int(rand.Reader, mod)
	if err1 != nil {
		log.Fatalln("generate random number a in beaver failed")
	}

	b, err2 := rand.Int(rand.Reader, mod)
	if err2 != nil {
		log.Fatalln("generate random number b in beaver failed")
	}

	// compute c
	c := new(big.Int).Mul(a, b)
	c.Mod(c, mod)

	// share a, b and c
	ashares := share.ShareSecret(a, t, k)
	bshares := share.ShareSecret(b, t, k)
	cshares := share.ShareSecret(c, t, k)

	// compute d and e
	d := new(big.Int).Sub(alpha, a)
	e := new(big.Int).Sub(beta, b)

	// construct k beaver
	beavers := []Beaver{}
	for i := 0; i < int(k); i++ {
		be := Beaver{}
		be.D = d
		be.E = e
		be.Ashare = ashares.ShareValues[i]
		be.Bshare = bshares.ShareValues[i]
		be.Cshare = cshares.ShareValues[i]
		beavers = append(beavers, be)
	}

	return beavers

}

func CalculateMul(beaver Beaver) share.Share {

	mul := new(big.Int).Mul(beaver.D, beaver.E)
	mul = mul.Add(mul, new(big.Int).Mul(beaver.D, beaver.Bshare.Y))
	mul = mul.Add(mul, new(big.Int).Mul(beaver.E, beaver.Ashare.Y))
	mul = mul.Add(mul, beaver.Cshare.Y)
	s := share.Share{
		X: beaver.Ashare.X,
		Y: mul,
	}

	return s
}
