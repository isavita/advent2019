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
	finalImage := make([]rune, layerSize)

	for i := range finalImage {
		finalImage[i] = '2'
	}

	for i := 0; i < len(imageData); i += layerSize {
		layer := imageData[i:min(i+layerSize, len(imageData))]

		for j, pixel := range layer {
			if finalImage[j] == '2' {
				finalImage[j] = pixel
			}
		}
	}

	fmt.Println("Decoded image:")
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			pixel := finalImage[i*width+j]
			if pixel == '0' {
				fmt.Print(" ")
			} else if pixel == '1' {
				fmt.Print("#")
			}
		}
		fmt.Println()
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
