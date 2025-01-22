package secondclass

import (
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
	LinkId       string
	Command      string
	Status       int64
	Jitter       time.Duration
	DecidedTime  time.Time
	FinishedTime time.Time
	Output       string
}

func NewLink(command string) *Link {
	link_id := uuid.New().String()
	link := Link{Command: command, LinkId: link_id, Status: EXECUTE, DecidedTime: time.Now(), Jitter: 1 * time.Second}
	return &link
}
