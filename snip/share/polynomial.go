package share

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"snip/utils"
)

// Poly represents a big integer polynom.
type Poly struct {
	// coeff are the coefficients of the polynom. Coeff[i] is the coefficient of x^i.
	Cof []*big.Int
	// mod is the modulus for polynom arithmetics calculations.
	Mod *big.Int
}

// CofValue returns the i'th coefficient
// Can panic with index out of range when i >= p.Deg().
func (p Poly) CofValue(i int) *big.Int {
	return cp(p.Cof[i])
}

// cp copies a big.Int.
func cp(v *big.Int) *big.Int {
	var u big.Int
	u.Set(v)
	return &u
}

// SetCof sets the i'th coefficient.
// Can panic with index out of range when i >= p.Deg().
func (p Poly) SetCof(i int, v *big.Int) {
	v = cp(v)
	p.Cof[i] = v.Mod(v, p.Mod)
}

// ValueAt returns the y value of the poly on a given x0 value.
func (p Poly) ValueAt(x0 *big.Int) *big.Int {
	val := big.NewInt(0)
	for i := len(p.Cof) - 1; i >= 0; i-- {
		val.Mul(val, x0)
		val.Add(val, p.Cof[i])
		val.Mod(val, p.Mod)
	}
	return val
}

// Degree return poly degree
func (p Poly) Degree() int {
	return len(p.Cof) - 1
}

// RandomPoly returns a new random polynomial of the given degree, which is subjected to arithmetics of
// the given modulus.
// degree t is the times of the highest order
func RandomPoly(degree int64, mod *big.Int) Poly {
	if degree <= 0 {
		panic("deg must be positive number")
	}
	var (
		err   error
		coeff = make([]*big.Int, degree+1)
	)
	for i := range coeff {
		coeff[i], err = rand.Int(rand.Reader, mod)
		if err != nil {
			panic(fmt.Sprintf("creating random int: %s", err))
		}
	}
	return Poly{Cof: coeff, Mod: mod}
}

// ConstructPoly returns a new random polynom of the given degree, which is subjected to arithmetics of
// the given modulus.
func ConstructPoly(cof []*big.Int, mod *big.Int) Poly {
	return Poly{Cof: cof, Mod: mod}
}

// Interpolate returns the y value at x0 of a polynom that lies on points (x[i], y[i]), with modulus
// arithmetics for the given modulus.
func Interpolate(x0 *big.Int, x []*big.Int, y []*big.Int, modulus *big.Int) (y0 *big.Int) {
	if len(x) != len(y) {
		return nil // x and y lists must have the same length.
	}

	nums := make([]*big.Int, len(x))
	dens := make([]*big.Int, len(x))

	for i := range x {
		nums[i] = product(x, x0, i)
		dens[i] = product(x, x[i], i)
	}

	den := product(dens, nil, -1)

	num := big.NewInt(0)
	for i := range nums {
		nums[i].Mul(nums[i], den)
		nums[i].Mul(nums[i], y[i])
		nums[i].Mod(nums[i], modulus)
		v := divmod(nums[i], dens[i], modulus)
		if v == nil {
			return nil // x values are not distinct.
		}
		num.Add(num, v)
	}

	y0 = divmod(num, den, modulus)
	y0.Add(y0, modulus)
	y0.Mod(y0, modulus)
	return y0
}

// product returns the product of vals. If sub is given, the returned product is of (sub-vals[i]).
// If skip is given, the i'th value will be ignored.
func product(vals []*big.Int, sub *big.Int, skip int) *big.Int {
	p := big.NewInt(1)
	for i := range vals {
		if i == skip {
			continue
		}
		v := cp(vals[i])
		if sub != nil {
			v.Sub(sub, v)
		}
		p.Mul(p, v)
	}
	return p
}

// divmod computes num / den modulo mod.
func divmod(a, b, mod *big.Int) *big.Int {
	n := new(big.Int).ModInverse(b, mod)
	if n == nil {
		return nil
	}
	tmp := new(big.Int).Mul(a, n)
	tmp = tmp.Mod(tmp, mod)
	return tmp
}

// MultPoly multiply two polynomials and return the coefficients of the values
func MultPoly(a, b Poly) Poly {

	C := []*big.Int{}
	for i := 0; i < len(a.Cof)+len(b.Cof)-1; i++ {
		C = append(C, big.NewInt(0))
	}

	for i := range a.Cof {
		for j := range b.Cof {
			C[i+j].Add(C[i+j], new(big.Int).Mul(a.Cof[i], b.Cof[j]))
		}
	}

	for i := range C {
		C[i] = C[i].Mod(C[i], a.Mod)
	}

	return Poly{
		Cof: C,
		Mod: a.Mod,
	}
}

// RecoverPolyByNTL recover the polynomial using NTL C++ library
func RecoverPolyByNTL(x []*big.Int, y []*big.Int, modulus *big.Int) Poly {
	if len(x) != len(y) {
		log.Println("x length: ", len(x))
		log.Println("y length: ", len(y))
		log.Fatalln("Construct polynomial failed: x array length not equal to y length")
		return Poly{}
	}
	n := len(x)
	poly := Poly{}
	poly.Mod = modulus

	cof := utils.PolyRecover(n, modulus, x, y)
	poly.Cof = cof
	return poly
}

// ConstructPolyByPoints2 uses lagrange interpolation to gather coefficients of the polynomial
func ConstructPolyByPoints2(x []*big.Int, y []*big.Int, modulus *big.Int) Poly {

	if len(x) != len(y) {
		log.Println("x length: ", len(x))
		log.Println("y length: ", len(y))
		log.Fatalln("Construct polynomial failed: x array length not equal to y length")
		return Poly{}
	}

	// compute k polys and store them in p_array
	p_array := []Poly{}
	for i := 0; i < len(x); i++ {

		// first compute numerator of k-1 amounts
		xc := []*big.Int{}
		for j := 0; j < len(x); j++ {
			if j == i {
				continue
			}
			xc = append(xc, x[j])
		}

		poly1 := MulPolys(xc, modulus)

		// next compute numerator multiply y[i]
		c0 := []*big.Int{y[i]}
		py := Poly{
			Cof: c0,
			Mod: modulus,
		}
		poly1 = MultPoly(poly1, py)

		// third compute poly div & mod denominator
		productl := big.NewInt(1)
		for l := 0; l < len(x); l++ {
			if l == i {
				continue
			}
			tmp := new(big.Int).Sub(x[i], x[l])
			productl = productl.Mul(productl, tmp)
		}

		for m := 0; m < len(poly1.Cof); m++ {
			poly1.Cof[m] = divmod(poly1.Cof[m], productl, modulus)
		}

		p_array = append(p_array, poly1)
	}

	pfinal := AddPolys(p_array)

	for i := 0; i < len(pfinal.Cof); i++ {
		pfinal.Cof[i] = pfinal.Cof[i].Mod(pfinal.Cof[i], modulus)
	}

	return pfinal
}

func MulPolys(x []*big.Int, mod *big.Int) Poly {

	n := len(x)

	if n == 1 {
		return Poly{
			Cof: []*big.Int{new(big.Int).Neg(x[0]), big.NewInt(1)},
			Mod: mod,
		}
	}

	c0 := []*big.Int{new(big.Int).Neg(x[0]), big.NewInt(1)}
	c1 := []*big.Int{new(big.Int).Neg(x[1]), big.NewInt(1)}
	poly0 := Poly{
		Cof: c0,
		Mod: mod,
	}
	poly1 := Poly{
		Cof: c1,
		Mod: mod,
	}
	poly := MultPoly(poly0, poly1)

	if n == 2 {
		return poly
	}

	for i := 2; i < len(x); i++ {
		ci := []*big.Int{new(big.Int).Neg(x[i]), big.NewInt(1)}
		polyi := Poly{
			Cof: ci,
			Mod: mod,
		}
		poly = MultPoly(poly, polyi)
	}

	for i := range poly.Cof {
		poly.Cof[i] = poly.Cof[i].Mod(poly.Cof[i], poly.Mod)
	}

	return poly
}

func AddPolys(p []Poly) Poly {

	n := len(p)
	poly := p[0]
	for i := 1; i < n; i++ {
		poly = AddPoly(poly, p[i])
	}

	return poly
}

func AddPoly(p1, p2 Poly) Poly {

	cof := make([]*big.Int, len(p1.Cof))
	for i := 0; i < len(p1.Cof); i++ {
		cof[i] = new(big.Int).Add(p1.Cof[i], p2.Cof[i])
	}
	poly := Poly{
		Cof: cof,
		Mod: p1.Mod,
	}
	return poly
}

//// ConstructPolyByPoints this is usage of using matrix to construct polynomial, and also uses float64 as W parameters, not recommended
//// we assume the degree of polynomial is t, with t+1 shares
//func ConstructPolyByPoints(x []*big.Int, y []*big.Int, modulus *big.Int) {
//
//	k := len(x)
//
//	x1 := ConvertBigToFloat64(x)
//	y1 := ConvertBigToFloat64(y)
//
//	x_mat_array := []float64{}
//	for i := range x1 {
//		for j := k - 1; j >= 0; j-- {
//			if j == 0 {
//				x_mat_array = append(x_mat_array, float64(1))
//			} else {
//				tmp := math.Pow(x1[i], float64(j))
//				x_mat_array = append(x_mat_array, tmp)
//			}
//		}
//	}
//
//	//fmt.Println("k: ", k, "x_array: ", x_mat_array, "x1: ", x1)
//	xmat := mat.NewDense(k, k, x_mat_array)
//
//	ymat := mat.NewDense(k, 1, y1)
//	//fmt.Println("y matrix: ", ymat)
//
//	var w mat.Dense
//
//	err := w.Solve(xmat, ymat)
//	if err != nil {
//		fmt.Println("no solution for computing W = X^{-1} * Y", err, "X: ", xmat, "Y: ", y)
//	}
//
//	// Print the result using the formatter.
//	fx := mat.Formatted(&w, mat.Prefix("    "), mat.Squeeze())
//	fmt.Printf("w = %.1f", fx)
//
//	// a = -20
//	a := w.At(1, 0)
//	fmt.Println("momo: ", math.Mod(a, 23))
//	bb := math.Float64bits(w.At(1, 0))
//	bbb := int64(bb)
//	now := new(big.Int).Mod(big.NewInt(bbb), big.NewInt(23))
//
//	fmt.Println("now: ", now)
//
//}
//
//func ReverseSlice(s interface{}) {
//	size := reflect.ValueOf(s).Len()
//	swap := reflect.Swapper(s)
//	for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
//		swap(i, j)
//	}
//}

func ConvertBigToFloat64(x []*big.Int) []float64 {
	data := make([]float64, len(x))
	for i := range x {
		data[i], _ = new(big.Float).SetInt(x[i]).Float64()
	}
	return data
}
