package main

import (
	"fmt"
	"math/big"
	"snip/circuit"
	"snip/mpc"
	"snip/share"
	"snip/utils"
	"time"
)

func main() {
	fmt.Println(">> Omnilytics plus <<")
	tfinal := time.Now()
	t := time.Now()

	err, numbers := utils.ReadFile("data/betafc1.txt", false)
	if err != nil {
		panic(err)
	}
	var gradient []*big.Int
	gradient = scaleAndTruncate(numbers[0])

	sum := new(big.Int)
	for _, val := range gradient {
		square := new(big.Int).Mul(val, val)
		sum.Add(sum, square)
	}

	fmt.Println("The sum of squares is:", sum)
	//for _, row := range numbers {
	//	for _, val := range row {
	//		gradient = append(gradient, val)
	//	}
	//}
	// this is the setup of no-text input, we have 5 servers and one client in memory space,
	// the input of client x is [3,3,5] and p=40
	client_x_array := gradient
	fmt.Println("client array length: ", len(client_x_array))
	p := 1000

	// assume that first client construct and share f, g, h
	mod := share.Prime128Value()
	fmt.Println("mod: ", mod)

	a := circuit.FormCircuitByInput(client_x_array, p, mod)
	fmt.Println("circuit output: ", a.Output)

	elapsed := time.Since(t)
	fmt.Println("Form circuit run time", elapsed)

	t2 := time.Now()
	// form f, g, h
	f, g, h := FormFGH(a)
	fmt.Println("form fgh succeed, polynomial degree: ", len(f.Cof))

	elapsed = time.Since(t2)
	fmt.Println("Form fgh run time", elapsed)

	t3 := time.Now()
	hpolys := SharePoly(h, 2, 5)
	fmt.Println("share poly succeed")
	elapsed = time.Since(t3)
	fmt.Println("Share poly run time", elapsed)

	// share x_array
	tmp := ShareToServer(client_x_array, 4, 5)

	// generate beavers with r=17
	alpha := f.ValueAt(big.NewInt(17))
	beta := g.ValueAt(big.NewInt(17))
	beavers := mpc.GenerateBeaver(alpha, beta, mod, 4, 5)
	MulGates := circuit.GatesOfType(a, circuit.Gate_Mul)
	fmt.Println("mul gates num: ", len(MulGates))

	// construct three servers and each have f, g, h share
	s := make([]Server, 5)

	for i := 0; i < 5; i++ {
		s[i].Xshare = tmp[i]
		s[i].Hpoly = hpolys[i]
		s[i].FinalWire = big.NewInt(int64(len(MulGates)))
		s[i].CalOutputWire()
		t4 := time.Now()
		s[i].Fpoly, s[i].Gpoly = s[i].RecoverfgPoly(p)
		elapsed = time.Since(t4)
		fmt.Println("Recover poly run time", elapsed)
		s[i].BeaverS = beavers[i]
		t5 := time.Now()
		s[i].CalculateSub(big.NewInt(17))
		elapsed = time.Since(t5)
		fmt.Println("Calculate R time", elapsed)
	}

	// recover output wire and r value from servers
	xpoint := []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3), big.NewInt(4), big.NewInt(5)}

	outputShares := []*big.Int{}
	rShares := []*big.Int{}

	for i := 0; i < 5; i++ {
		outputShares = append(outputShares, s[i].OutputWire)
		rShares = append(rShares, s[i].SubNum)
	}
	fmt.Println("Output shares in smart contract: ", outputShares)
	fmt.Println("R shares in smart contract: ", rShares)

	output := share.Interpolate(big.NewInt(0), xpoint, outputShares, mod)
	r := share.Interpolate(big.NewInt(0), xpoint, rShares, mod)

	finaltime := time.Since(tfinal)
	fmt.Println("total time of one client: ", finaltime)
	fmt.Println("Output: ", output, "r: ", r)
}

func scaleAndTruncate(nums []*big.Int) []*big.Int {
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(20), nil)
	var scaledNums []*big.Int

	for _, num := range nums {
		// Divide the number by 10^20
		scaledNum := new(big.Int).Div(num, divisor)
		scaledNums = append(scaledNums, scaledNum)
	}
	return scaledNums
}
