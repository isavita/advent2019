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

var masses []int  // An array to store our masses
var total float64 // Our calculated total mass

func processLine(line string) { // A helper func to parse one single line from stdin
	var err error
	m, err := strconv.Atoi(strings.TrimSpace(line)) // trim whitespace and convert string to int
	if err != nil {
		fmt.Println("Error parsing line")
		return
	}

	masses = append(masses, m) // Add the current mass to the list
}

func getTotal() { // Calcuate & return total fuel requirement
	var tempTotal float64 // temp variable to hold intermediate calcualtion

	for i := range masses { // Loop over our stored masses
		tempTotal += (math.Floor(float64(masses[i])/3) - 2) // Calculate fuel for single item
	}

	total = tempTotal
	return
}

func main() {
	// Open & read data from our source
	file, _ := os.Open("day1/input.txt")
	defer file.Close() // Always remember to close files!
	reader := bufio.NewReader(file)

	// Loop reading in lines from our file into our processor func
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF { // Exit when we reach end of file
			break
		} else if err != nil {
			fmt.Println("Error while reading line")
			continue
		}

		processLine(string(line))
	}
	getTotal()

	// Finalise calculations & print results
	fmt.Println("Total Fuel Requirement: ", total)
}
