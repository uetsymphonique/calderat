package main

import (
	"calderat/objects"
	"calderat/utils/envdetector"
	logger "calderat/utils/logger"
	"flag"
	"fmt"
)

func main() {

	logLevelFlag := flag.String("log-level", "INFO", "Set the log level (TRACE, DEBUG, INFO, WARN, ERROR)")
	flag.Parse()

	// Initialize a centralized logger with a specified log level
	log, err := logger.New(*logLevelFlag)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v", err)
		return
	}

	env, err := envdetector.DetectEnvironment(log)
	if err != nil {
		log.Log(logger.ERROR, "Failed to detect environment: %v", err)
		return
	}

	fmt.Printf("Operating System: %s\n", env.OS)
	fmt.Printf("Shells: %s\n", env.ShortnameShells)
	fmt.Println("Available Shells:")
	for _, shell := range env.AvailableShells {
		fmt.Printf("- %s\n", shell)
	}
	ipaddrs, err := env.GetAllIPAddresses()

	fmt.Printf("Available IP Addresses: %v\n", ipaddrs)

	// ----------------------------------------------------------------

	// // Load abilities from YAML
	// abilities, err := objects.LoadMultipleFromYAML("data/26c8b8b5-7b5b-4de1-a128-7d37fb14f517.yml", log)
	// if err != nil {
	// 	return
	// }

	// // Print loaded abilities
	// for _, ability := range abilities {
	// 	fmt.Printf("Loaded Ability: %s (ID: %s)\n", ability.Name, ability.AbilityId)
	// 	fmt.Println(ability.Executors)
	// }

	adversary := objects.Adversary{}
	adversary.Logger = log
	adversary.LoadFromYAML("data/adversary.yml")
	fmt.Printf("Loaded Adversary: %s (ID: %s)\n", adversary.Name, adversary.AdversaryId)
}
