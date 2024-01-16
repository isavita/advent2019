package main

import (
	"bufio"
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

type Dir int

const (
	N Dir = iota
	E
	S
	W
)

var point = map[Dir]image.Point{N: {0, 1}, E: {1, 0}, S: {0, -1}, W: {-1, 0}}
var pointReversed = map[Dir]image.Point{N: {0, -1}, E: {1, 0}, S: {0, 1}, W: {-1, 0}}

// Y-axis goes up.
func (d Dir) Point() image.Point {
	return point[d]
}

// Y-axis goes down.
func (d Dir) PointR() image.Point {
	return pointReversed[d]
}

var Neighbors4 = []image.Point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

type Pt3 struct {
	X, Y, Z int
}

type Item struct {
	Obj      interface{}
	Priority int
}

func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

var fromPoint = map[image.Point]Dir{{0, 1}: N, {1, 0}: E, {0, -1}: S, {-1, 0}: W}

func DirFromPoint(p image.Point) Dir {
	return fromPoint[p]
}

func (d Dir) Next() Dir {
	return (d + 1) % 4
}

func (d Dir) Prev() Dir {
	return (d + 3) % 4
}

var fromByte = map[byte]Dir{
	'N': N,
	'E': E,
	'S': S,
	'W': W,
	'U': N,
	'R': E,
	'D': S,
	'L': W,
	'^': N,
	'>': E,
	'v': S,
	'<': W,
}

func DirFromByte(b byte) Dir {
	return fromByte[b]
}

func readAll(filepath string) string {
	input, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(input))
}

func Atoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

func main() {
	var program []int
	for _, n := range strings.Split(readAll("day17/input.txt"), ",") {
		program = append(program, Atoi(n))
	}
	scaffolding, _, _ := parse(program)

	fmt.Println(sumAlign(scaffolding))
}

func parse(program []int) (map[image.Point]struct{}, image.Point, Dir) {
	in := make(chan int)
	close(in)
	out := Run(program, in)
	var sb strings.Builder
	for o := range out {
		fmt.Fprintf(&sb, "%c", o)
	}
	scaffolding := make(map[image.Point]struct{})
	var robot image.Point
	var dir Dir
	s := bufio.NewScanner(strings.NewReader(sb.String()))
	var y int
	for s.Scan() {
		line := s.Text()
		for x := 0; x < len(line); x++ {
			switch line[x] {
			case '^', 'v', '<', '>':
				robot = image.Pt(x, y)
				dir = DirFromByte(line[x])
				fallthrough
			case '#':
				scaffolding[image.Pt(x, y)] = struct{}{}
			}
		}
		y++
	}
	return scaffolding, robot, dir
}

func sumAlign(grid map[image.Point]struct{}) int {
	var sum int
	for p := range grid {
		x := true
		for _, n := range Neighbors4 {
			if _, ok := grid[p.Add(n)]; !ok {
				x = false
			}
		}
		if x {
			sum += p.X * p.Y
		}
	}
	return sum
}

func path(scaffolding map[image.Point]struct{}, robot image.Point, dir Dir) string {
	var dist int
	var d byte
	var sections []string
	for {
		if _, ok := scaffolding[robot.Add(dir.PointR())]; ok {
			robot = robot.Add(dir.PointR())
			dist++
			continue
		}
		if dist > 0 {
			sections = append(sections, fmt.Sprintf("%c,%d", d, dist))
		}
		if _, ok := scaffolding[robot.Add(dir.Next().PointR())]; ok {
			robot = robot.Add(dir.Next().PointR())
			dir = dir.Next()
			dist = 1
			d = 'R'
		} else if _, ok := scaffolding[robot.Add(dir.Prev().PointR())]; ok {
			robot = robot.Add(dir.Prev().PointR())
			dir = dir.Prev()
			dist = 1
			d = 'L'
		} else {
			break
		}
	}
	return strings.Join(sections, ",")
}

func encode(path string) (seq, a, b, c string) {
loop:
	for i := 2; i <= 21; i++ {
		for j := 2; j <= 21; j++ {
			for k := 2; k <= 21; k++ {
				next := path + ","
				a = next[:i]
				next = strings.ReplaceAll(next, a, "")
				b = next[:j]
				next = strings.ReplaceAll(next, b, "")
				c = next[:k]
				next = strings.ReplaceAll(next, c, "")
				if next == "" {
					break loop
				}
			}
		}
	}
	a, b, c = strings.Trim(a, ","), strings.Trim(b, ","), strings.Trim(c, ",")
	path = strings.ReplaceAll(path, a, "A")
	path = strings.ReplaceAll(path, b, "B")
	path = strings.ReplaceAll(path, c, "C")
	path = strings.Trim(path, ",")
	return path, a, b, c
}
