package objects

import (
	"calderat/execute"
	"calderat/objects/secondclass"
	"calderat/utils/logger"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

const (
	FINISHED        = 0
	RUNNING         = 1
	WAITING_TO_STOP = 2
)

type Operation struct {
	OperationID       string
	Name              string
	Adversary         Adversary
	Abilities         map[string]Ability
	Autonomous        bool
	Links             []secondclass.Link
	Logger            *logger.Logger
	Ignored           []Ability
	Status            int
	shells            []string
	ExecutingServices map[string]execute.ExecutingService
	os                string
	attireLog         AttireLog
	ip                string
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
	for _, ability_id := range o.Adversary.AtomicOrdering {
		if ability, exists := o.Abilities[ability_id]; exists {
			if ability.IsAvailable(o.shells) {
				links := ability.CreateLinks(o.Logger)
				o.Links = append(o.Links, links...)
				o.Logger.Log(logger.DEBUG, "Creating links of ability %s", ability.Name)
				for _, link := range links {
					link.Execute(o.ExecutingServices[link.Executor.Name])
				}
			}
		}
	}

	o.Logger.Log(logger.DEBUG, "Operation (%s - %s) successfully executed!", o.Name, o.OperationID)
}

func NewOperation(adversary Adversary, autonomous bool, abilities []Ability, shells []string, os string, ip string, log *logger.Logger) *Operation {
	operation := Operation{
		OperationID:       uuid.New().String(),
		Name:              adversary.Name,
		Adversary:         adversary,
		Autonomous:        autonomous,
		Abilities:         map[string]Ability{},
		Links:             []secondclass.Link{},
		Ignored:           []Ability{},
		Logger:            log,
		Status:            FINISHED,
		shells:            shells,
		os:                os,
		attireLog:         *NewAttireLog(ip),
		ExecutingServices: map[string]execute.ExecutingService{},
	}
	operation.AddAbilities(abilities)
	operation.addingExecutingServices()
	return &operation
}

func (o *Operation) addingExecutingServices() {
	if o.os == "windows" {
		if slices.Contains(o.shells, "psh") {
			o.ExecutingServices["psh"] = execute.NewPowerShell(o.Logger)
		}
		if slices.Contains(o.shells, "cmd") {
			o.ExecutingServices["cmd"] = execute.NewCmd(o.Logger)
		}
	} else {
		if slices.Contains(o.shells, "sh") {
			o.ExecutingServices["sh"] = execute.NewSh(o.Logger)
		}
	}
}
