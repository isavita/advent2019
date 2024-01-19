package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	Side   = 5
	Square = Side * Side
)

func parse() [Square]bool {
	res := [Square]bool{}

	file, err := os.Open("day24/input.txt")
	if err != nil {
		panic(err)
	}
	s := bufio.NewScanner(file)
	for row := 0; s.Scan(); row++ {
		line := s.Text()
		for col := 0; col < Side; col++ {
			if line[col] == '#' {
				res[row*Side+col] = true
			} else {
				res[row*Side+col] = false
			}
		}
	}
	return res
}

func main() {
	appeared := map[[Square]bool]struct{}{}

	grid := parse()
	appeared[grid] = struct{}{}
	for {
		grid = next1(grid)
		_, ok := appeared[grid]

		if ok {
			fmt.Println(biodiversity(grid))
			return
		}
		appeared[grid] = struct{}{}
	}
}

func next1(grid [Square]bool) [Square]bool {
	newGrid := [Square]bool{}

	for i := 0; i < Square; i++ {
		row, col := i/Side, i%Side
		neighbours := 0

		if row > 0 {
			if grid[i-Side] {
				neighbours++
			}
		}
		if row < Side-1 {
			if grid[i+Side] {
				neighbours++
			}
		}
		if col > 0 {
			if grid[i-1] {
				neighbours++
			}
		}
		if col < Side-1 {
			if grid[i+1] {
				neighbours++
			}
		}

		if grid[i] && neighbours != 1 {
			newGrid[i] = false
			continue
		}

		if !grid[i] && (neighbours == 1 || neighbours == 2) {
			newGrid[i] = true
			continue
		}

		newGrid[i] = grid[i]
	}

	return newGrid
}

func biodiversity(grid [Square]bool) int {
	bio := 0
	for i := 0; i < Square; i++ {
		if grid[i] {
			bio += 1 << i
		}
	}
	return bio
}
