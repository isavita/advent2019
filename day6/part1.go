package main

import (
	"fmt"
	"os"
	"strings"
)

func countOrbits(orbitMap map[string][]string, start string, depth int) int {
	orbits, exists := orbitMap[start]
	if !exists {
		return depth
	}
	count := depth
	for _, orbit := range orbits {
		count += countOrbits(orbitMap, orbit, depth+1)
	}
	return count
}

func main() {
	data, err := os.ReadFile("day6/input.txt")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	orbitMap := make(map[string][]string)
	for _, line := range lines {
		parts := strings.Split(line, ")")
		center, orbiter := parts[0], parts[1]
		orbitMap[center] = append(orbitMap[center], orbiter)
	}

	totalOrbits := countOrbits(orbitMap, "COM", 0)
	fmt.Println(totalOrbits)
}
