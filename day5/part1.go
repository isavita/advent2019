package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func getMode(instruction, position int) int {
	return instruction / int(math.Pow10(position+1)) % 10
}

func getParam(program []int, pointer, mode int) int {
	if mode == 0 {
		return program[program[pointer]]
	}
	return program[pointer]
}

func runProgram(program []int, input int) int {
	output := 0
	for pointer := 0; pointer < len(program); {
		instruction := program[pointer]
		opcode := instruction % 100

		switch opcode {
		case 1, 2:
			param1 := getParam(program, pointer+1, getMode(instruction, 1))
			param2 := getParam(program, pointer+2, getMode(instruction, 2))
			result := 0
			if opcode == 1 {
				result = param1 + param2
			} else {
				result = param1 * param2
			}
			program[program[pointer+3]] = result
			pointer += 4
		case 3:
			program[program[pointer+1]] = input
			pointer += 2
		case 4:
			output = getParam(program, pointer+1, getMode(instruction, 1))
			pointer += 2
		case 99:
			return output
		default:
			panic(fmt.Sprintf("Unknown opcode: %d", opcode))
		}
	}
	return output
}

func main() {
	data, err := os.ReadFile("day5/input.txt")
	if err != nil {
		panic(err)
	}
	strProgram := strings.Split(strings.TrimSpace(string(data)), ",")
	program := make([]int, len(strProgram))
	for i, s := range strProgram {
		program[i], err = strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(runProgram(program, 1))
}
