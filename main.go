package main

import (
	"math/rand"
	"strings"
	"time"
)

type GatesManager struct {
	avaliableStates []StateInput

	currentStates []State
	points        int
}

type State struct {
	Id       string
	Name     string
	Type     string
	Layer    *int
	Cost     int
	Conflict []string
}

func (ths *GatesManager) generateAdditionalState() State {
	blockedSlots := map[string]int{}
	filledSlots := map[string]int{}
	for _, state := range ths.currentStates {
		for _, slot := range state.Conflict {
			blockedSlots[slot] = 0
		}
		filledSlots[state.Type] = 0
	}

	possibleNewStates := []State{}
loop:
	for _, state := range ths.avaliableStates {
		if _, blocked := blockedSlots[state.Type]; blocked {
			continue loop
		}
		for _, slot := range state.Conflict {
			if _, conflict := filledSlots[slot]; conflict {
				continue loop
			}
		}
		if state.Mods != nil {
			mod := state.Mods[rand.Intn(len(state.Mods))]
			modId := state.Id + "_" + mod
			for _, existing := range ths.currentStates {
				if modId == existing.Id {
					continue loop
				}
			}
			possibleNewStates = append(possibleNewStates, State{
				Id:       modId,
				Name:     strings.ReplaceAll(state.Name, "{}", mod),
				Type:     state.Type,
				Layer:    state.Layer,
				Cost:     state.Cost,
				Conflict: state.Conflict,
			})
		} else {
			for _, existing := range ths.currentStates {
				if state.Id == existing.Id {
					continue loop
				}
			}
			possibleNewStates = append(possibleNewStates, State{
				Id:       state.Id,
				Name:     state.Name,
				Type:     state.Type,
				Layer:    state.Layer,
				Cost:     state.Cost,
				Conflict: state.Conflict,
			})
		}
	}

	if len(possibleNewStates) == 0 {
		return State{Name: "Nil"}
	}

	return possibleNewStates[rand.Intn(len(possibleNewStates))]
}

func (ths *GatesManager) GenerateInitialState(lenStates int) {
	for i := 0; i < lenStates; i++ {
		ths.currentStates = append(ths.currentStates, ths.generateAdditionalState())
	}
}

func (ths *GatesManager) RemoveState() State {
	highestLayer := 0
	for _, state := range ths.currentStates {
		if state.Layer != nil && *state.Layer > highestLayer {
			highestLayer = *state.Layer
		}
	}

	possibleRemovals := []int{}
	for i, state := range ths.currentStates {
		if state.Layer != nil && *state.Layer < highestLayer {
			continue
		}
		if state.Cost > ths.points {
			continue
		}
		possibleRemovals = append(possibleRemovals, i)
	}

	if len(possibleRemovals) == 0 {
		return State{
			Name: "Nothing",
		}
	}

	removeIndex := possibleRemovals[rand.Intn(len(possibleRemovals))]
	removed := ths.currentStates[removeIndex]
	ths.points = ths.points - removed.Cost
	ths.currentStates[removeIndex] = ths.currentStates[len(ths.currentStates)-1]
	ths.currentStates = ths.currentStates[:len(ths.currentStates)-1]
	return removed
}

func (ths *GatesManager) PrintCurrentState() {
	for _, state := range ths.currentStates {
		println("You will " + state.Name)
	}
}

func NewGatesManager(input *GatesInput) *GatesManager {
	return &GatesManager{
		avaliableStates: input.States,
		currentStates:   []State{},
	}
}

func main() {
	seed := time.Now().Unix()
	println(seed)
	rand.Seed(seed)
	input, err := NewGatesInput("taskdef.yaml")
	if err != nil {
		panic(err)
	}

	man := NewGatesManager(input)
	man.GenerateInitialState(5)
	man.PrintCurrentState()

	man.points += 1
	println("Stop " + man.RemoveState().Name)
	man.points += 2
	println("Stop " + man.RemoveState().Name)
	man.points += 3
	println("Stop " + man.RemoveState().Name)
	man.points += 4
	println("Stop " + man.RemoveState().Name)
	man.points += 5
	println("Stop " + man.RemoveState().Name)

	// b, _ := yaml.Marshal(man.currentStates)
	// println(string(b))

}
