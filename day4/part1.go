package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("day4/input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	rangeStr := scanner.Text()
	parts := strings.Split(rangeStr, "-")
	start, _ := strconv.Atoi(parts[0])
	end, _ := strconv.Atoi(parts[1])

	count := 0
	for i := start; i <= end; i++ {
		s := strconv.Itoa(i)
		if hasDoubleAndIncreasingDigits(s) {
			count++
		}
	}

	fmt.Println(count)
}

func hasDoubleAndIncreasingDigits(s string) bool {
	hasDouble := false
	for i := 0; i < len(s)-1; i++ {
		if s[i] == s[i+1] {
			hasDouble = true
		}
		if s[i] > s[i+1] {
			return false
		}
	}
	return hasDouble
}
