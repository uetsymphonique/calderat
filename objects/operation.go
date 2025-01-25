package objects

import (
	"calderat/objects/secondclass"
	"calderat/utils/logger"

	"github.com/google/uuid"
)

type Operation struct {
	OperationID string
	Name        string
	Adversary   Adversary
	Abilities   map[string]Ability
	Autonomous  bool
	Links       []secondclass.Link
	Logger      *logger.Logger
}

func (o *Operation) AddAbility(ability Ability) {
	o.Logger.Log(logger.TRACE, "Add ability %s", ability.Name)
	o.Abilities[ability.AbilityId] = ability
	o.Logger.Log(logger.DEBUG, "Added ability (%s - %s) successfully!", ability.Name, ability.AbilityId)
}

func (o *Operation) RemoveAbility(ability_id string) {
	delete(o.Abilities, ability_id)
}

func (o *Operation) AddAbilities(abilities []Ability) {
	for _, a := range abilities {
		o.AddAbility(a)
	}
}

func (o *Operation) Run() {
	o.Logger.Log(logger.TRACE, "Running operation %s", o.Name)
	for _, a := range o.Abilities {
		if a.IsAvailable() {
			o.Logger.Log(logger.DEBUG, "Creating links of ability %s", a.Name)
			o.Links = append(o.Links, a.CreateLinks()...)
		}
	}
	o.Logger.Log(logger.DEBUG, "Operation (%s - %s) successfully executed!", o.Name, o.OperationID)
}

func NewOperation(adversary Adversary, autonomous bool, abilities []Ability, log *logger.Logger) *Operation {
	operation := Operation{
		OperationID: uuid.New().String(),
		Name:        adversary.Name,
		Adversary:   adversary,
		Autonomous:  autonomous,
		Abilities:   map[string]Ability{},
		Links:       []secondclass.Link{},
		Logger:      log,
	}
	operation.AddAbilities(abilities)
	return &operation
}
