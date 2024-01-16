package main

import (
	"container/heap"
	"fmt"
	"image"
	"math"
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

func (d Dir) Point() image.Point {
	return point[d]
}

var Neighbors4 = []image.Point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

type Pt3 struct {
	X, Y, Z int
}

func Manhattan(p, q image.Point) int {
	return Abs(p.X-q.X) + Abs(p.Y-q.Y)
}

type Item struct {
	Obj      interface{}
	Priority int
}

type PQ []*Item

func (pq PQ) Len() int            { return len(pq) }
func (pq PQ) Less(i, j int) bool  { return pq[i].Priority > pq[j].Priority }
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(*Item)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[:n-1]
	return item
}

func (p1 Pt3) Add(p2 Pt3) Pt3 {
	return Pt3{p1.X + p2.X, p1.Y + p2.Y, p1.Z + p2.Z}
}

func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

type pathfinder struct {
	m       *Machine
	grid    map[image.Point]byte
	in, out chan int
	dirmap  map[Dir]int
	p       image.Point
	oxygen  image.Point
}

func newPathfinder(program []int) *pathfinder {
	pf := &pathfinder{
		grid:   make(map[image.Point]byte),
		in:     make(chan int),
		out:    make(chan int),
		dirmap: map[Dir]int{N: 1, S: 2, W: 3, E: 4},
	}
	pf.grid[pf.p] = '.'
	pf.m = New(program, pf.in, pf.out)
	go pf.m.Run()
	return pf
}

func (pf *pathfinder) tryMove(dir Dir) bool {
	pf.in <- pf.dirmap[dir]
	next := pf.p.Add(dir.Point())
	switch <-pf.out {
	case 0:
		pf.grid[next] = '#'
		return false
	case 1:
		pf.grid[next] = '.'
	case 2:
		pf.grid[next] = 'O'
		pf.oxygen = next
	}
	pf.p = next
	return true
}

func (pf *pathfinder) explore() {
	for len(pf.open()) > 0 {
		if _, ok := pf.open()[pf.p]; !ok {
			min := math.MaxInt32
			var next image.Point
			for to := range pf.open() {
				dist := Manhattan(pf.p, to)
				if dist < min {
					min, next = dist, to
				}
			}
			minpath := pf.shortestPath(pf.p, next)
			for _, m := range minpath {
				if !pf.tryMove(m) {
					panic("bad path")
				}
			}
		}
		for {
			var d Dir
			for _, n := range Neighbors4 {
				if _, ok := pf.grid[pf.p.Add(n)]; !ok {
					d = DirFromPoint(n)
					break
				}
			}
			if !pf.tryMove(d) {
				break
			}
		}
	}
}

var fromPoint = map[image.Point]Dir{{0, 1}: N, {1, 0}: E, {0, -1}: S, {-1, 0}: W}

func DirFromPoint(p image.Point) Dir {
	return fromPoint[p]
}

func (pf *pathfinder) open() map[image.Point]struct{} {
	ps := make(map[image.Point]struct{})
	for p, b := range pf.grid {
		if b == '#' {
			continue
		}
		for _, n := range Neighbors4 {
			if _, ok := pf.grid[p.Add(n)]; !ok {
				ps[p] = struct{}{}
				break
			}
		}
	}
	return ps
}

func (pf *pathfinder) shortestPath(from, to image.Point) []Dir {
	pq := PQ{&Item{from, 0}}
	pathmap := map[image.Point][]Dir{from: {}}
	heap.Init(&pq)
	for pq.Len() > 0 {
		curr := heap.Pop(&pq).(*Item).Obj.(image.Point)
		currpath := pathmap[curr]
		if curr == to {
			break
		}
		for _, n := range Neighbors4 {
			next := curr.Add(n)
			if b, ok := pf.grid[next]; !ok || b == '#' {
				continue
			}
			if path, ok := pathmap[next]; !ok || len(path) > 1+len(currpath) {
				nextpath := make([]Dir, 1+len(currpath))
				copy(nextpath, currpath)
				nextpath[len(nextpath)-1] = DirFromPoint(n)
				heap.Push(&pq, &Item{next, -len(nextpath)})
				pathmap[next] = nextpath
			}
		}
	}
	path, ok := pathmap[to]
	if !ok {
		panic("no path")
	}
	return path
}

func (pf *pathfinder) longestPath(from image.Point) int {
	pq := PQ{&Item{from, 0}}
	distmap := map[image.Point]int{from: 0}
	heap.Init(&pq)
	for pq.Len() > 0 {
		curr := heap.Pop(&pq).(*Item).Obj.(image.Point)
		currdist := distmap[curr]
		for _, n := range Neighbors4 {
			next := curr.Add(n)
			if b, ok := pf.grid[next]; !ok || b == '#' {
				continue
			}
			if dist, ok := distmap[next]; !ok || dist > 1+currdist {
				nextdist := 1 + currdist
				heap.Push(&pq, &Item{next, -nextdist})
				distmap[next] = nextdist
			}
		}
	}
	var max int
	for _, d := range distmap {
		max = Max(max, d)
	}
	return max
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

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	var program []int
	for _, n := range strings.Split(readAll("day15/input.txt"), ",") {
		program = append(program, Atoi(n))
	}

	pf := newPathfinder(program)
	pf.explore()
	fmt.Println(len(pf.shortestPath(image.Pt(0, 0), pf.oxygen)))
}
