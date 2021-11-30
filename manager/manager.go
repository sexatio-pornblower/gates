package manager

type GatesManager struct {
	CurrentStates []State

	avaliableStates  []StateInput
	avaliableGates   []GateInput
	lastGateId       string
	currentIntensity int
	currentLevel     int
}

func NewGatesManager(input *GatesInput) *GatesManager {
	return &GatesManager{
		avaliableStates: input.States,
		avaliableGates:  input.Gates,
		CurrentStates:   []State{},
	}
}
