package circuit

import (
	"math/big"
)

// Supported gate types.
// In the first step of achieving SNIP, we only use Gate_AddConst and Gate_Mul
const (
	Gate_Input    = iota
	Gate_Add      = iota
	Gate_AddConst = iota
	Gate_Mul      = iota
	Gate_MulConst = iota
)

// Gate represents a gate in an arithmetic circuit, and possibly
// holds the value on its output wire.
type Gate struct {
	GateType int

	LChild *Gate
	RChild *Gate

	OutputValue *big.Int
}

// Circuit represents an arithmetic circuit over a particular finite
// note: in this circuit we have only one output
type Circuit struct {
	Mod    *big.Int
	Gates  []*Gate
	Root   *Gate
	Output *big.Int
}

func FormMulGate(L *Gate, R *Gate, mod *big.Int) *Gate {

	var g = new(Gate)

	g.GateType = Gate_Mul

	g.LChild = L
	g.RChild = R

	// calculate output value
	left := L.OutputValue
	right := R.OutputValue
	g.OutputValue = new(big.Int).Mul(left, right)
	g.OutputValue = g.OutputValue.Mod(g.OutputValue, mod)

	return g
}

func FormAddGate(L *Gate, R *Gate, mod *big.Int) *Gate {

	var g = new(Gate)

	g.GateType = Gate_Add

	g.LChild = L
	g.RChild = R

	// calculate output value
	left := L.OutputValue
	right := R.OutputValue
	g.OutputValue = new(big.Int).Add(left, right)
	g.OutputValue = g.OutputValue.Mod(g.OutputValue, mod)

	return g
}

func FormInputGate(value *big.Int, mod *big.Int) *Gate {

	var g = new(Gate)

	g.GateType = Gate_Input

	g.LChild = nil
	g.RChild = nil

	g.OutputValue = value
	g.OutputValue = g.OutputValue.Mod(g.OutputValue, mod)

	return g
}

func FormAddConstGate(L *Gate, value *big.Int, mod *big.Int) *Gate {

	var g = new(Gate)

	g.GateType = Gate_AddConst

	g.LChild = L
	g.RChild = nil

	value = value.Mod(value, mod)
	// calculate the output value of Gate AddConst
	g.OutputValue = new(big.Int).Add(L.OutputValue, value)
	g.OutputValue = g.OutputValue.Mod(g.OutputValue, mod)

	return g
}

func GatesOfType(c Circuit, t int) []*Gate {
	gates := make([]*Gate, 0)
	for _, gate := range c.Gates {
		switch gate.GateType {
		case t:
			gates = append(gates, gate)
		}
	}
	return gates
}

func cpGate(gate *Gate) *Gate {
	g2 := new(Gate)

	g2.GateType = gate.GateType

	g2.LChild = gate.LChild
	g2.RChild = gate.RChild

	g2.OutputValue = gate.OutputValue

	return g2
}

// FormPowAddCircuit Form a circuit that is
func FormCircuitByInput(x_array []*big.Int, p int, mod *big.Int) Circuit {

	// InputGatesX is the input gates containing x1, x2, ... xn
	var InputGatesX []*Gate
	num := len(x_array)
	for i := 0; i < num; i++ {
		temp := FormInputGate(x_array[i], mod)
		InputGatesX = append(InputGatesX, temp)
	}

	// PowGateArray is the mul gates of x_i^2
	var PowGateArray []*Gate
	for i := 0; i < num; i++ {
		temp := FormMulGate(InputGatesX[i], InputGatesX[i], mod)
		PowGateArray = append(PowGateArray, temp)
	}

	// AddGateArray is the add gate of x_i^2 + x_i+1^2
	var AddGateArray []*Gate
	for i := 0; i < num-1; i++ {
		if i == 0 {
			t := FormAddGate(PowGateArray[0], PowGateArray[1], mod)
			AddGateArray = append(AddGateArray, t)
			continue
		}
		temp := FormAddGate(AddGateArray[i-1], PowGateArray[i+1], mod)
		AddGateArray = append(AddGateArray, temp)
	}

	AddGateRoot := AddGateArray[num-2]

	// Add2 is c1-1, c1-2, ...
	// Mul2 is (c1-1)(c1-2)...
	var Add2 []*Gate
	var Mul2 []*Gate

	for i := 0; i < p; i++ {
		temp := FormAddConstGate(cpGate(AddGateRoot), big.NewInt(int64(-i-1)), mod)
		Add2 = append(Add2, temp)
	}

	for i := 0; i < p-1; i++ {
		if i == 0 {
			t := FormMulGate(Add2[0], Add2[1], mod)
			Mul2 = append(Mul2, t)
			continue
		}
		temp := FormMulGate(Mul2[i-1], Add2[i+1], mod)
		Mul2 = append(Mul2, temp)
	}
	Root := Mul2[p-2]

	circuit := Circuit{}
	circuit.Gates = append(circuit.Gates, InputGatesX...)
	circuit.Gates = append(circuit.Gates, PowGateArray...)
	circuit.Gates = append(circuit.Gates, AddGateArray...)
	circuit.Gates = append(circuit.Gates, Add2...)
	circuit.Gates = append(circuit.Gates, Mul2...)
	circuit.Mod = mod
	circuit.Root = Root
	circuit.Output = Root.OutputValue

	return circuit
}
