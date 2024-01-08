package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	X, Y int
}

func main() {
	data, _ := os.ReadFile("day3/input.txt")
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	wire1 := getPointsWithSteps(lines[0])
	wire2 := getPointsWithSteps(lines[1])

	minSteps := int(^uint(0) >> 1)
	for p, steps1 := range wire1 {
		if steps2, found := wire2[p]; found {
			totalSteps := steps1 + steps2
			if totalSteps < minSteps {
				minSteps = totalSteps
			}
		}
	}

	fmt.Println(minSteps)
}

func getPointsWithSteps(path string) map[Point]int {
	points := make(map[Point]int)
	current := Point{0, 0}
	steps := 0
	for _, move := range strings.Split(path, ",") {
		dir := move[0]
		dist, _ := strconv.Atoi(move[1:])
		for i := 0; i < dist; i++ {
			steps++
			switch dir {
			case 'U':
				current.Y++
			case 'D':
				current.Y--
			case 'L':
				current.X--
			case 'R':
				current.X++
			}
			if _, found := points[current]; !found {
				points[current] = steps
			}
		}
	}
	return points
}
