package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type asteroid struct {
	x, y  int
	angle float64
	dist  float64
}

func main() {
	asteroids := readAsteroids("day10/input.txt")
	station, _ := findBestAsteroidLocation(asteroids)
	vaporized := vaporizeAsteroids(asteroids, station)
	if len(vaporized) >= 200 {
		result := vaporized[199].x*100 + vaporized[199].y
		fmt.Println(result)
	} else {
		fmt.Println("Less than 200 asteroids were vaporized.")
	}
}

func readAsteroids(filename string) (asteroids [][]bool) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		asteroidRow := make([]bool, len(line))
		for i, char := range line {
			asteroidRow[i] = char == '#'
		}
		asteroids = append(asteroids, asteroidRow)
	}
	return asteroids
}
func vaporizeAsteroids(asteroids [][]bool, station [2]int) []asteroid {
	var targets []asteroid
	for y, row := range asteroids {
		for x, isAsteroid := range row {
			if isAsteroid && !(x == station[0] && y == station[1]) {
				angle := math.Atan2(float64(y-station[1]), float64(x-station[0]))
				dist := math.Hypot(float64(x-station[0]), float64(y-station[1]))
				if angle < -math.Pi/2 {
					angle += 2 * math.Pi // Adjust angle for clockwise rotation
				}
				targets = append(targets, asteroid{x, y, angle, dist})
			}
		}
	}

	sort.Slice(targets, func(i, j int) bool {
		if targets[i].angle == targets[j].angle {
			return targets[i].dist < targets[j].dist
		}
		return targets[i].angle < targets[j].angle
	})

	var vaporized []asteroid
	for len(targets) > 0 {
		var lastAngle float64 = -math.MaxFloat64
		for i := 0; i < len(targets); {
			if targets[i].angle != lastAngle {
				vaporized = append(vaporized, targets[i])
				lastAngle = targets[i].angle
				targets = append(targets[:i], targets[i+1:]...)
			} else {
				i++
			}
		}
	}
	return vaporized
}

func findBestAsteroidLocation(asteroids [][]bool) (bestLocation [2]int, maxCount int) {
	for y, row := range asteroids {
		for x, isAsteroid := range row {
			if isAsteroid {
				count := countVisibleAsteroids(asteroids, x, y)
				if count > maxCount {
					maxCount = count
					bestLocation = [2]int{x, y}
				}
			}
		}
	}
	return bestLocation, maxCount
}

func countVisibleAsteroids(asteroids [][]bool, x, y int) int {
	angles := make(map[float64]bool)
	for otherY, row := range asteroids {
		for otherX, isAsteroid := range row {
			if isAsteroid && !(otherX == x && otherY == y) {
				angle := math.Atan2(float64(otherY-y), float64(otherX-x))
				angles[angle] = true
			}
		}
	}
	return len(angles)
}
