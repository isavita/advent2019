package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	asteroids := readAsteroids("day10/input.txt")
	maxCount := findBestAsteroidLocation(asteroids)
	fmt.Println(maxCount)
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

func findBestAsteroidLocation(asteroids [][]bool) int {
	maxCount := 0
	for y, row := range asteroids {
		for x, isAsteroid := range row {
			if isAsteroid {
				count := countVisibleAsteroids(asteroids, x, y)
				if count > maxCount {
					maxCount = count
				}
			}
		}
	}
	return maxCount
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
