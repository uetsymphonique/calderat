package objects

import (
	"calderat/objects/secondclass"
	"calderat/utils/random"
	"encoding/json"
	"os"
	"os/user"
	"time"
)

type Procedure struct {
	ProcedureName        string      `json:"procedure-name"`
	ProcedureDescription string      `json:"procedure-description"`
	ProcedureId          ProcedureId `json:"procedure-id"`
	MitreTechniqueId     string      `json:"mitre-technique-id"`
	Order                int         `json:"order"`
	Steps                []Step      `json:"steps"`
}

func NewProcedure(link *secondclass.Link, order int) *Procedure {
	return &Procedure{
		ProcedureName:        link.ProcedureName,
		ProcedureDescription: link.ProcedureName,
		ProcedureId: ProcedureId{
			Type: "guid",
			Id:   link.ProcedureId,
		},
		MitreTechniqueId: link.MitreTechniqueId,
		Order:            order,
		Steps:            []Step{},
	}
}

func (p *Procedure) AddStep(link *secondclass.Link, order int) {
	p.Steps = append(p.Steps, *NewStep(link, order))
}

type ProcedureId struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

type Step struct {
	Command   string        `json:"command"`
	Executor  string        `json:"executor"`
	Order     int           `json:"order"`
	Output    []OutputBlock `json:"output"`
	TimeStart string        `json:"time-start"`
	TimeStop  string        `json:"time-stop"`
}

func NewStep(link *secondclass.Link, order int) *Step {
	output := []OutputBlock{}
	output = append(output, *NewOutputBlock(link.Out, "STDOUT"))
	if link.Err != "" {
		output = append(output, *NewOutputBlock(link.Err, "STDERR"))
	}
	return &Step{
		Command:   link.Command,
		Executor:  link.Executor.Name,
		Order:     order,
		Output:    []OutputBlock{},
		TimeStart: link.DecidedTime.UTC().Format("2006-01-02T15:04:05.000Z"),
		TimeStop:  link.FinishedTime.UTC().Format("2006-01-02T15:04:05.000Z"),
	}
}

type OutputBlock struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Level   string `json:"level"`
}

func NewOutputBlock(content, level string) *OutputBlock {
	return &OutputBlock{
		Type:    "console",
		Content: content,
		Level:   level,
	}
}

type AttireLog struct {
	AttireVersion string                 `json:"attire-version"`
	ExecutionData map[string]interface{} `json:"execution-data"`
	Procedures    []Procedure            `json:"procedures"`
}

func NewAttireLog(ip string) *AttireLog {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	currentUser, err := user.Current()
	var user string
	if err != nil {
		user = "unknown"
	} else {
		user = currentUser.Username
	}
	return &AttireLog{
		AttireVersion: "1.1",
		ExecutionData: map[string]interface{}{
			"execution-command": "rabbitqm",
			"execution-id":      random.SecureRandomString(32),
			"execution-source":  "Purple-TeamTest",
			"execution-category": map[string]interface{}{
				"name":         "Purple Team VCS",
				"abbreviation": "PPT",
			},
			"target":         hostname,
			"ip":             ip,
			"path":           os.Getenv("PATH"),
			"user":           user,
			"time-generated": time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
		},
		Procedures: []Procedure{},
	}
}

func (al *AttireLog) AddProcedure(procedure Procedure) {
	al.Procedures = append(al.Procedures, procedure)
}

func (a *AttireLog) DumpToFile(filename string) error {
	// Open file for writing
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Convert struct to JSON with indentation
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print JSON

	// Write JSON to file
	return encoder.Encode(a)
}
