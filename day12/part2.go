package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

type Vec3 struct {
	x, y, z int
}

type Moon struct {
	pos, vel Vec3
}

func applyGravity(moons []Moon, axis string) {
	for i := 0; i < len(moons); i++ {
		for j := i + 1; j < len(moons); j++ {
			switch axis {
			case "x":
				if moons[i].pos.x > moons[j].pos.x {
					moons[i].vel.x--
					moons[j].vel.x++
				} else if moons[i].pos.x < moons[j].pos.x {
					moons[i].vel.x++
					moons[j].vel.x--
				}
			case "y":
				if moons[i].pos.y > moons[j].pos.y {
					moons[i].vel.y--
					moons[j].vel.y++
				} else if moons[i].pos.y < moons[j].pos.y {
					moons[i].vel.y++
					moons[j].vel.y--
				}
			case "z":
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
}

func applyVelocity(moons []Moon, axis string) {
	for i := range moons {
		switch axis {
		case "x":
			moons[i].pos.x += moons[i].vel.x
		case "y":
			moons[i].pos.y += moons[i].vel.y
		case "z":
			moons[i].pos.z += moons[i].vel.z
		}
	}
}

func findCycle(moons []Moon, initialMoons []Moon, axis string) int {
	for steps := 1; ; steps++ {
		applyGravity(moons, axis)
		applyVelocity(moons, axis)

		match := true
		for i, m := range moons {
			switch axis {
			case "x":
				if m.pos.x != initialMoons[i].pos.x || m.vel.x != initialMoons[i].vel.x {
					match = false
				}
			case "y":
				if m.pos.y != initialMoons[i].pos.y || m.vel.y != initialMoons[i].vel.y {
					match = false
				}
			case "z":
				if m.pos.z != initialMoons[i].pos.z || m.vel.z != initialMoons[i].vel.z {
					match = false
				}
			}
		}

		if match {
			return steps
		}
	}
}

func lcm(a, b int) *big.Int {
	bigA := big.NewInt(int64(a))
	bigB := big.NewInt(int64(b))
	return new(big.Int).Div(new(big.Int).Mul(bigA, bigB), new(big.Int).GCD(nil, nil, bigA, bigB))
}

func main() {
	file, err := os.Open("day12/input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var moons, initialMoons []Moon
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var x, y, z int
		fmt.Sscanf(scanner.Text(), "<x=%d, y=%d, z=%d>", &x, &y, &z)
		moons = append(moons, Moon{Vec3{x, y, z}, Vec3{0, 0, 0}})
		initialMoons = append(initialMoons, Moon{Vec3{x, y, z}, Vec3{0, 0, 0}})
	}

	cycleX := findCycle(moons, initialMoons, "x")
	cycleY := findCycle(moons, initialMoons, "y")
	cycleZ := findCycle(moons, initialMoons, "z")

	lcmXY := lcm(cycleX, cycleY)
	lcmXYZ := lcm(int(lcmXY.Int64()), cycleZ)

	fmt.Println(lcmXYZ)
}
