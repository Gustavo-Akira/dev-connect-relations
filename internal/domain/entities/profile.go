package entities

import "errors"

type Profile struct {
	ConnectId int32
	Name      string
}

func NewProfile(id int32, name string) (*Profile, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	if id == 0 {
		return nil, errors.New("id cannot be 0")
	}
	return &Profile{ConnectId: id, Name: name}, nil
}
