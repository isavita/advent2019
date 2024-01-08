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
	wire1 := getPoints(lines[0])
	wire2 := getPoints(lines[1])

	intersections := make(map[Point]bool)
	for p := range wire1 {
		if wire2[p] {
			intersections[p] = true
		}
	}

	minDistance := int(^uint(0) >> 1)
	for p := range intersections {
		distance := abs(p.X) + abs(p.Y)
		if distance < minDistance {
			minDistance = distance
		}
	}

	fmt.Println(minDistance)
}

func getPoints(path string) map[Point]bool {
	points := make(map[Point]bool)
	current := Point{0, 0}
	for _, move := range strings.Split(path, ",") {
		dir := move[0]
		steps, _ := strconv.Atoi(move[1:])
		for i := 0; i < steps; i++ {
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
			points[Point{current.X, current.Y}] = true
		}
	}
	return points
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
