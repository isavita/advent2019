package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	f, err := os.Open("day2/input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

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

	inputData[1] = 12
	inputData[2] = 2

	result := executeProgram(inputData)

	fmt.Println(result)
}

func executeProgram(data []int) int {
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
