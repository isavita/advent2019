package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
)

type Cmd int

func (c Cmd) OpCode() int {
	// last 2 digits
	return int(c) % 100
}

func (c Cmd) Immediate(paramNum int) bool {
	// 10^2 =   100  for 1
	// 10^3 =  1000  for 2
	// 10^4 = 10000  for 3
	digit := int(math.Pow10(paramNum + 1))
	return int(c)/digit%10 == 1
}

type VM struct {
	Code   []int // source code and memory, list of int
	Ip     int   // Instruction Pointer
	Input  <-chan int
	Output chan<- int
}

func (v *VM) Load(filename string) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	listStr := strings.Split(strings.TrimSpace(string(buf)), ",")
	listInt := make([]int, len(listStr))

	for i := range listStr {
		var err error
		listInt[i], err = strconv.Atoi(listStr[i])
		if err != nil {
			log.Fatal(err)
		}
	}
	v.Code = listInt
	v.Ip = 0
}

func (v *VM) Run() {
	for {
		cmd := Cmd(v.Code[v.Ip])

		switch cmd.OpCode() {

		case 1:
			// add
			param1 := v.getParam(v.Ip+1, cmd.Immediate(1))
			param2 := v.getParam(v.Ip+2, cmd.Immediate(2))
			address := v.getParam(v.Ip+3, true)
			v.Code[address] = param1 + param2
			v.Ip += 4

		case 2:
			// multiply
			param1 := v.getParam(v.Ip+1, cmd.Immediate(1))
			param2 := v.getParam(v.Ip+2, cmd.Immediate(2))
			address := v.getParam(v.Ip+3, true)
			v.Code[address] = param1 * param2
			v.Ip += 4

		case 3:
			// read
			address := v.getParam(v.Ip+1, true)
			v.Code[address] = <-v.Input
			v.Ip += 2

		case 4:
			// write
			param1 := v.getParam(v.Ip+1, cmd.Immediate(1))
			v.Output <- param1
			v.Ip += 2

		case 5:
			// jump not zero
			param1 := v.getParam(v.Ip+1, cmd.Immediate(1))
			param2 := v.getParam(v.Ip+2, cmd.Immediate(2))
			if param1 != 0 {
				v.Ip = param2
			} else {
				v.Ip += 3
			}

		case 6:
			// jump zero
			param1 := v.getParam(v.Ip+1, cmd.Immediate(1))
			param2 := v.getParam(v.Ip+2, cmd.Immediate(2))
			if param1 == 0 {
				v.Ip = param2
			} else {
				v.Ip += 3
			}

		case 7:
			// less than
			param1 := v.getParam(v.Ip+1, cmd.Immediate(1))
			param2 := v.getParam(v.Ip+2, cmd.Immediate(2))
			address := v.getParam(v.Ip+3, true)
			if param1 < param2 {
				v.Code[address] = 1
			} else {
				v.Code[address] = 0
			}
			v.Ip += 4

		case 8:
			// equal
			param1 := v.getParam(v.Ip+1, cmd.Immediate(1))
			param2 := v.getParam(v.Ip+2, cmd.Immediate(2))
			address := v.getParam(v.Ip+3, true)
			if param1 == param2 {
				v.Code[address] = 1
			} else {
				v.Code[address] = 0
			}
			v.Ip += 4

		case 99:
			// halt
			return

		default:
			log.Fatalf("not an opcode %v", cmd)
		}
	}
}

func (v *VM) getParam(address int, immediate bool) int {
	param := v.Code[address]
	if immediate {
		return param
	}
	return v.Code[param]
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

func permutations(arr []int) [][]int {
	var helper func([]int, int)
	res := [][]int{}

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}

func main() {
	max := 0
	for _, phase := range permutations([]int{5, 6, 7, 8, 9}) {
		res := runLoop(phase)
		if res > max {
			max = res
		}
	}
	fmt.Println(max)
}

func runLoop(phase []int) int {
	// ch0       ch1       ch2       ch3        ch4       ch0 (loop)
	// -->  VM0  -->  VM1  -->  VM2  -->   VM3  -->  VM4  -->

	var chs = make([]chan int, 6)

	for i := 0; i < 6; i++ {
		chs[i] = make(chan int, 1)
	}

	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		chs[i] <- phase[i]
		go func(i int) {
			NewVM("day7/input.txt", chs[i], chs[(i+1)%5]).Run()
			wg.Done()
		}(i)
	}

	chs[0] <- 0
	wg.Wait()
	return <-chs[0]
}
