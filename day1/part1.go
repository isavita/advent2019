package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

var masses []int
var total float64

func processLine(line string) {
	var err error
	m, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		fmt.Println("Error parsing line")
		return
	}

	masses = append(masses, m)
}

func getTotal() {
	var tempTotal float64

	for i := range masses {
		tempTotal += (math.Floor(float64(masses[i])/3) - 2)
	}

	total = tempTotal
	return
}

func main() {
	file, _ := os.Open("day1/input.txt")
	defer file.Close()
	reader := bufio.NewReader(file)

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error while reading line")
			continue
		}

		processLine(string(line))
	}
	getTotal()

	fmt.Println(total)
}
