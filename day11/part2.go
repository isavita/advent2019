package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Mode int

const (
	Position Mode = iota
	Immediate
	Relative
)

type Cmd int

func (c Cmd) OpCode() int {
	// last 2 digits
	return int(c) % 100
}

func (c Cmd) Modes(arity int) []Mode {
	// remove last two digits
	modeSection := int(c) / 100

	modes := make([]Mode, arity)
	for i := 0; i < arity; i++ {
		modes[i] = Mode(modeSection / int(math.Pow10(i)) % 10)
	}
	return modes
}

type C struct{ X, Y int }

const (
	Up = iota
	Right
	Down
	Left
)

const (
	Black = iota
	White
)

const (
	TurnLeft = iota
	TurnRight
)

type Robot struct {
	Pos       C
	Direction int
}

func NewRobot() *Robot {
	return &Robot{
		Pos:       C{0, 0},
		Direction: Up,
	}
}

func (r *Robot) Turn(where int) {
	switch where {
	case TurnLeft:
		r.Direction = (r.Direction + 3) % 4
	case TurnRight:
		r.Direction = (r.Direction + 1) % 4
	}
}

func (r *Robot) MoveForward() {
	// move forward
	switch r.Direction {
	case Up:
		r.Pos.Y++
	case Right:
		r.Pos.X++
	case Down:
		r.Pos.Y--
	case Left:
		r.Pos.X--
	}
}

func (r *Robot) Paint(grid map[C]int, color int) {
	grid[r.Pos] = color
}

type VM struct {
	Code         map[int]int // source code and memory
	Ip           int         // Instruction Pointer
	Input        <-chan int
	Output       chan<- int
	RelativeBase int
}

func NewVM(filename string, input <-chan int, output chan<- int) *VM {
	vm := &VM{
		Code:   nil,
		Ip:     0,
		Input:  input,
		Output: output,
	}

	vm.Load(filename)

	return vm
}

func (vm *VM) Load(filename string) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	listStr := strings.Split(strings.TrimSpace(string(buf)), ",")
	code := make(map[int]int, len(listStr))

	for i := range listStr {
		var err error
		code[i], err = strconv.Atoi(listStr[i])
		if err != nil {
			log.Fatal(err)
		}
	}
	vm.Code = code
	vm.Ip = 0
	vm.RelativeBase = 0
}

func (vm *VM) Run() {
	var arity int

	for {
		cmd := Cmd(vm.Code[vm.Ip])

		switch cmd.OpCode() {

		case 1:
			// add
			arity = 3
			params := vm.getParamsAddresses(vm.Ip, cmd, arity)
			vm.Code[params[2]] = vm.Code[params[0]] + vm.Code[params[1]]

		case 2:
			// multiply
			arity = 3
			params := vm.getParamsAddresses(vm.Ip, cmd, arity)
			vm.Code[params[2]] = vm.Code[params[0]] * vm.Code[params[1]]

		case 3:
			// read
			arity = 1
			params := vm.getParamsAddresses(vm.Ip, cmd, arity)
			vm.Code[params[0]] = <-vm.Input

		case 4:
			// write
			arity = 1
			params := vm.getParamsAddresses(vm.Ip, cmd, arity)
			vm.Output <- vm.Code[params[0]]

		case 5:
			// jump not zero
			arity = 2
			params := vm.getParamsAddresses(vm.Ip, cmd, arity)
			if vm.Code[params[0]] != 0 {
				vm.Ip = vm.Code[params[1]]
				continue
			}

		case 6:
			// jump zero
			arity = 2
			params := vm.getParamsAddresses(vm.Ip, cmd, arity)
			if vm.Code[params[0]] == 0 {
				vm.Ip = vm.Code[params[1]]
				continue
			}

		case 7:
			// less than
			arity = 3
			params := vm.getParamsAddresses(vm.Ip, cmd, arity)

			if vm.Code[params[0]] < vm.Code[params[1]] {
				vm.Code[params[2]] = 1
			} else {
				vm.Code[params[2]] = 0
			}

		case 8:
			// equal
			arity = 3
			params := vm.getParamsAddresses(vm.Ip, cmd, arity)

			if vm.Code[params[0]] == vm.Code[params[1]] {
				vm.Code[params[2]] = 1
			} else {
				vm.Code[params[2]] = 0
			}
		case 9:
			// change relative base
			arity = 1
			params := vm.getParamsAddresses(vm.Ip, cmd, arity)
			vm.RelativeBase += vm.Code[params[0]]

		case 99:
			// halt
			return

		default:
			log.Fatalf("not an opcode %v", cmd)
		}

		vm.Ip += arity + 1
	}
}

func (vm *VM) getParamsAddresses(pos int, cmd Cmd, arity int) []int {
	modes := cmd.Modes(arity)
	results := make([]int, arity)
	for i := 0; i < arity; i++ {
		results[i] = vm.getParamAddress(pos+i+1, modes[i])
	}
	return results
}

func (vm *VM) getParamAddress(pos int, mode Mode) int {
	switch mode {
	case Position:
		return vm.Code[pos]
	case Immediate:
		return pos
	case Relative:
		return vm.RelativeBase + vm.Code[pos]
	}
	log.Fatal("wrong mode")
	return -1
}

func run(grid map[C]int) {
	// 1 because we should be able to write input
	// although the program terminated and can't read it
	camera := make(chan int, 1)
	output := make(chan int, 0)

	go func() {
		NewVM("day11/input.txt", camera, output).Run()
		close(output)
	}()

	robot := NewRobot()

	for {
		// send color from camera
		camera <- grid[robot.Pos]

		// receive color
		newColor, open := <-output
		exit(!open, grid)

		// receive rotation
		turn := <-output

		robot.Paint(grid, newColor)
		robot.Turn(turn)
		robot.MoveForward()
	}
}

func main() {
	grid := make(map[C]int)

	// first position must be white
	grid[C{0, 0}] = White

	run(grid)
}

// if output channel is closed,
// program stopped to work
func exit(closed bool, grid map[C]int) {
	if closed {
		print(grid)
		os.Exit(0)
	}
}

func print(grid map[C]int) {
	// image is mirrored, so print upside down
	for j := 1; j > -6; j-- {

		// left to right
		for i := 0; i < 41; i++ {

			switch grid[C{i, j}] {
			case Black:
				fmt.Print(" ")
			case White:
				fmt.Print("█")
			}

		}
		fmt.Println()
	}
}
