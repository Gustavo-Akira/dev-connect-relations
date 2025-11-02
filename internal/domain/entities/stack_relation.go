package entities

import "errors"

type StackRelation struct {
	StackName string
	ProfileID int64
}

func NewStackRelation(stackName string, profileID int64) (*StackRelation, error) {
	if stackName == "" || profileID == 0 {
		return nil, errors.New("invalid parameters for StackRelation")
	}
	return &StackRelation{
		StackName: stackName,
		ProfileID: profileID,
	}, nil
}
