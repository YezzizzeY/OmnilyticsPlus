package utils

import (
	"fmt"
	"math"
	"scientificgo.org/fft"
)

func CalculateFFT() {
	// Given data points
	yValues := []complex128{7, 6, 0, 0} // Padded to the next power of 2

	// Compute the inverse FFT to get the coefficients of the polynomial
	coefficients := fft.Fft(yValues, true)

	// Normalize the coefficients by dividing by the length
	n := len(yValues)
	for i := range coefficients {
		coefficients[i] /= complex(float64(n), 0)
	}

	// Convert coefficients to integer representation
	intCoefficients := make([]int, n)
	for i, coeff := range coefficients {
		intCoefficients[i] = int(math.Round(real(coeff)))
	}

	// Print the integer coefficients
	fmt.Println("Interpolated polynomial integer coefficients:", intCoefficients)
}
