package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	// Read and parse the input file
	file, err := os.Open("day16/input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	input := scanner.Text()

	// Prepare the repeated input signal
	repeatedInput := repeatInput(input, 10000)

	// Find the message offset
	offset, err := strconv.Atoi(input[:7])
	if err != nil {
		panic(err)
	}

	// Apply FFT for 100 phases starting from the offset
	for phase := 0; phase < 100; phase++ {
		sum := 0
		for i := len(repeatedInput) - 1; i >= offset; i-- {
			sum += repeatedInput[i]
			repeatedInput[i] = sum % 10
		}
	}

	// Print the eight-digit message
	for i := offset; i < offset+8; i++ {
		fmt.Print(repeatedInput[i])
	}
	fmt.Println()
}

func repeatInput(input string, times int) []int {
	digits := make([]int, len(input)*times)
	for t := 0; t < times; t++ {
		for i, r := range input {
			digit, err := strconv.Atoi(string(r))
			if err != nil {
				panic(err)
			}
			digits[t*len(input)+i] = digit
		}
	}
	return digits
}
