package knowledge

import (
	"calderat/secondclass"
	"calderat/utils/logger"
	"regexp"
)

const (
	FACTRGX = `#{(.*?)}`
)

type KnowledgeService struct {
	Logger *logger.Logger
}

func NewKnowledgeService(logger *logger.Logger) *KnowledgeService {
	return &KnowledgeService{Logger: logger}
}

func (ks *KnowledgeService) RequiredTraits(command string) []string {
	re := regexp.MustCompile(FACTRGX)
	matches := re.FindAllStringSubmatch(command, -1)
	traits := []string{}

	// Print extracted values
	for _, match := range matches {
		if len(match) > 1 {
			ks.Logger.Log(logger.DEBUG, "Command %s requires fact: #{%s}", command, match[1])
			traits = append(traits, match[1])
		}
	}
	return traits
}

func GenerateCombinations(keys []string, facts map[string][]*secondclass.Fact, index int, current map[string]*secondclass.Fact, results *[]string, template string) {
	// Base case: all keys are replaced
	if index == len(keys) {
		// Replace placeholders in the template
		finalStr := template
		for key, fact := range current {
			finalStr = regexp.MustCompile(`#{`+key+`}`).ReplaceAllLiteralString(finalStr, fact.Value)
		}
		*results = append(*results, finalStr)
		return
	}

	// Get the current key and its values
	key := keys[index]
	values := facts[key]

	// Iterate over all replacement values for the current key
	for _, value := range values {
		current[key] = value
		GenerateCombinations(keys, facts, index+1, current, results, template)
	}
}

func (ks *KnowledgeService) ReplaceFacts(command string, facts map[string][]*secondclass.Fact) []string {
	requiredTraits := ks.RequiredTraits(command)
	var results []string
	GenerateCombinations(requiredTraits, facts, 0, make(map[string]*secondclass.Fact), &results, command)
	return results
}
