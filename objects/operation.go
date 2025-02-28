package objects

import (
	"calderat/secondclass"
	"calderat/service/execute"
	"calderat/service/knowledge"
	"calderat/utils/colorprint"
	"calderat/utils/logger"
	"fmt"

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
	fmt.Println(colorprint.ColorString("\n------------------------ EXPLOIT PHASE ------------------------", colorprint.YELLOW))
	for index, ability_id := range o.Adversary.AtomicOrdering {
		if ability, exists := o.Abilities[ability_id]; exists {
			if ability.IsAvailable(o.shells) {
				fmt.Println(colorprint.ColorString(fmt.Sprintf("\n[+] Running ability (%d/%d) %s", index, len(o.Adversary.AtomicOrdering), ability.Name), colorprint.YELLOW))
				fmt.Println(colorprint.ColorString(fmt.Sprintf("    [-] %s: %s(%s)", ability.Tactic, ability.Technique, ability.TechniqueId), colorprint.YELLOW))
				links, cleanupLinks := ability.CreateLinks(o.Logger, o.shells, o.Facts)
				o.Links = append(o.Links, links...)
				o.Logger.Log(logger.DEBUG, "Creating links of ability %s", ability.Name)
				o.CleanupLinks = append(o.CleanupLinks, cleanupLinks...)
				for _, link := range links {
					link.Execute(o.ExecutingServices[link.Executor.Name])
					o.attireLog.AddLinkResult(&link)
					o.attireLog.DumpToFile("log.json")
					if !o.Cleanup {
						secondclass.DumpLinksToJson(o.CleanupLinks, "cleanups.json", o.Logger)
					}
				}
			}
		}
	}
	o.Logger.Log(logger.INFO, "Operation (%s - %s) successfully executed!", o.Name, o.OperationID)
	if o.Cleanup {
		fmt.Println(colorprint.ColorString("\n------------------------ CLEANUP PHASE ------------------------", colorprint.YELLOW))
		o.CleanupOperation()
	}

}

func (o *Operation) CleanupOperation() {
	o.Logger.Log(logger.TRACE, "Cleaning up operation %s", o.Name)
	for i := len(o.CleanupLinks) - 1; i >= 0; i-- {
		link := o.CleanupLinks[i]
		o.Logger.Log(logger.INFO, "Cleaning up link of ability %s(%s)", link.ProcedureName, link.MitreTechniqueId)
		link.Execute(o.ExecutingServices[link.Executor.Name])
		o.attireLog.AddLinkResult(&link)
		o.attireLog.DumpToFile("log.json")
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

func NewCleanupOperation(cleanupLinks []secondclass.Link, shells []string, os, ip string, log *logger.Logger) *Operation {
	operation := Operation{
		OperationID:       uuid.New().String(),
		Name:              "Cleanup Operation",
		CleanupLinks:      cleanupLinks,
		Ignored:           []Ability{},
		Logger:            log,
		shells:            shells,
		os:                os,
		attireLog:         *NewAttireLog(ip),
		ExecutingServices: map[string]execute.ExecutingService{},
	}
	operation.addingExecutingServices()
	return &operation
}

func (o *Operation) RunningCleanupOperation() {
	o.Logger.Log(logger.TRACE, "Running cleanup operation")
	for i := len(o.CleanupLinks) - 1; i >= 0; i-- {
		link := o.CleanupLinks[i]
		o.Logger.Log(logger.INFO, "Running cleanup link of ability %s(%s)", link.ProcedureName, link.MitreTechniqueId)
		link.Execute(o.ExecutingServices[link.Executor.Name])
		o.attireLog.AddLinkResult(&link)
		o.attireLog.DumpToFile("cleanup_log.json")

		o.CleanupLinks = append(o.CleanupLinks[:i], o.CleanupLinks[i+1:]...)
		secondclass.DumpLinksToJson(o.CleanupLinks, "not_completed_cleanups.json", o.Logger)
	}
	o.Logger.Log(logger.INFO, "Cleanup operation successfully executed!")
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
