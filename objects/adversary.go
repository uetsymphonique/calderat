package objects

import (
	logger "calderat/utils"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Adversary struct {
	adversaryID    string   `yaml:"adversary_id"`
	name           string   `yaml:"name"`
	description    string   `yaml:"description"`
	atomicOrdering []string `yaml:"atomic_ordering"`
	Logger         *logger.Logger
}

func NewAdversary(adversaryID, name, description string, atomicOrdering []string, log *logger.Logger) *Adversary {
	return &Adversary{
		adversaryID:    adversaryID,
		name:           name,
		description:    description,
		atomicOrdering: atomicOrdering,
		Logger:         log,
	}
}

func NewAdversaryWithLogger(logger *logger.Logger) *Adversary {
	return &Adversary{
		Logger: logger,
	}
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

	a.Logger.Log(logger.DEBUG, "Successfully loaded Adversary from file: %s", filePath)
	return nil

}

func (a *Adversary) Name() string {
	return a.name
}

func (a *Adversary) Description() string {
	return a.description
}

func (a *Adversary) AtomicOrdering() []string {
	return a.atomicOrdering
}

func (a *Adversary) AdversaryID() string {
	return a.adversaryID
}

func (a *Adversary) String() string {
	return fmt.Sprintf("Adversary(id=%s, name=%s, description=%s, atomicOrdering=%v)",
		a.adversaryID, a.name, a.description, a.atomicOrdering)
}
