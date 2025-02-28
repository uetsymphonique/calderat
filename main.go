package main

import (
	"calderat/objects"
	"calderat/secondclass"
	"calderat/service/knowledge"
	"calderat/utils/colorprint"
	"calderat/utils/data"
	"calderat/utils/envdetector"
	logger "calderat/utils/logger"
	"flag"
	"fmt"
	"strings"
)

func main() {

	logLevelFlag := flag.String("log-level", "INFO", "Set the log level (TRACE, DEBUG, INFO, WARN, ERROR)")
	nonCleanupMode := flag.Bool("non-cleanup", false, "Disable cleanup operation")
	nonAutonomousMode := flag.Bool("non-auto", false, "Enable non-auto mode")
	cleanupOp := flag.Bool("cleanup-op", false, "Cleanup current operation")
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
	fmt.Println(colorprint.ColorString("Agent information:", colorprint.BLUE))
	fmt.Printf("%s[+] Operating System: %s\n", colorprint.CYAN, env.OS)
	fmt.Printf("[+] Shells: %s\n", strings.Join(env.ShortnameShells, ", "))
	ipaddrs, err := env.GetAllIPAddresses()

	fmt.Printf("Available IP Addresses: %s%s\n", strings.Join(ipaddrs, ", "), colorprint.RESET)

	if *cleanupOp {
		cleanupLinks, err := secondclass.LoadCleanupLinksFromJson("cleanups.json", log)
		if err != nil {
			log.Log(logger.ERROR, "Failed to load cleanup links: %v", err)
			return
		}
		operation := objects.NewCleanupOperation(cleanupLinks, env.ShortnameShells, env.OS, ipaddrs[0], log)
		operation.RunningCleanupOperation()
		return
	}

	// ----------------------------------------------------------------
	knowledgeService := knowledge.NewKnowledgeService(log)
	abilities, err := data.ProcessYmlAbilities("data/abilities/", log, knowledgeService)

	if err != nil {
		log.Log(logger.ERROR, "Failed to load abilities: %v", err)
		return
	}

	adversary := objects.Adversary{}
	adversary.Logger = log
	adversary.LoadFromYAML("data/adversary.yml")

	operation := objects.NewOperation(adversary, !*nonAutonomousMode, !*nonCleanupMode, abilities, env.ShortnameShells, env.OS, ipaddrs[0], log, knowledgeService)

	operation.Run()

}
