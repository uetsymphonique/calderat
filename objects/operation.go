package objects

import (
	"calderat/objects/secondclass"

	"github.com/google/uuid"
)

type Operation struct {
	operationID string
	name        string
	adversary   Adversary
	autonomous  bool
	links       []secondclass.Link
}

func NewOperation(adversary Adversary, autonomous bool) *Operation {
	operation := Operation{
		operationID: uuid.New().String(),
		name:        adversary.Name(),
		adversary:   adversary,
		autonomous:  autonomous,
		links:       []secondclass.Link{},
	}
	return &operation
}
