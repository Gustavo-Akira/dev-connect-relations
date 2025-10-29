package entities

import "fmt"

type Stack struct {
	Name string
}

func NewStack(name string) (*Stack, error) {
	if name == "" {
		return nil, fmt.Errorf("stack name cannot be empty")
	}
	return &Stack{Name: name}, nil
}
