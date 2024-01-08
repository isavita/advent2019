package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func isValidPassword(password int) bool {
	s := strconv.Itoa(password)
	hasDouble := false

	for i := 0; i < len(s)-1; i++ {
		if s[i] > s[i+1] {
			return false
		}
		if s[i] == s[i+1] {
			// Check if the pair is part of a larger group
			if (i == 0 || s[i] != s[i-1]) && (i+2 >= len(s) || s[i] != s[i+2]) {
				hasDouble = true
			}
		}
	}

	return hasDouble
}

func main() {
	data, err := os.ReadFile("day4/input.txt")
	if err != nil {
		panic(err)
	}

	rangeStr := strings.TrimSpace(string(data))
	ranges := strings.Split(rangeStr, "-")
	start, _ := strconv.Atoi(ranges[0])
	end, _ := strconv.Atoi(ranges[1])

	count := 0
	for i := start; i <= end; i++ {
		if isValidPassword(i) {
			count++
		}
	}

	fmt.Println(count)
}
