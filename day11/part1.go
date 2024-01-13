package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PanelColor int

const (
	Black PanelColor = 0
	White PanelColor = 1
)

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

type Position struct {
	X, Y int
}

type Robot struct {
	Position
	Direction
}

type Grid map[Position]PanelColor

func (r *Robot) turnAndMove(turnDirection int) {
	if turnDirection == 0 {
		r.Direction = (r.Direction + 3) % 4 // Turn left
	} else {
		r.Direction = (r.Direction + 1) % 4 // Turn right
	}

	switch r.Direction {
	case Up:
		r.Y--
	case Right:
		r.X++
	case Down:
		r.Y++
	case Left:
		r.X--
	}
}

type Intcode struct {
	memory []int
	ip     int
	input  []int
	output []int
	halted bool
}

func NewIntcode(program []int) *Intcode {
	memory := make([]int, len(program))
	copy(memory, program)
	return &Intcode{memory: memory}
}

func (ic *Intcode) AddInput(input int) {
	ic.input = append(ic.input, input)
}

func (ic *Intcode) Run() {
	ic.output = []int{}
	for {
		opcode := ic.memory[ic.ip] % 100
		switch opcode {
		case 1, 2, 7, 8: // Add, Multiply, Less than, Equals
			ic.ensureMemory(ic.ip + 3)
			params := ic.getParams(3)
			val1, val2 := ic.readMemory(params[0]), ic.readMemory(params[1])
			if opcode == 1 {
				ic.writeMemory(params[2], val1+val2)
			} else if opcode == 2 {
				ic.writeMemory(params[2], val1*val2)
			} else if opcode == 7 && val1 < val2 || opcode == 8 && val1 == val2 {
				ic.writeMemory(params[2], 1)
			} else {
				ic.writeMemory(params[2], 0)
			}
			ic.ip += 4
		case 3, 4: // Input, Output
			ic.ensureMemory(ic.ip + 1)
			params := ic.getParams(1)
			if opcode == 3 {
				if len(ic.input) == 0 {
					return // Wait for more input
				}
				ic.writeMemory(params[0], ic.input[0])
				ic.input = ic.input[1:]
			} else {
				ic.output = append(ic.output, ic.readMemory(params[0]))
			}
			ic.ip += 2
		case 5, 6: // Jump-if-true, Jump-if-false
			ic.ensureMemory(ic.ip + 2)
			params := ic.getParams(2)
			val, target := ic.readMemory(params[0]), ic.readMemory(params[1])
			if (opcode == 5 && val != 0) || (opcode == 6 && val == 0) {
				ic.ip = target
			} else {
				ic.ip += 3
			}
		case 99: // Halt
			ic.halted = true
			return
		default:
			panic(fmt.Sprintf("unknown opcode: %d", opcode))
		}
	}
}

func (ic *Intcode) readMemory(address int) int {
	ic.ensureMemory(address)
	return ic.memory[address]
}

func (ic *Intcode) writeMemory(address, value int) {
	ic.ensureMemory(address)
	ic.memory[address] = value
}

func (ic *Intcode) ensureMemory(address int) {
	if address >= len(ic.memory) {
		newMemory := make([]int, address+1)
		copy(newMemory, ic.memory)
		ic.memory = newMemory
	}
}

func (ic *Intcode) Outputs() []int {
	return ic.output
}

func (ic *Intcode) Halted() bool {
	return ic.halted
}

func (ic *Intcode) getParams(count int) []int {
	paramModes := ic.memory[ic.ip] / 100
	params := make([]int, count)
	for i := 0; i < count; i++ {
		if paramModes%10 == 1 {
			params[i] = ic.ip + i + 1
		} else {
			params[i] = ic.memory[ic.ip+i+1]
		}
		paramModes /= 10
	}
	return params
}

func main() {
	// Read and parse Intcode program from input.txt
	data, err := os.ReadFile("day11/input.txt")
	if err != nil {
		panic(err)
	}
	codeStr := strings.Split(strings.TrimSpace(string(data)), ",")
	program := make([]int, len(codeStr))
	for i, s := range codeStr {
		program[i], err = strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
	}

	grid := make(Grid)
	robot := Robot{Position{0, 0}, Up}
	intcode := NewIntcode(program)

	// Run the Intcode program and paint the panels
	for !intcode.Halted() {
		currentColor := grid[robot.Position]
		intcode.AddInput(int(currentColor))
		intcode.Run()
		outputs := intcode.Outputs()

		if len(outputs) == 2 {
			grid[robot.Position] = PanelColor(outputs[0])
			robot.turnAndMove(outputs[1])
		}
	}

	// Count the number of panels painted at least once
	count := len(grid)
	fmt.Println(count)
}
