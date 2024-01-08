package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readIntcode(filename string) []int {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	text := scanner.Text()
	codes := strings.Split(text, ",")
	intcode := make([]int, len(codes))
	for i, code := range codes {
		intcode[i], err = strconv.Atoi(code)
		if err != nil {
			panic(err)
		}
	}
	return intcode
}

func runAmplifier(intcode []int, phaseSetting, inputSignal int) int {
	ip := 0
	for {
		opcode := intcode[ip] % 100
		modes := intcode[ip] / 100
		param1Mode := modes % 10
		param2Mode := (modes / 10) % 10

		switch opcode {
		case 1:
			param1 := intcode[ip+1]
			param2 := intcode[ip+2]
			resultPos := intcode[ip+3]
			var value1, value2 int
			if param1Mode == 0 {
				value1 = intcode[param1]
			} else {
				value1 = param1
			}
			if param2Mode == 0 {
				value2 = intcode[param2]
			} else {
				value2 = param2
			}
			intcode[resultPos] = value1 + value2
			ip += 4
		case 2:
			param1 := intcode[ip+1]
			param2 := intcode[ip+2]
			resultPos := intcode[ip+3]
			var value1, value2 int
			if param1Mode == 0 {
				value1 = intcode[param1]
			} else {
				value1 = param1
			}
			if param2Mode == 0 {
				value2 = intcode[param2]
			} else {
				value2 = param2
			}
			intcode[resultPos] = value1 * value2
			ip += 4
		case 3:
			resultPos := intcode[ip+1]
			if ip == 0 {
				intcode[resultPos] = phaseSetting
			} else {
				intcode[resultPos] = inputSignal
			}
			ip += 2
		case 4:
			param1 := intcode[ip+1]
			var value int
			if param1Mode == 0 {
				value = intcode[param1]
			} else {
				value = param1
			}
			ip += 2
			return value
		case 99:
			return -1 // Halt
		default:
			panic("Invalid opcode")
		}
	}
}

func main() {
	intcode := readIntcode("day7/input.txt")
	maxOutput := 0

	for phaseSettingSeq := 0; phaseSettingSeq < 120; phaseSettingSeq++ {
		phases := []int{phaseSettingSeq % 5, (phaseSettingSeq / 5) % 5, (phaseSettingSeq / 25) % 5, (phaseSettingSeq / 125) % 5, phaseSettingSeq / 625}
		inputSignal := 0
		for i := 0; i < 5; i++ {
			intcodeCopy := make([]int, len(intcode))
			copy(intcodeCopy, intcode)
			outputSignal := runAmplifier(intcodeCopy, phases[i], inputSignal)
			inputSignal = outputSignal
		}
		if inputSignal > maxOutput {
			maxOutput = inputSignal
		}
	}

	fmt.Println(maxOutput)
}
