package objects

import (
	"calderat/utils/logger"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Adversary struct {
	AdversaryId    string   `yaml:"adversary_id"`
	Name           string   `yaml:"name"`
	Description    string   `yaml:"description"`
	AtomicOrdering []string `yaml:"atomic_ordering"`
	Logger         *logger.Logger
}

func (a *Adversary) LoadFromYAML(filePath string) error {
	a.Logger.Log(logger.TRACE, "Loading from yaml file: %s", filePath)

	rawData, err := os.ReadFile(filePath)
	if err != nil {
		a.Logger.Log(logger.ERROR, "Failed to read file '%s': %v", filePath, err)
		return fmt.Errorf("error reading file '%s': %w", filePath, err)
	}

	err = yaml.Unmarshal(rawData, a)
	if err != nil {
		a.Logger.Log(logger.ERROR, "Failed to unmarshal YAML for file '%s': %v", filePath, err)
		return fmt.Errorf("error unmarshalling YAML for file '%s': %w", filePath, err)
	}

	a.Logger.Log(logger.TRACE, "Successfully loaded Adversary from file: %s", filePath)

	return nil

}

func NewAdversary(adversary_id, name, description string, atomicOrdering []string, log *logger.Logger) *Adversary {
	return &Adversary{
		AdversaryId:    adversary_id,
		Name:           name,
		Description:    description,
		AtomicOrdering: atomicOrdering,
		Logger:         log,
	}
}

func NewAdversaryWithLogger(logger *logger.Logger) *Adversary {
	adversary := Adversary{}
	adversary.Logger = logger
	return &adversary
}
