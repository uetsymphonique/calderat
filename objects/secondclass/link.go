package secondclass

import (
	"calderat/service/execute"
	"calderat/utils/logger"
	"calderat/utils/random"
	"fmt"
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
	ProcedureName    string
	ProcedureId      string
	MitreTechniqueId string
	LinkId           string
	Command          string
	Status           int64
	Jitter           time.Duration
	Executor         Executor
	DecidedTime      time.Time
	FinishedTime     time.Time
	Out              string
	Err              string
	Timeout          time.Duration
	IsCleanup        bool
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
	link.Logger.Log(logger.INFO, "Command finished after %s", link.Duration())
	fmt.Println("--------------------------------")
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
