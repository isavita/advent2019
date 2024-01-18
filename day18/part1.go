package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	x, y int
}

type State struct {
	pos  Point
	keys int
}

func main() {
	file, _ := os.Open("day18/input.txt")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var grid []string
	var start Point
	keyMap := make(map[byte]int)
	keyCounter := 0

	for y := 0; scanner.Scan(); y++ {
		line := scanner.Text()
		grid = append(grid, line)
		for x := range line {
			if line[x] == '@' {
				start = Point{x, y}
			} else if line[x] >= 'a' && line[x] <= 'z' {
				keyMap[line[x]] = keyCounter
				keyCounter++
			}
		}
	}

	fmt.Println(findShortestPath(grid, start, keyMap))
}

func findShortestPath(grid []string, start Point, keyMap map[byte]int) int {
	dirs := []Point{{0, -1}, {-1, 0}, {0, 1}, {1, 0}}
	visited := make(map[State]bool)
	queue := []State{{start, 0}}
	steps := 0

	for len(queue) > 0 {
		size := len(queue)
		for i := 0; i < size; i++ {
			current := queue[0]
			queue = queue[1:]

			if current.keys == (1<<len(keyMap))-1 {
				return steps
			}

			for _, d := range dirs {
				next := Point{current.pos.x + d.x, current.pos.y + d.y}
				if next.x >= 0 && next.x < len(grid[0]) && next.y >= 0 && next.y < len(grid) {
					char := grid[next.y][next.x]
					if char != '#' && !(char >= 'A' && char <= 'Z' && current.keys&(1<<keyMap[char+'a'-'A']) == 0) {
						newState := State{next, current.keys}
						if char >= 'a' && char <= 'z' {
							newState.keys |= 1 << keyMap[char]
						}
						if !visited[newState] {
							visited[newState] = true
							queue = append(queue, newState)
						}
					}
				}
			}
		}
		steps++
	}
	return -1
}
