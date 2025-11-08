package entities

import (
	"errors"
	"strings"
)

type StackRelation struct {
	StackName string
	ProfileID int64
}

func NewStackRelation(stackName string, profileID int64) (*StackRelation, error) {
	if stackName == "" || profileID == 0 {
		return nil, errors.New("invalid parameters for StackRelation")
	}
	stackName = strings.ToLower(stackName)
	return &StackRelation{
		StackName: stackName,
		ProfileID: profileID,
	}, nil
}
