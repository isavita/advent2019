package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("day5/input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	programStr := scanner.Text()
	program := []int{}
	for _, s := range strings.Split(programStr, ",") {
		i, _ := strconv.Atoi(s)
		program = append(program, i)
	}

	input := 5
	output := 0
	i := 0
	for {
		opcode := program[i] % 100
		modes := program[i] / 100
		param1Mode := modes % 10
		modes /= 10
		param2Mode := modes % 10

		switch opcode {
		case 1:
			p1 := getValue(program, i+1, param1Mode)
			p2 := getValue(program, i+2, param2Mode)
			p3 := program[i+3]
			program[p3] = p1 + p2
			i += 4
		case 2:
			p1 := getValue(program, i+1, param1Mode)
			p2 := getValue(program, i+2, param2Mode)
			p3 := program[i+3]
			program[p3] = p1 * p2
			i += 4
		case 3:
			program[program[i+1]] = input
			i += 2
		case 4:
			output = getValue(program, i+1, param1Mode)
			fmt.Println(output)
			i += 2
		case 5:
			p1 := getValue(program, i+1, param1Mode)
			p2 := getValue(program, i+2, param2Mode)
			if p1 != 0 {
				i = p2
			} else {
				i += 3
			}
		case 6:
			p1 := getValue(program, i+1, param1Mode)
			p2 := getValue(program, i+2, param2Mode)
			if p1 == 0 {
				i = p2
			} else {
				i += 3
			}
		case 7:
			p1 := getValue(program, i+1, param1Mode)
			p2 := getValue(program, i+2, param2Mode)
			p3 := program[i+3]
			if p1 < p2 {
				program[p3] = 1
			} else {
				program[p3] = 0
			}
			i += 4
		case 8:
			p1 := getValue(program, i+1, param1Mode)
			p2 := getValue(program, i+2, param2Mode)
			p3 := program[i+3]
			if p1 == p2 {
				program[p3] = 1
			} else {
				program[p3] = 0
			}
			i += 4
		case 99:
			return
		default:
			panic("Invalid opcode")
		}
	}
}

func getValue(program []int, pos int, mode int) int {
	if mode == 0 {
		return program[program[pos]]
	} else {
		return program[pos]
	}
}
