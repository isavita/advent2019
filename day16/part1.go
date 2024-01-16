package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	// Read input from file
	file, err := os.Open("day16/input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	input := scanner.Text()

	// Convert input string to slice of ints
	digits := make([]int, len(input))
	for i, r := range input {
		digits[i], err = strconv.Atoi(string(r))
		if err != nil {
			panic(err)
		}
	}

	// Apply FFT algorithm for 100 phases
	for phase := 0; phase < 100; phase++ {
		digits = applyFFT(digits)
	}

	// Output the first eight digits
	for i := 0; i < 8; i++ {
		fmt.Print(digits[i])
	}
	fmt.Println()
}

func applyFFT(input []int) []int {
	basePattern := []int{0, 1, 0, -1}
	output := make([]int, len(input))
	for i := range input {
		sum := 0
		for j, val := range input {
			patternValue := basePattern[((j+1)/(i+1))%len(basePattern)]
			sum += val * patternValue
		}
		output[i] = abs(sum % 10)
	}
	return output
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
