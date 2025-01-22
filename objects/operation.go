package objects

import (
	"calderat/objects/secondclass"

	"github.com/google/uuid"
)

type Operation struct {
	OperationID string
	Name        string
	Adversary   Adversary
	Abilities   map[string]Ability
	Autonomous  bool
	Links       []secondclass.Link
}

func (o *Operation) AddAbility(ability Ability) {
	o.Abilities[ability.AbilityId] = ability
}

func (o *Operation) RemoveAbility(ability_id string) {
	delete(o.Abilities, ability_id)
}

func (o *Operation) AddAbilities(abilities []Ability) {
	for _, a := range abilities {
		o.AddAbility(a)
	}
}

func NewOperation(adversary Adversary, autonomous bool, abilities []Ability) *Operation {
	operation := Operation{
		OperationID: uuid.New().String(),
		Name:        adversary.Name,
		Adversary:   adversary,
		Autonomous:  autonomous,
		Links:       []secondclass.Link{},
	}
	operation.AddAbilities(abilities)
	return &operation
}
