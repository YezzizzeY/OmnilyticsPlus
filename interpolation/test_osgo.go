package main

import (
	"fmt"
	"math/big"
	"os/exec"
	"strings"
)

func main() {
	// Executable file path
	cppExecutable := "./test"

	// Input data
	p := big.NewInt(23)
	n := 3
	x := []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}
	y := []*big.Int{big.NewInt(7), big.NewInt(16), big.NewInt(6)}

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
		return
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
			return
		}
		result = append(result, num)
	}
	
	// Print the result
	fmt.Println(result)
}
