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

func processLine(line string) int {
	n, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		fmt.Println("Error parsing line")
		return 0
	}

	return n
}

func getTotal(masses []int) int {
	var total int

	for i := range masses {
		total += calcFuelMass(masses[i])
	}

	return int(total)
}

func calcFuelMass(mass int) int {
	fuel := (math.Floor(float64(mass)/3) - 2)
	if fuel <= 0 {
		return 0
	}

	return int(fuel) + calcFuelMass(int(fuel))
}

func main() {
	file, _ := os.Open("day1/input.txt")
	defer file.Close()
	reader := bufio.NewReader(file)

	var masses []int
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error while reading line")
			continue
		}
		n := processLine(string(line))
		masses = append(masses, n)
	}
	total := getTotal(masses)

	fmt.Println("Total Fuel Requirement: ", total)
}
