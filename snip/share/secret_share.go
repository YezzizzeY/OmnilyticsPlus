package share

import (
	"math/big"
)

var prime128 = Prime128Value()

func Prime128Value() *big.Int {
	p := big.NewInt(2)
	p.Exp(p, big.NewInt(127), nil)
	p.Sub(p, big.NewInt(1))
	return p
}

// Share is a part of a secret.
type Share struct {
	X, Y *big.Int
}

// Shares is a part of a secret.
type Shares struct {
	ShareValues []Share
	Length      int
	Xs          []*big.Int
	Ys          []*big.Int
}

// ShareSecret creates k shares. The secret argument is optional. It returns the shares and the
// secret for the shares.
// t is the degree of polynomial, we need at least t+1 shares to recover secret
// p = 1.70141183E38
func ShareSecret(secret *big.Int, t, k int64) Shares {
	if k < t+1 {
		panic("irrecoverable: not enough shares to reconstruct the secret.")
	}
	if k <= 0 {
		panic("number of shares must be positive.")
	}
	p := RandomPoly(t, prime128)

	// Set the first coefficient to the secret (the value at x=0) if the secret was given. And
	// anyway store the first coefficient in the secret variable.
	if secret != nil {
		if secret.Cmp(prime128) > 0 {
			panic("secret value is too big (must be lower than 2^127 - 1)")
		}
		p.SetCof(0, secret)
	}
	secret = p.CofValue(0)

	// Create the shares which are the value of p at any point but x != 0. Choose x in [1..n].
	shares := make([]Share, 0, k)
	for i := int64(1); i <= k; i++ {
		x := big.NewInt(i)
		y := p.ValueAt(x)
		shares = append(shares, Share{X: x, Y: y})
	}

	var shareArray Shares
	shareArray.ShareValues = shares
	shareArray.Length = len(shares)
	for i := int64(0); i < k; i++ {
		shareArray.Xs = append(shareArray.Xs, shareArray.ShareValues[i].X)
		shareArray.Ys = append(shareArray.Ys, shareArray.ShareValues[i].Y)
	}

	return shareArray
}

// RecoverSecret the secret from shares. Notice that the number of shares that is used should be at least
// the recover amount (k) that was used in order to create them in the New function.
func RecoverSecret(s Shares) (secret *big.Int) {
	// Evaluate the polynom that goes through all (x[i], y[i]) points at x=0.
	return Interpolate(big.NewInt(0), s.Xs, s.Ys, prime128)
}

// Recovershares the secret from shares. Notice that the number of shares that is used should be at least
// the recover amount (k) that was used in order to create them in the New function.
func Recovershares(s []Share) (secret *big.Int) {
	xs := []*big.Int{}
	ys := []*big.Int{}
	for _, v := range s {
		xs = append(xs, v.X)
		ys = append(ys, v.Y)
	}
	return Interpolate(big.NewInt(0), xs, ys, prime128)
}
