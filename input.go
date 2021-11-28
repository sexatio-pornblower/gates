package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type GatesInput struct {
	States []StateInput
}

type StateInput struct {
	Id       string
	Name     string
	Type     string
	Layer    *int
	Cost     int
	Conflict []string
	Mods     []string
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
