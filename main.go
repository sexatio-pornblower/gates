package main

import (
	"bufio"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/sexatio-pornblower/gates/manager"
)

func initGame() *manager.GatesManager {
	input, err := manager.NewGatesInput("taskdef.yaml")
	if err != nil {
		panic(err)
	}

	man := manager.NewGatesManager(input)
	man.GenerateInitialState(input.Manditory.States, 40)

	// for _, state := range man.CurrentStates {
	// 	println("You are " + state.Name)
	// }

	return man
}

func main() {
	console := bufio.NewScanner(os.Stdin)
	seed := time.Now().Unix()
	println(seed)
	rand.Seed(seed)

	man := initGame()

	running := true
	var target string
	points := 0
	for running {
		println("\n\n\n")
		// Current state
		for _, state := range man.CurrentStates {
			println(state.Id + ":\t" + state.Name + " (" + strconv.Itoa(state.Cost) + ")")
		}
		println("win:\tEnd the game (20)")
		println("You have " + strconv.Itoa(points) + " points")
		println("Enter id to use points buy options or 'c' to continue")
		console.Scan()
		in := console.Text()
		switch in {
		case "win":
			if points >= 20 {
				println("Game Over")
				running = false
			}
		case "c":
			gate := man.NextGate(target)
			target = gate.Target
			println(gate.Description())

			console.Scan()
			print("Cum check: ")
			console.Scan()
			res := console.Text()
			if res == "y" {
				man.OpenGate(gate, 1)
			} else {
				man.OpenGate(gate, 0)
			}
			points += gate.Value
			if gate.Reset {
				target = ""
			}
		default: //buy
			success, state := man.RemoveState(in, points)
			if success {
				points -= state.Cost
				println("Removed " + state.Name)
			} else {
				println("Removal unsuccessful")
			}
		}
	}
}
