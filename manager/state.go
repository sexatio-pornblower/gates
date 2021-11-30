package manager

import (
	"math/rand"
	"strings"
)

type State struct {
	Id       string
	Name     string
	Type     string
	Stack    map[string]int
	Cost     int
	Conflict []string
	Block    []string
}

func (ths *GatesManager) addState(id string) int {
	for _, state := range ths.avaliableStates {
		if state.Id == id {
			if state.Mods != nil {
				mod := state.Mods[rand.Intn(len(state.Mods))]
				modId := state.Id + "_" + mod
				for _, existing := range ths.CurrentStates {
					if modId == existing.Id {
						continue
					}
				}
				ths.CurrentStates = append(ths.CurrentStates, State{
					Id:       modId,
					Name:     strings.ReplaceAll(state.Name, "{}", mod),
					Type:     state.Type,
					Stack:    state.Stack,
					Cost:     state.Cost,
					Conflict: state.Conflict,
					Block:    state.Block,
				})
				return state.Cost
			} else {
				for _, existing := range ths.CurrentStates {
					if state.Id == existing.Id {
						continue
					}
				}
				ths.CurrentStates = append(ths.CurrentStates, State{
					Id:       state.Id,
					Name:     state.Name,
					Type:     state.Type,
					Stack:    state.Stack,
					Cost:     state.Cost,
					Conflict: state.Conflict,
					Block:    state.Block,
				})
				return state.Cost
			}
		}
	}
	return 0
}

func (ths *GatesManager) generateAdditionalState() State {
	blockedSlots := map[string]int{}
	filledSlots := map[string]int{}
	for _, state := range ths.CurrentStates {
		for _, slot := range state.Conflict {
			blockedSlots[slot] = 0
		}
		filledSlots[state.Type] = 0
	}

	possibleNewStates := []State{}
loop:
	for _, state := range ths.avaliableStates {
		if state.Disabled {
			continue loop
		}
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
			for _, existing := range ths.CurrentStates {
				if modId == existing.Id {
					continue loop
				}
			}
			possibleNewStates = append(possibleNewStates, State{
				Id:       modId,
				Name:     strings.ReplaceAll(state.Name, "{}", mod),
				Type:     state.Type,
				Stack:    state.Stack,
				Cost:     state.Cost,
				Conflict: state.Conflict,
				Block:    state.Block,
			})
		} else {
			for _, existing := range ths.CurrentStates {
				if state.Id == existing.Id {
					continue loop
				}
			}
			possibleNewStates = append(possibleNewStates, State{
				Id:       state.Id,
				Name:     state.Name,
				Type:     state.Type,
				Stack:    state.Stack,
				Cost:     state.Cost,
				Conflict: state.Conflict,
				Block:    state.Block,
			})
		}
	}

	if len(possibleNewStates) == 0 {
		return State{Name: "Nil"}
	}

	return possibleNewStates[rand.Intn(len(possibleNewStates))]
}

func (ths *GatesManager) GenerateInitialState(manditoryStates []string, pointMin int) {
	p := 0
	for _, manditoryState := range manditoryStates {
		p += ths.addState(manditoryState)
	}

	for p < pointMin {
		state := ths.generateAdditionalState()
		p += state.Cost
		ths.CurrentStates = append(ths.CurrentStates, state)
	}
}

func (ths *GatesManager) AddRandomState() State {
	state := ths.generateAdditionalState()
	ths.CurrentStates = append(ths.CurrentStates, state)
	return state
}

func (ths *GatesManager) RemoveState(id string, points int) (bool, State) {
	highestStacks := make(map[string]int)
	for _, state := range ths.CurrentStates {
		for stack, height := range state.Stack {
			if _, ok := highestStacks[stack]; ok && highestStacks[stack] >= height {
				continue
			}
			highestStacks[stack] = height
		}
	}

	removeIndex := -1
	for i, state := range ths.CurrentStates {
		if state.Id == id {
			for stack, height := range state.Stack {
				if height < highestStacks[stack] {
					return false, state
				}
			}
			if state.Cost > points {
				return false, state
			}
			removeIndex = i
			break
		}
	}

	if removeIndex == -1 {
		return false, State{}
	}

	removed := ths.CurrentStates[removeIndex]
	ths.CurrentStates[removeIndex] = ths.CurrentStates[len(ths.CurrentStates)-1]
	ths.CurrentStates = ths.CurrentStates[:len(ths.CurrentStates)-1]
	return true, removed
}

func (ths *GatesManager) RemoveRandomState(points int) State {
	highestStacks := make(map[string]int)
	for _, state := range ths.CurrentStates {
		for stack, height := range state.Stack {
			if _, ok := highestStacks[stack]; ok && highestStacks[stack] >= height {
				continue
			}
			highestStacks[stack] = height
		}
	}

	possibleRemovals := []int{}
loop:
	for i, state := range ths.CurrentStates {
		for stack, height := range state.Stack {
			if height < highestStacks[stack] {
				continue loop
			}
		}
		if state.Cost > points {
			continue loop
		}
		possibleRemovals = append(possibleRemovals, i)
	}

	if len(possibleRemovals) == 0 {
		return State{
			Name: "Nothing",
		}
	}

	removeIndex := possibleRemovals[rand.Intn(len(possibleRemovals))]
	removed := ths.CurrentStates[removeIndex]
	ths.CurrentStates[removeIndex] = ths.CurrentStates[len(ths.CurrentStates)-1]
	ths.CurrentStates = ths.CurrentStates[:len(ths.CurrentStates)-1]
	return removed
}
