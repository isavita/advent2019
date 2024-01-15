package main

import (
	"fmt"
	"image"
	"os"
	"strconv"
	"strings"
)

type mode int

const (
	position mode = iota
	immediate
	relative
)

type opcode int

const (
	add opcode = 1 + iota
	mul
	input
	output
	jt
	jf
	lt
	eq
	rbo
	halt opcode = 99
)

func decode(n int) (op opcode, modes [3]mode) {
	op = opcode(n % 100)
	n /= 100
	for i := 0; i < 3; i++ {
		modes[i] = mode(n % 10)
		n /= 10
	}
	return
}

type Machine struct {
	data    map[int]int
	ip      int
	in      <-chan int
	out     chan<- int
	relbase int
}

func New(program []int, in <-chan int, out chan<- int) *Machine {
	m := &Machine{
		data: make(map[int]int),
		in:   in,
		out:  out,
	}
	for i, n := range program {
		m.data[i] = n
	}
	return m
}

func (m *Machine) get(i int, mo mode) int {
	switch mo {
	case immediate:
		return m.data[i]
	case position:
		return m.data[m.data[i]]
	case relative:
		return m.data[m.relbase+m.data[i]]
	default:
		panic(fmt.Sprintf("Unknown mode: %v", mo))
	}
	return 0
}

func (m *Machine) set(i int, mo mode, val int) {
	switch mo {
	case position:
		m.data[m.data[i]] = val
	case relative:
		m.data[m.relbase+m.data[i]] = val
	default:
		panic(fmt.Sprintf("Unknown mode: %v", mo))
	}
}

func (m *Machine) Step() bool {
	op, modes := decode(m.data[m.ip])
	switch op {
	case add:
		val := m.get(m.ip+1, modes[0]) + m.get(m.ip+2, modes[1])
		m.set(m.ip+3, modes[2], val)
		m.ip += 4
	case mul:
		val := m.get(m.ip+1, modes[0]) * m.get(m.ip+2, modes[1])
		m.set(m.ip+3, modes[2], val)
		m.ip += 4
	case input:
		m.set(m.ip+1, modes[0], <-m.in)
		m.ip += 2
	case output:
		m.out <- m.get(m.ip+1, modes[0])
		m.ip += 2
	case jt:
		if m.get(m.ip+1, modes[0]) != 0 {
			m.ip = m.get(m.ip+2, modes[1])
		} else {
			m.ip += 3
		}
	case jf:
		if m.get(m.ip+1, modes[0]) == 0 {
			m.ip = m.get(m.ip+2, modes[1])
		} else {
			m.ip += 3
		}
	case lt:
		if m.get(m.ip+1, modes[0]) < m.get(m.ip+2, modes[1]) {
			m.set(m.ip+3, modes[2], 1)
		} else {
			m.set(m.ip+3, modes[2], 0)
		}
		m.ip += 4
	case eq:
		if m.get(m.ip+1, modes[0]) == m.get(m.ip+2, modes[1]) {
			m.set(m.ip+3, modes[2], 1)
		} else {
			m.set(m.ip+3, modes[2], 0)
		}
		m.ip += 4
	case rbo:
		m.relbase += m.get(m.ip+1, modes[0])
		m.ip += 2
	case halt:
		return false
	default:
		panic(fmt.Sprintf("Unknown opcode: %v", op))
	}
	return true
}

func (m *Machine) Run() {
	for m.Step() {
	}
	close(m.out)
}

func Run(program []int, in <-chan int) chan int {
	out := make(chan int)
	m := New(program, in, out)
	go m.Run()
	return out
}

type tile int

const (
	empty tile = iota
	wall
	block
	paddle
	ball
)

type move int

const (
	left move = iota - 1
	none
	right
)

func main() {
	var program []int
	input := readAll("day13/input.txt")
	for _, n := range strings.Split(input, ",") {
		program = append(program, toInt(n))
	}

	fmt.Println(countBlocks(program))
}

func countBlocks(program []int) int {
	grid := make(map[image.Point]tile)
	in := make(chan int)
	close(in)
	out := Run(program, in)
	for {
		x, ok := <-out
		if !ok {
			break
		}
		y := <-out
		grid[image.Pt(x, y)] = tile(<-out)
	}
	var n int
	for _, t := range grid {
		if t == block {
			n++
		}
	}
	return n
}

func readAll(filepath string) string {
	file, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(file))
}

func toInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}
