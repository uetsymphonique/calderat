package data

import (
	"calderat/objects"
	"calderat/service/knowledge"
	"calderat/utils/logger"
	"fmt"
	"os"
	"path/filepath"
)

func ProcessYmlAbilities(folder string, log *logger.Logger, knowledgeService *knowledge.KnowledgeService) ([]objects.Ability, error) {
	ret_abilities := []objects.Ability{}
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		// Check if it's a file and has .yml extension
		if !info.IsDir() && filepath.Ext(path) == ".yml" {
			log.Log(logger.TRACE, "Processing file: %s", path)
			// Add your file processing logic here
			abilities, err := objects.LoadMultipleAbilityFromYAML(path, log, knowledgeService)
			if err != nil {
				return err
			}
			ret_abilities = append(ret_abilities, abilities...)
		}
		return nil
	})
	if err != nil {
		return ret_abilities, fmt.Errorf("error walking the path %s: %w", folder, err)
	}
	return ret_abilities, nil
}
