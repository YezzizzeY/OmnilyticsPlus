package share

import (
	"fmt"
	"math/big"
	"testing"
)

func TestRand(t *testing.T) {
	poly := RandomPoly(4, big.NewInt(23))
	fmt.Println(poly)
}

func TestConstruct(t *testing.T) {
	cof := []*big.Int{big.NewInt(2), big.NewInt(3), big.NewInt(2)}
	mod := big.NewInt(23)
	poly := ConstructPoly(cof, mod)
	fmt.Println(poly)
	fmt.Println("poly degree: ", poly.Degree())
	value := poly.ValueAt(big.NewInt(0))
	fmt.Println("y value at x=0 is ", value)

	value2 := poly.ValueAt(big.NewInt(1))
	fmt.Println("y value at x=1 is ", value2)

	value3 := poly.ValueAt(big.NewInt(2))
	fmt.Println("y value at x=2 is ", value3)
}

func TestMultPoly(t *testing.T) {

	cofA := []*big.Int{big.NewInt(5), big.NewInt(2), big.NewInt(3)}
	cofB := []*big.Int{big.NewInt(2), big.NewInt(1), big.NewInt(5)}

	polyA := Poly{
		Cof: cofA,
		Mod: big.NewInt(97),
	}

	polyB := Poly{
		Cof: cofB,
		Mod: big.NewInt(97),
	}

	PolyC := MultPoly(polyA, polyB)

	fmt.Println("PolyC: ", PolyC)
}

func TestConstructPolyByPoints(t *testing.T) {

	x_array := []*big.Int{big.NewInt(1), big.NewInt(3), big.NewInt(4)}
	y_array := []*big.Int{big.NewInt(7), big.NewInt(6), big.NewInt(0)}

	ConstructPolyByPoints2(x_array, y_array, big.NewInt(23))

}

func TestSmall(t *testing.T) {

	x1 := big.NewInt(5)
	y1 := big.NewInt(11)

	x1 = x1.ModInverse(x1, big.NewInt(23))
	out := new(big.Int).Mul(x1, y1)
	out = out.Mod(out, big.NewInt(23))
	fmt.Println("out: ", out)
}

func TestConstructPolyByPoints2(t *testing.T) {

	x_array := []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(4)}
	y_array := []*big.Int{big.NewInt(7), big.NewInt(16), big.NewInt(0)}

	p := ConstructPolyByPoints2(x_array, y_array, big.NewInt(23))
	fmt.Println(p)
}

func TestWholePoly(t *testing.T) {

	mod := Prime128Value()
	//mod = big.NewInt(23)
	fmt.Println("mod: ", mod)

	poly := RandomPoly(5, mod)
	//poly = ConstructPoly([]*big.Int{big.NewInt(1), big.NewInt(2)}, mod)
	fmt.Println("polynomial: ", poly)

	var x []*big.Int
	var y []*big.Int
	for i := 1; i <= 6; i++ {
		tmp := poly.ValueAt(big.NewInt(int64(i)))
		x = append(x, big.NewInt(int64(i)))
		y = append(y, tmp)
	}
	fmt.Println("x: ", x)
	fmt.Println("y: ", y)

	s := Interpolate(big.NewInt(0), x, y, mod)
	fmt.Println("recovered secret: ", s)

	p := ConstructPolyByPoints2(x, y, mod)
	fmt.Println("recovered polynomial: ", p)
}
