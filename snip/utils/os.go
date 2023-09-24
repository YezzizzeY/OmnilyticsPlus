package utils

import (
	"fmt"
	"math/big"
	"os/exec"
	"strings"
)

func PolyRecover(n int, p *big.Int, x []*big.Int, y []*big.Int) []*big.Int {
	// Executable file path
	cppExecutable := "/home/yezzi/Desktop/Omnilytics+/c++/test"

	// Construct the command and its arguments
	cmd := exec.Command(cppExecutable)

	// Prepare input data
	inputData := fmt.Sprintf("%s\n%d\n", p.String(), n)

	for i := range x {
		inputData += x[i].String() + " "
	}
	inputData += "\n"

	for i := range y {
		inputData += y[i].String() + " "
	}

	// Pass input data to the command's standard input
	cmd.Stdin = strings.NewReader(inputData)

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error during the execution of the C++ program:", err)
	}

	// Parse the output into []*big.Int
	outputStr := string(output)
	outputStr = strings.ReplaceAll(outputStr, "[", "")
	outputStr = strings.ReplaceAll(outputStr, "]", "")
	outputStr = strings.TrimSpace(outputStr)

	numStrs := strings.Split(outputStr, " ")

	// Create a []*big.Int slice and convert strings to big.Int
	var result []*big.Int
	for _, numStr := range numStrs {
		num := new(big.Int)
		_, ok := num.SetString(numStr, 10)
		if !ok {
			fmt.Printf("Unable to parse string %s as big.Int\n", numStr)
		}
		result = append(result, num)
	}

	// Print the result
	return result
}
