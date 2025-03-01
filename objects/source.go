package objects

import (
	"calderat/secondclass"
	"calderat/utils/logger"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Source struct {
	Facts  []secondclass.Fact `yaml:"facts"`
	Logger *logger.Logger
}

func NewSource(facts []secondclass.Fact, log *logger.Logger) *Source {
	return &Source{
		Facts:  facts,
		Logger: log,
	}
}

func (s *Source) LoadFromYAML(filePath string) error {
	s.Logger.Log(logger.TRACE, "Loading from yaml file: %s", filePath)

	rawData, err := os.ReadFile(filePath)
	if err != nil {
		s.Logger.Log(logger.ERROR, "Failed to read file '%s': %v", filePath, err)
		return fmt.Errorf("error reading file '%s': %w", filePath, err)
	}

	err = yaml.Unmarshal(rawData, s)
	if err != nil {
		s.Logger.Log(logger.ERROR, "Failed to unmarshal YAML for file '%s': %v", filePath, err)
		return fmt.Errorf("error unmarshalling YAML for file '%s': %w", filePath, err)
	}

	s.Logger.Log(logger.TRACE, "Successfully loaded Source from file: %s", filePath)

	return nil
}
