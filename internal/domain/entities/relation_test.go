package entities

import "testing"

func TestShouldRelationBeCreatedWhenAllParametersAreValid(t *testing.T) {
	var id int32 = 1
	var toId int32 = 2
	var relationType RelationType = RelationFriend
	relation, err := NewRelation(id, toId, relationType)
	if err != nil || relation == nil {
		t.Error("Should not throw error in these cases")
	}
}

func TestShouldThrowErrorWhenToIdIsEqualFromID(t *testing.T) {
	var id int32 = 1
	var toId int32 = 1
	var relationType RelationType = RelationFriend
	relation, err := NewRelation(id, toId, relationType)
	if err.Error() != "cannot create relation with self" || relation != nil {
		t.Error("Should throw self relation impossible")
	}
}

func TestShouldThrowErrorWhenRelationTypeIsNone(t *testing.T) {
	var id int32 = 2
	var toId int32 = 1
	var relationType RelationType
	relation, err := NewRelation(id, toId, relationType)
	if err.Error() != "relation type is required" || relation != nil {
		t.Error("Should throw relation type is required")
	}
}
