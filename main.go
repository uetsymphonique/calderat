package main

import (
	"calderat/objects"
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

	operation := objects.NewOperation(adversary, true, abilities, env.ShortnameShells, env.OS, ipaddrs[0], log, knowledgeService)

	operation.Run()

}
