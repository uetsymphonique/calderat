package objects

import (
	"calderat/objects/secondclass"
	"calderat/utils/logger"

	"github.com/google/uuid"
)

const (
	FINISHED        = 0
	RUNNING         = 1
	WAITING_TO_STOP = 2
)

type Operation struct {
	OperationID string
	Name        string
	Adversary   Adversary
	Abilities   map[string]Ability
	Autonomous  bool
	Links       []secondclass.Link
	Logger      *logger.Logger
	Ignored     []Ability
	Status      int
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
	o.Status = RUNNING
	o.Logger.Log(logger.TRACE, "Running operation %s", o.Name)
	for o.Status == RUNNING {
		if o.Autonomous {
			for _, a := range o.Abilities {
				if a.IsAvailable() {
					o.Logger.Log(logger.DEBUG, "Creating links of ability %s", a.Name)
					o.Links = append(o.Links, a.CreateLinks()...)
					o.RemoveAbility(a.AbilityId)
				}
			}
		} else {
			// TODO: waiting keyboard input

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
		Ignored:     []Ability{},
		Logger:      log,
		Status:      FINISHED,
	}
	operation.AddAbilities(abilities)
	return &operation
}
