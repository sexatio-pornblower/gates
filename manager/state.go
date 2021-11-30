package manager

import (
	"math/rand"
	"strings"
)

type State struct {
	Id       string
	Name     string
	Type     string
	Layer    *int
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
					Layer:    state.Layer,
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
					Layer:    state.Layer,
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
				Layer:    state.Layer,
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
				Layer:    state.Layer,
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

func (ths *GatesManager) RemoveState(id string, points int) (bool, State) {
	highestLayer := 0
	for _, state := range ths.CurrentStates {
		if state.Layer != nil && *state.Layer > highestLayer {
			highestLayer = *state.Layer
		}
	}

	removeIndex := -1
	for i, state := range ths.CurrentStates {
		if state.Id == id {
			if state.Layer != nil && *state.Layer < highestLayer {
				return false, state
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
	highestLayer := 0
	for _, state := range ths.CurrentStates {
		if state.Layer != nil && *state.Layer > highestLayer {
			highestLayer = *state.Layer
		}
	}

	possibleRemovals := []int{}
	for i, state := range ths.CurrentStates {
		if state.Layer != nil && *state.Layer < highestLayer {
			continue
		}
		if state.Cost > points {
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
	removed := ths.CurrentStates[removeIndex]
	ths.CurrentStates[removeIndex] = ths.CurrentStates[len(ths.CurrentStates)-1]
	ths.CurrentStates = ths.CurrentStates[:len(ths.CurrentStates)-1]
	return removed
}
