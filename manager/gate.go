package manager

import (
	"math/rand"
	"strconv"
	"strings"
)

const BLANK = "blank"

type Gate struct {
	Id      string
	Text    string
	Value   int
	Type    string
	Target  string
	Penalty *int
	Timer   *GateTimer
	Reset   bool
	Level   int
}

type GateTimer struct {
	Max int
	Min int
}

func (ths *Gate) Description() string {
	if ths.Type == BLANK {
		return "No activity found"
	}
	if ths.Timer != nil {
		time := rand.Intn(ths.Timer.Max-ths.Timer.Min) + ths.Timer.Min
		return "You will be " + strings.ReplaceAll(ths.Text, "[]", strconv.Itoa(time))
	}
	return "You will be " + ths.Text
}

func (ths *GatesManager) expandTargets(gate GateInput, target string) []Gate {
	gates := []Gate{}
	if gate.Target != nil {
		for _, gateTarget := range gate.Target {
			if gate.Id+"_"+gateTarget.Name == ths.lastGateId {
				continue
			}
			if target != "" && gateTarget.Name != target {
				continue
			}
			if gateTarget.Level > ths.currentLevel+1 {
				continue
			}
			gates = append(gates, Gate{
				Id:      gate.Id + "_" + gateTarget.Name,
				Value:   gate.Value + gateTarget.Value,
				Type:    gate.Type,
				Target:  gateTarget.Name,
				Penalty: gate.Penalty,
				Text:    strings.ReplaceAll(gate.Name, "{}", gateTarget.Name),
				Timer:   gate.Timer,
				Reset:   gate.Resetintensity,
				Level:   gateTarget.Level,
			})
		}
	} else {
		if gate.Id != ths.lastGateId {
			gates = append(gates, Gate{
				Id:      gate.Id,
				Value:   gate.Value,
				Type:    gate.Type,
				Target:  target,
				Penalty: gate.Penalty,
				Text:    gate.Name,
				Timer:   gate.Timer,
				Reset:   gate.Resetintensity,
			})
		}
	}
	return gates
}

func (ths *GatesManager) NextGate(target string) Gate {
	blockedGateTypes := map[string]int{}
	for _, state := range ths.CurrentStates {
		if state.Block != nil {
			for _, blockedType := range state.Block {
				blockedGateTypes[blockedType] = 0
			}
		}
	}

	gatesWithGreaterIntensity := []GateInput{}
	for _, gate := range ths.avaliableGates {
		if _, blocked := blockedGateTypes[gate.Type]; blocked {
			continue
		}
		if gate.Disabled {
			continue
		}
		if gate.Intensity >= ths.currentIntensity {
			gatesWithGreaterIntensity = append(gatesWithGreaterIntensity, gate)
		}
	}

	// reset intensity once list is finished
	if len(gatesWithGreaterIntensity) == 0 {
		ths.currentIntensity = 0
		return ths.NextGate(target)
	}

	possibleGates := []Gate{}
	if len(ths.CurrentStates) == 0 {
		for _, gate := range gatesWithGreaterIntensity {
			if gate.Type != "finish" {
				continue
			}
			possibleGates = append(possibleGates, ths.expandTargets(gate, target)...)
		}
	} else {
		for _, gate := range gatesWithGreaterIntensity {
			if gate.Intensity != ths.currentIntensity {
				continue
			}
			possibleGates = append(possibleGates, ths.expandTargets(gate, target)...)
		}
	}

	if len(possibleGates) == 0 {
		return Gate{
			Id:     BLANK,
			Text:   "No activity found",
			Value:  1,
			Type:   BLANK,
			Target: target,
			Reset:  true,
		}
	}

	return possibleGates[rand.Intn(len(possibleGates))]
}

func (ths *GatesManager) OpenGate(gate Gate, penalty int) {
	if gate.Penalty != nil {
		for i := 0; i < *gate.Penalty; i++ {
			ths.CurrentStates = append(ths.CurrentStates, ths.generateAdditionalState())
		}
	}
	for i := 0; i < penalty; i++ {
		ths.CurrentStates = append(ths.CurrentStates, ths.generateAdditionalState())
	}
	ths.lastGateId = gate.Id
	ths.currentIntensity++
	if gate.Reset {
		ths.currentIntensity = 0
	}
	if ths.currentLevel < gate.Level {
		ths.currentLevel = gate.Level
	}
}
