package secondclass

import (
	"calderat/service/execute"
	"calderat/utils/logger"
	"calderat/utils/random"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	EXECUTE = -3
	DISCARD = -2
	PAUSED  = -1
	SUCCESS = 0
	ERROR   = 1
	TIMEOUT = 128
)

type Link struct {
	ProcedureName    string `json:"procedure-name"`
	ProcedureId      string `json:"procedure-id"`
	MitreTechniqueId string `json:"mitre-technique-id"`
	LinkId           string `json:"link-id"`
	Command          string `json:"command"`
	Status           int64
	Jitter           time.Duration `json:"jitter"`
	Executor         Executor      `json:"executor"`
	DecidedTime      time.Time
	FinishedTime     time.Time
	Out              string
	Err              string
	Timeout          time.Duration `json:"timeout"`
	IsCleanup        bool          `json:"is-cleanup"`
	Logger           *logger.Logger
}

func NewLink(procedureName, procedureId, mitreTechniqueId, command string, executor Executor, timeout time.Duration, log *logger.Logger, isCleanup bool) *Link {
	link_id := uuid.New().String()
	link := Link{
		ProcedureName:    procedureName,
		ProcedureId:      procedureId,
		MitreTechniqueId: mitreTechniqueId,
		Command:          command,
		LinkId:           link_id,
		Status:           EXECUTE,
		Jitter:           time.Duration(random.SecureRandomInt(5)) * time.Second,
		Timeout:          timeout,
		Executor:         executor,
		Err:              "",
		IsCleanup:        isCleanup,
		Logger:           log,
	}
	return &link
}

func (link *Link) Execute(executingService execute.ExecutingService) {
	link.Logger.Log(logger.INFO, "Waiting for %s", link.Jitter)
	time.Sleep(link.Jitter)
	link.Decide()
	output, err := executingService.Execute(link.Command, link.Timeout)
	link.Finish()
	link.Logger.Log(logger.INFO, "Command finished after %s\n", link.Duration())
	link.Out = output
	if err != nil {
		link.Err = err.Error()
		if strings.Contains(err.Error(), "timeout") {
			link.Status = TIMEOUT
		} else {
			link.Status = ERROR
		}
	} else {
		link.Status = SUCCESS
	}
}

func (link *Link) Decide() {
	link.DecidedTime = time.Now()
}

func (link *Link) Finish() {
	link.FinishedTime = time.Now()
}

func (link *Link) Duration() time.Duration {
	return link.FinishedTime.Sub(link.DecidedTime)
}

type LinkJSON struct {
	ProcedureName    string   `json:"procedure-name"`
	ProcedureId      string   `json:"procedure-id"`
	MitreTechniqueId string   `json:"mitre-technique-id"`
	LinkId           string   `json:"link-id"`
	Command          string   `json:"command"`
	Jitter           string   `json:"jitter"`
	Executor         Executor `json:"executor"`
	Timeout          string   `json:"timeout"`
	IsCleanup        bool     `json:"is-cleanup"`
}

// Dump only marked attributes to JSON file
func DumpLinksToJson(links []Link, filename string, log *logger.Logger) {
	file, err := os.Create(filename)
	if err != nil {
		log.Log(logger.ERROR, "Error creating file: %s", err)
		return
	}
	defer file.Close()

	// Convert Links to JSON-friendly structure
	var jsonLinks []LinkJSON
	for _, link := range links {
		jsonLinks = append(jsonLinks, LinkJSON{
			ProcedureName:    link.ProcedureName,
			ProcedureId:      link.ProcedureId,
			MitreTechniqueId: link.MitreTechniqueId,
			LinkId:           link.LinkId,
			Command:          link.Command,
			Executor:         link.Executor,
			Jitter:           link.Jitter.String(),
			Timeout:          link.Timeout.String(),
			IsCleanup:        link.IsCleanup,
		})
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(jsonLinks); err != nil {
		log.Log(logger.DEBUG, "Error encoding JSON: %s", err)
		return
	}

	log.Log(logger.TRACE, "Successfully created cleanup links JSON file %s.", filename)
}

// Load cleanup links from JSON file
func LoadCleanupLinksFromJson(filename string, log *logger.Logger) ([]Link, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Unmarshal JSON into an array of LinkJSON objects
	var jsonLinks []LinkJSON
	err = json.Unmarshal(data, &jsonLinks)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Convert back to Link struct
	var links []Link
	for _, jsonLink := range jsonLinks {
		jitter, _ := time.ParseDuration(jsonLink.Jitter)
		timeout, _ := time.ParseDuration(jsonLink.Timeout)

		links = append(links, Link{
			ProcedureName:    jsonLink.ProcedureName,
			ProcedureId:      jsonLink.ProcedureId,
			MitreTechniqueId: jsonLink.MitreTechniqueId,
			LinkId:           jsonLink.LinkId,
			Command:          jsonLink.Command,
			Executor:         jsonLink.Executor,
			Jitter:           jitter,
			Timeout:          timeout,
			IsCleanup:        jsonLink.IsCleanup,
			Logger:           log,
		})
	}

	return links, nil
}
