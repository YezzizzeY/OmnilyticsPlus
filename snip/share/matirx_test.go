package share

import (
	"fmt"

	"testing"
)

func Test1(t *testing.T) {

	X := Matrix{}

	row := []float64{0, 0, 1, 1, 1, 1, 4, 2, 1}

	X.Set(3, 3, row)

	fmt.Println("X: ", X)

	Y := Matrix{}
	Y.Set(1, 3, []float64{2, 7, 16})

	fmt.Println("Y: ", Y)

	x_inverse := Inverse(X)

	fmt.Println("X_inverse: ", x_inverse)

	Result := Mul(x_inverse, Y)
	fmt.Println("result: ", Result)

	Result = Transpose(Result)
	fmt.Println("Transparent result: ", Result.Data[0])

}
