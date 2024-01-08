package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	data, _ := os.ReadFile("day2/input.txt")
	strs := strings.Split(strings.TrimSpace(string(data)), ",")
	original := make([]int, len(strs))
	for i, s := range strs {
		original[i], _ = strconv.Atoi(s)
	}

	for noun := 0; noun <= 99; noun++ {
		for verb := 0; verb <= 99; verb++ {
			memory := make([]int, len(original))
			copy(memory, original)
			memory[1] = noun
			memory[2] = verb
			if execute(memory) == 19690720 {
				fmt.Println(100*noun + verb)
				return
			}
		}
	}
}

func execute(memory []int) int {
	for i := 0; i < len(memory); i += 4 {
		switch memory[i] {
		case 1:
			memory[memory[i+3]] = memory[memory[i+1]] + memory[memory[i+2]]
		case 2:
			memory[memory[i+3]] = memory[memory[i+1]] * memory[memory[i+2]]
		case 99:
			return memory[0]
		}
	}
	return memory[0]
}
