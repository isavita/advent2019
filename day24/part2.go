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

type Space map[int]*[Square]bool

func main() {

	input := parse()

	space := Space{
		0: &input,
	}

	for i := 0; i < 200; i++ {
		space = next2(space)
	}

	count := 0
	for _, grid := range space {
		for i := 0; i < Square; i++ {
			if grid[i] {
				count++
			}
		}
	}
	fmt.Println(count)
}

func next2(space Space) Space {
	newSpace := Space{}

	minLevel, maxLevel := minMaxLevel(space)

	for level := minLevel - 1; level <= maxLevel+1; level++ {
		newSpace[level] = &[Square]bool{}

		for cell := 0; cell < Square; cell++ {

			if cell == 12 {
				continue
			}

			row, col := cell/Side, cell%Side
			neighbours := 0

			if row == 0 {
				if infested(space, level-1, 7) {
					neighbours++
				}
			}

			if col == 0 {
				if infested(space, level-1, 11) {
					neighbours++
				}
			}

			if col == 4 {
				if infested(space, level-1, 13) {
					neighbours++
				}
			}

			if row == 4 {
				if infested(space, level-1, 17) {
					neighbours++
				}
			}

			if cell == 7 {
				for i := 0; i < Side; i++ {
					if infested(space, level+1, i) {
						neighbours++
					}
				}
			}

			if cell == 11 {
				for i := 0; i < Side; i++ {
					if infested(space, level+1, 5*i) {
						neighbours++
					}
				}
			}

			if cell == 13 {
				for i := 0; i < Side; i++ {
					if infested(space, level+1, 5*i+Side-1) {
						neighbours++
					}
				}
			}

			if cell == 17 {
				for i := 0; i < Side; i++ {
					if infested(space, level+1, (Side-1)*Side+i) {
						neighbours++
					}
				}
			}

			if row > 0 && cell != 17 {
				if infested(space, level, cell-Side) {
					neighbours++
				}
			}

			if col > 0 && cell != 13 {
				if infested(space, level, cell-1) {
					neighbours++
				}
			}

			if col < Side-1 && cell != 11 {
				if infested(space, level, cell+1) {
					neighbours++
				}
			}

			if row < Side-1 && cell != 7 {
				if infested(space, level, cell+Side) {
					neighbours++
				}
			}

			if infested(space, level, cell) && neighbours != 1 {
				newSpace[level][cell] = false
				continue
			}

			if !infested(space, level, cell) && (neighbours == 1 || neighbours == 2) {
				newSpace[level][cell] = true
				continue
			}

			newSpace[level][cell] = infested(space, level, cell)
		}
	}

	clean(newSpace)

	return newSpace
}

func clean(space Space) {
	min, max := minMaxLevel(space)

	countMin, countMax := 0, 0
	for cell := 0; cell < Square; cell++ {
		if space[min][cell] {
			countMin++
		}
		if space[max][cell] {
			countMax++
		}
	}
	if countMin == 0 {
		delete(space, min)
	}
	if countMax == 0 {
		delete(space, max)
	}
}

func infested(space Space, level, cell int) bool {
	if space[level] == nil {
		return false
	}
	return space[level][cell]
}

func minMaxLevel(space Space) (int, int) {
	min, max := +999999, -999999
	for level := range space {
		if level < min {
			min = level
		}
		if level > max {
			max = level
		}
	}
	return min, max
}
