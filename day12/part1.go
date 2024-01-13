package main

import (
	"bufio"
	"fmt"
	"os"
)

type Vec3 struct {
	x, y, z int
}

type Moon struct {
	pos, vel Vec3
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func applyGravity(moons []Moon) {
	for i := 0; i < len(moons); i++ {
		for j := i + 1; j < len(moons); j++ {
			if moons[i].pos.x > moons[j].pos.x {
				moons[i].vel.x--
				moons[j].vel.x++
			} else if moons[i].pos.x < moons[j].pos.x {
				moons[i].vel.x++
				moons[j].vel.x--
			}

			if moons[i].pos.y > moons[j].pos.y {
				moons[i].vel.y--
				moons[j].vel.y++
			} else if moons[i].pos.y < moons[j].pos.y {
				moons[i].vel.y++
				moons[j].vel.y--
			}

			if moons[i].pos.z > moons[j].pos.z {
				moons[i].vel.z--
				moons[j].vel.z++
			} else if moons[i].pos.z < moons[j].pos.z {
				moons[i].vel.z++
				moons[j].vel.z--
			}
		}
	}
}

func applyVelocity(moons []Moon) {
	for i := range moons {
		moons[i].pos.x += moons[i].vel.x
		moons[i].pos.y += moons[i].vel.y
		moons[i].pos.z += moons[i].vel.z
	}
}

func totalEnergy(moons []Moon) int {
	total := 0
	for _, m := range moons {
		pot := Abs(m.pos.x) + Abs(m.pos.y) + Abs(m.pos.z)
		kin := Abs(m.vel.x) + Abs(m.vel.y) + Abs(m.vel.z)
		total += pot * kin
	}
	return total
}

func main() {
	file, err := os.Open("day12/input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var moons []Moon
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var x, y, z int
		fmt.Sscanf(scanner.Text(), "<x=%d, y=%d, z=%d>", &x, &y, &z)
		moons = append(moons, Moon{Vec3{x, y, z}, Vec3{0, 0, 0}})
	}

	for step := 0; step < 1000; step++ {
		applyGravity(moons)
		applyVelocity(moons)
	}

	fmt.Println(totalEnergy(moons))
}
