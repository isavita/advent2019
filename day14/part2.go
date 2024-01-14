package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type chemical struct {
	name   string
	amount int
}

func parseChemical(s string) chemical {
	parts := strings.Split(s, " ")
	amount, _ := strconv.Atoi(parts[0])
	return chemical{name: parts[1], amount: amount}
}

func main() {
	file, err := os.Open("day14/input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	reactions := make(map[string]chemical)
	ingredients := make(map[string][]chemical)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " => ")
		output := parseChemical(parts[1])
		var inputs []chemical
		for _, in := range strings.Split(parts[0], ", ") {
			inputs = append(inputs, parseChemical(in))
		}
		reactions[output.name] = output
		ingredients[output.name] = inputs
	}

	const oreAvailable int64 = 1000000000000
	fmt.Println(maxFuel(reactions, ingredients, oreAvailable))
}

func calculateOre(chem string, amount int, reactions map[string]chemical, ingredients map[string][]chemical, surplus map[string]int) int {
	if chem == "ORE" {
		return amount
	}

	if surplus[chem] >= amount {
		surplus[chem] -= amount
		return 0
	}

	amount -= surplus[chem]
	surplus[chem] = 0
	reaction := reactions[chem]
	times := (amount + reaction.amount - 1) / reaction.amount
	ore := 0

	for _, ingredient := range ingredients[chem] {
		ore += calculateOre(ingredient.name, ingredient.amount*times, reactions, ingredients, surplus)
	}

	surplus[chem] += times*reaction.amount - amount
	return ore
}

func maxFuel(reactions map[string]chemical, ingredients map[string][]chemical, oreAvailable int64) int64 {
	low, high := int64(0), oreAvailable
	for low < high {
		mid := (low + high + 1) / 2
		if calculateOre("FUEL", int(mid), reactions, ingredients, make(map[string]int)) > int(oreAvailable) {
			high = mid - 1
		} else {
			low = mid
		}
	}
	return low
}
