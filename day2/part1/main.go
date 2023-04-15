package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Open input file
	f, err := os.Open("day2/input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	// Read input data
	var inputData []int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		vals := strings.Split(text, ",")
		for _, val := range vals {
			num, _ := strconv.Atoi(val)
			inputData = append(inputData, num)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}

	// Replace positions 1 and 2
	inputData[1] = 12
	inputData[2] = 2

	// Execute program
	result := executeProgram(inputData)

	// Print result
	fmt.Println(result)
}

// Function to execute the program
func executeProgram(data []int) int {
	// Iterate over each element in the array
	for i := 0; i < len(data)-3; i += 4 {
		pos1 := data[i+1]
		pos2 := data[i+2]
		pos3 := data[i+3]
		switch data[i] {
		case 1:
			sum := data[pos1] + data[pos2]
			data[pos3] = sum
		case 2:
			product := data[pos1] * data[pos2]
			data[pos3] = product
		case 99:
			return data[0]
		default:
			panic("Invalid opcode")
		}
	}

	return data[0]
}
