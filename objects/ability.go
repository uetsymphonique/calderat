package objects

import (
	"calderat/secondclass"
	"calderat/service/knowledge"
	"calderat/utils/logger"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

const (
	DefaultTactic = "null_tactic"
)

// Ability represents a configurable ability loaded from a YAML file.
type Ability struct {
	AbilityId        string                 `yaml:"id"`
	Tactic           string                 `yaml:"tactic"`
	Technique        string                 `yaml:"technique_name"`
	TechniqueId      string                 `yaml:"technique_id"`
	Name             string                 `yaml:"name"`
	Description      string                 `yaml:"description"`
	Executors        []secondclass.Executor `yaml:"executors"`
	Privilege        string                 `yaml:"privilege"`
	DeletePayload    bool                   `yaml:"delete_payload"`
	KnowledgeService *knowledge.KnowledgeService
	Logger           *logger.Logger
}

// NewAbility creates a new Ability object with the given parameters.
func NewAbility(ability_id, tactic, technique, technique_id, name, description string, executors []secondclass.Executor, privilege string, delete_payload bool, log *logger.Logger) *Ability {
	if ability_id == "" {
		ability_id = uuid.New().String()
	}
	return &Ability{
		AbilityId:     ability_id,
		Tactic:        tactic,
		Technique:     technique,
		TechniqueId:   technique_id,
		Name:          name,
		Description:   description,
		Executors:     executors,
		Privilege:     privilege,
		DeletePayload: delete_payload,
		Logger:        log,
	}
}

// prehook processes and validates raw YAML data, adding defaults and modifying fields.
func prehook(data map[string]interface{}) ([]byte, error) {
	// Add a UUID if the `id` field is missing
	if _, exists := data["id"]; !exists {
		data["id"] = uuid.New().String()
	}

	// Ensure `tactic` is lowercase or set a default value
	if tactic, exists := data["tactic"]; exists {
		if tacticStr, ok := tactic.(string); ok {
			data["tactic"] = strings.ToLower(tacticStr)
		}
	} else {
		data["tactic"] = DefaultTactic
	}

	// Add additional field validations as needed (e.g., required fields)
	if name, exists := data["name"]; !exists || name == "" {
		return nil, errors.New("missing or empty 'name' field in YAML")
	}

	processedData, err := yaml.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to re-encode YAML after applying modifications: %w", err)
	}

	return processedData, nil
}

// LoadFromYAML loads an Ability from the specified YAML file.
func (a *Ability) LoadFromYAML(filePath string) error {

	a.Logger.Log(logger.TRACE, "Loading YAML file: %s", filePath)

	rawData, err := os.ReadFile(filePath)
	if err != nil {
		a.Logger.Log(logger.ERROR, "Failed to read file '%s': %v", filePath, err)
		return fmt.Errorf("error reading file '%s': %w", filePath, err)
	}

	var data map[string]interface{}
	err = yaml.Unmarshal(rawData, &data)
	if err != nil {
		return fmt.Errorf("failed to parse YAML : %w", err)
	}

	processedData, err := prehook(data)
	if err != nil {
		a.Logger.Log(logger.ERROR, "Error in prehook for file '%s': %v", filePath, err)
		return fmt.Errorf("error in prehook for file '%s': %w", filePath, err)
	}

	err = yaml.Unmarshal(processedData, a)
	if err != nil {
		a.Logger.Log(logger.ERROR, "Failed to unmarshal YAML for file '%s': %v", filePath, err)
		return fmt.Errorf("error unmarshalling YAML for file '%s': %w", filePath, err)
	}

	a.Logger.Log(logger.TRACE, "Successfully loaded Ability from file: %s", filePath)
	return nil
}

func (a *Ability) IsAvailable(shells []string) bool {
	for _, executor := range a.Executors {
		if slices.Contains(shells, executor.Name) {
			return true
		}
	}
	return false
}

func (a *Ability) CreateLinks(log *logger.Logger, shells []string, facts map[string][]*secondclass.Fact) ([]secondclass.Link, []secondclass.Link) {
	links := []secondclass.Link{}
	cleanupLinks := []secondclass.Link{}
	for _, executor := range a.Executors {
		if slices.Contains(shells, executor.Name) {
			commands := a.KnowledgeService.ReplaceFacts(executor.Command, facts)
			for _, command := range commands {
				links = append(links, *secondclass.NewLink(a.Name, a.AbilityId, a.TechniqueId, command, executor, time.Duration(executor.Timeout)*time.Second, log, false))
			}
			for i := len(executor.Cleanup) - 1; i >= 0; i-- {
				commands = a.KnowledgeService.ReplaceFacts(executor.Cleanup[i], facts)
				for _, command := range commands {
					cleanupLinks = append(cleanupLinks, *secondclass.NewLink(a.Name, a.AbilityId, a.TechniqueId, command, executor, time.Duration(executor.Timeout)*time.Second, log, true))
				}
			}
			break
		}
	}
	return links, cleanupLinks
}

// LoadMultipleFromYAML loads multiple abilities from the specified YAML file.
func LoadMultipleAbilityFromYAML(filePath string, log *logger.Logger, knowledgeService *knowledge.KnowledgeService) ([]Ability, error) {
	log.Log(logger.TRACE, "Loading YAML file: %s", filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Log(logger.ERROR, "Failed to read file '%s': %v", filePath, err)
		return nil, fmt.Errorf("error reading file '%s': %w", filePath, err)
	}

	// data, err = prehook(data)
	// if err != nil {
	// 	log.Log(logger.ERROR, "Error in prehook for file '%s': %v", filePath, err)
	// 	return nil, fmt.Errorf("error in prehook for file '%s': %w", filePath, err)
	// }

	var abilities []Ability
	err = yaml.Unmarshal(data, &abilities)
	if err != nil {
		log.Log(logger.ERROR, "Failed to unmarshal YAML for file '%s': %v", filePath, err)
		return nil, fmt.Errorf("error unmarshalling YAML for file '%s': %w", filePath, err)
	}

	// Set the logger for each ability
	for i := range abilities {
		abilities[i].Logger = log
		abilities[i].KnowledgeService = knowledgeService
	}

	log.Log(logger.TRACE, "Successfully loaded %d abilities from file: %s", len(abilities), filePath)
	return abilities, nil
}
