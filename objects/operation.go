package objects

import (
	"calderat/objects/secondclass"
	"calderat/service/execute"
	"calderat/service/knowledge"
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
	Source            Source
	Facts             map[string][]*secondclass.Fact
	Autonomous        bool
	Cleanup           bool
	Links             []secondclass.Link
	CleanupLinks      []secondclass.Link
	Logger            *logger.Logger
	Ignored           []Ability
	Status            int
	shells            []string
	ExecutingServices map[string]execute.ExecutingService
	KnowledgeService  *knowledge.KnowledgeService
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
				o.Logger.Log(logger.INFO, "Running ability %s", ability.Name)
				links, cleanupLinks := ability.CreateLinks(o.Logger, o.shells, o.Facts)
				o.Links = append(o.Links, links...)
				o.Logger.Log(logger.DEBUG, "Creating links of ability %s", ability.Name)
				o.CleanupLinks = append(o.CleanupLinks, cleanupLinks...)
				for _, link := range links {
					link.Execute(o.ExecutingServices[link.Executor.Name])
					o.attireLog.AddLinkResult(&link)
					o.attireLog.DumpToFile("log.json")
				}
			}
		}
	}
	if o.Cleanup {
		o.CleanupOperation()
	}

	o.Logger.Log(logger.INFO, "Operation (%s - %s) successfully executed!", o.Name, o.OperationID)
}

func (o *Operation) CleanupOperation() {
	o.Logger.Log(logger.TRACE, "Cleaning up operation %s", o.Name)
	for i := len(o.CleanupLinks) - 1; i >= 0; i-- {
		link := o.CleanupLinks[i]
		link.Execute(o.ExecutingServices[link.Executor.Name])
		o.attireLog.AddLinkResult(&link)
		o.attireLog.DumpToFile("log.json")
		o.Logger.Log(logger.DEBUG, "Cleaning up link of ability %s", link.Executor.Name)
	}
	o.Logger.Log(logger.INFO, "Operation (%s - %s) cleanup successfully executed!", o.Name, o.OperationID)
}
func NewOperation(adversary Adversary, autonomous, cleanup bool, abilities []Ability, shells []string, os string, ip string, log *logger.Logger, knowledgeService *knowledge.KnowledgeService) *Operation {
	operation := Operation{
		OperationID:       uuid.New().String(),
		Name:              adversary.Name,
		Adversary:         adversary,
		Autonomous:        autonomous,
		Cleanup:           cleanup,
		Abilities:         map[string]Ability{},
		Source:            Source{Logger: log},
		Facts:             map[string][]*secondclass.Fact{},
		Links:             []secondclass.Link{},
		CleanupLinks:      []secondclass.Link{},
		Ignored:           []Ability{},
		Logger:            log,
		Status:            FINISHED,
		shells:            shells,
		os:                os,
		attireLog:         *NewAttireLog(ip),
		ExecutingServices: map[string]execute.ExecutingService{},
		KnowledgeService:  knowledgeService,
	}
	operation.AddAbilities(abilities)
	operation.addingExecutingServices()
	operation.Source.LoadFromYAML("data/source.yml")
	operation.addingFacts()
	return &operation
}

func (o *Operation) addingFacts() {
	for _, fact := range o.Source.Facts {
		if facts, exists := o.Facts[fact.Trait]; exists {
			o.Facts[fact.Trait] = append(facts, &fact)
		} else {
			o.Facts[fact.Trait] = []*secondclass.Fact{&fact}
		}
	}
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
