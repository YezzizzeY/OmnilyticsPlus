package utils

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"strings"
)

func ReadFile(path string, lineByLine bool) (error, [][]*big.Int) {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err, [][]*big.Int{}
	}
	defer file.Close()

	var numbers [][]*big.Int

	// 用于乘以解析的浮点数以获得整数表示
	multiplier := new(big.Float).SetFloat64(1e+18)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// 使用空格分隔每一行
		line := scanner.Text()
		parts := strings.Fields(line)

		var lineNumbers []*big.Int
		for _, part := range parts {
			// 将字符串转换为 *big.Float
			floatNum, _, err := big.ParseFloat(part, 10, 0, big.ToNearestEven)

			if err != nil {
				fmt.Println("Error parsing float:", err)
				continue
			}

			// 乘以multiplier以获得整数表示
			floatNum.Mul(floatNum, multiplier)

			// 将 *big.Float 转换为 *big.Int
			intNum := new(big.Int)
			floatNum.Int(intNum)

			lineNumbers = append(lineNumbers, intNum)
		}

		if lineByLine {
			numbers = append(numbers, lineNumbers)
		} else {
			if len(numbers) == 0 {
				numbers = append(numbers, []*big.Int{})
			}
			numbers[0] = append(numbers[0], lineNumbers...)
		}
	}

	if err := scanner.Err(); err != nil {
		return err, [][]*big.Int{}
	}

	return nil, numbers
}

// IntegerPart 获取 big.Float 的整数部分。
func IntegerPart(f *big.Float) *big.Int {
	intPart := new(big.Int)
	f.Int(intPart)
	return intPart
}
