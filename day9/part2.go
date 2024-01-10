package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Read input from file
	file, err := os.Open("day9/input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	program := strings.Split(scanner.Text(), ",")

	// Convert input to integers and extend memory
	memory := make(map[int]int)
	for i, s := range program {
		value, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		memory[i] = value
	}

	// Run the Intcode program in sensor boost mode
	fmt.Println(runIntcode(memory, 2))
}

func runIntcode(memory map[int]int, input int) int {
	var output, ip, relativeBase int

	for {
		opcode := memory[ip] % 100
		modes := strconv.Itoa(memory[ip] / 100)

		getParam := func(offset int) int {
			mode := 0
			if len(modes) >= offset {
				mode, _ = strconv.Atoi(string(modes[len(modes)-offset]))
			}

			param := memory[ip+offset]
			switch mode {
			case 0:
				return memory[param]
			case 1:
				return param
			case 2:
				return memory[relativeBase+param]
			default:
				panic("unknown parameter mode")
			}
		}

		setParam := func(offset int, value int) {
			mode := 0
			if len(modes) >= offset {
				mode, _ = strconv.Atoi(string(modes[len(modes)-offset]))
			}

			param := memory[ip+offset]
			switch mode {
			case 0:
				memory[param] = value
			case 2:
				memory[relativeBase+param] = value
			default:
				panic("unknown parameter mode")
			}
		}

		switch opcode {
		case 1:
			setParam(3, getParam(1)+getParam(2))
			ip += 4
		case 2:
			setParam(3, getParam(1)*getParam(2))
			ip += 4
		case 3:
			setParam(1, input) // The input is now 2 for sensor boost mode
			ip += 2
		case 4:
			output = getParam(1)
			ip += 2
		case 5:
			if getParam(1) != 0 {
				ip = getParam(2)
			} else {
				ip += 3
			}
		case 6:
			if getParam(1) == 0 {
				ip = getParam(2)
			} else {
				ip += 3
			}
		case 7:
			if getParam(1) < getParam(2) {
				setParam(3, 1)
			} else {
				setParam(3, 0)
			}
			ip += 4
		case 8:
			if getParam(1) == getParam(2) {
				setParam(3, 1)
			} else {
				setParam(3, 0)
			}
			ip += 4
		case 9:
			relativeBase += getParam(1)
			ip += 2
		case 99:
			return output // Final output
		default:
			panic(fmt.Sprintf("unknown opcode: %d", opcode))
		}
	}
}
