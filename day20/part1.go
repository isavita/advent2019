package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	Wall = '#'
	Free = '.'
)

type char uint8

func (c char) Letter() bool {
	return c >= 'A' && c <= 'Z'
}

type P struct{ X, Y uint16 }

func (p P) Neighbours() []P {
	return []P{
		{p.X, p.Y + 1},
		{p.X + 1, p.Y},
		{p.X, p.Y - 1},
		{p.X - 1, p.Y},
	}
}

type Map struct {
	XMax       uint16
	YMax       uint16
	Grid       map[P]char
	AA         P
	ZZ         P
	Teleport   map[P]P
	PortalName map[P]string
	IsOuter    map[P]bool
}

func parse() *Map {
	grid := map[P]char{}
	xMax, yMax := uint16(0), uint16(0)

	file, err := os.Open("day20/input.txt")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)

	i := uint16(0)
	for ; scanner.Scan(); i++ {
		line := scanner.Text()

		if uint16(len(line)) > yMax {
			yMax = uint16(len(line))
		}

		for j := uint16(0); j < uint16(len(line)); j++ {
			grid[P{i, j}] = char(line[j])
		}
	}
	xMax = i

	var aa, zz P
	isOuter := map[P]bool{}
	portalName := map[P]string{}
	teleport := map[P]P{}

	cache := map[string]P{}

	for i := uint16(0); i < xMax; i++ {
		for j := uint16(0); j < yMax; j++ {
			c := grid[P{i, j}]

			if !c.Letter() {
				continue
			}

			pName, pPoint, ok := extractPortal(grid, P{i, j})

			if !ok {

				continue
			}

			portalName[pPoint] = pName

			if pName == "AA" {
				aa = pPoint
				isOuter[pPoint] = true
				continue
			}

			if pName == "ZZ" {
				zz = pPoint
				isOuter[pPoint] = true
				continue
			}

			if target, ok := cache[pName]; ok {
				teleport[pPoint] = target
				teleport[target] = pPoint
			} else {
				cache[pName] = pPoint
			}

			switch {
			case j == 0 || i == 0,
				i == xMax-2 || j == yMax-2:

				isOuter[pPoint] = true
			default:

				isOuter[pPoint] = false
			}
		}
	}

	return &Map{
		XMax:       xMax,
		YMax:       yMax,
		Grid:       grid,
		AA:         aa,
		ZZ:         zz,
		Teleport:   teleport,
		PortalName: portalName,
		IsOuter:    isOuter,
	}
}

func extractPortal(grid map[P]char, p P) (string, P, bool) {
	c1 := grid[p]

	if c2 := grid[P{p.X + 1, p.Y}]; c2.Letter() {
		portalName := string(c1) + string(c2)

		portalPoint := P{p.X + 2, p.Y}
		if grid[portalPoint] == '.' {
			return portalName, portalPoint, true
		}

		portalPoint = P{p.X - 1, p.Y}
		if grid[portalPoint] == '.' {
			return portalName, portalPoint, true
		}
	}

	if c2 := grid[P{p.X, p.Y + 1}]; c2.Letter() {
		portalName := string(c1) + string(c2)

		portalPoint := P{p.X, p.Y + 2}
		if grid[portalPoint] == '.' {
			return portalName, portalPoint, true
		}

		portalPoint = P{p.X, p.Y - 1}
		if grid[portalPoint] == '.' {
			return portalName, portalPoint, true
		}
	}

	return "", P{}, false
}
func main() {
	m := parse()

	fmt.Println(BFS(m))
}

func BFS(m *Map) int {
	discovered := map[P]struct{}{}
	toDo := []P{}

	discovered[m.AA] = struct{}{}
	toDo = append(toDo, m.AA)

	depth := 0

	for len(toDo) > 0 {
		for levelSize := len(toDo); levelSize > 0; levelSize-- {
			var curr P
			curr, toDo = toDo[0], toDo[1:]

			if curr == m.ZZ {
				return depth
			}

			for _, n := range curr.Neighbours() {
				dest := m.Grid[n]

				switch {

				case dest == Wall:
					continue

				case dest == Free:

					if _, found := discovered[n]; !found {
						discovered[n] = struct{}{}
						toDo = append(toDo, n)
					}

				case dest.Letter():
					next := m.Teleport[curr]
					if _, found := discovered[next]; !found {
						discovered[next] = struct{}{}
						toDo = append(toDo, next)
					}
				}
			}
		}
		depth++
	}

	return -1
}
