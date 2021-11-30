package manager

import (
	"os"

	"gopkg.in/yaml.v2"
)

type GatesInput struct {
	States    []StateInput
	Gates     []GateInput
	Manditory ManditoryInput
}

type StateInput struct {
	Id       string
	Name     string
	Type     string
	Layer    *int
	Cost     int
	Conflict []string
	Mods     []string
	Block    []string
	Disabled bool
}

type GateInput struct {
	Id             string
	Name           string
	Type           string
	Value          int
	Intensity      int
	Target         []GateTarget
	Timer          *GateTimer
	Penalty        *int
	Resetintensity bool
	Disabled       bool
}

type GateTarget struct {
	Name  string
	Value int
	Level int
}

type ManditoryInput struct {
	States []string
}

func NewGatesInput(path string) (*GatesInput, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	input := &GatesInput{}
	err = yaml.Unmarshal(data, input)
	if err != nil {
		return nil, err
	}
	return input, nil
}
