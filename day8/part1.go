package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	data, err := os.ReadFile("day8/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	imageData := strings.TrimSpace(string(data))

	width, height := 25, 6
	layerSize := width * height

	minZeros := layerSize + 1
	result := 0

	for i := 0; i < len(imageData); i += layerSize {
		layer := imageData[i:min(i+layerSize, len(imageData))]
		zeroCount, oneCount, twoCount := 0, 0, 0

		for _, pixel := range layer {
			switch pixel {
			case '0':
				zeroCount++
			case '1':
				oneCount++
			case '2':
				twoCount++
			}
		}

		if zeroCount < minZeros {
			minZeros = zeroCount
			result = oneCount * twoCount
		}
	}

	fmt.Println(result)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
