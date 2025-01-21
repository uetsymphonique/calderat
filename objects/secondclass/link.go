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
	LinkID       string
	Command      string
	Status       int64
	Jitter       time.Duration
	DecidedTime  time.Time
	FinishedTime time.Time
	Output       string
}

func NewLink(command string) *Link {
	linkID := uuid.New().String()
	link := Link{Command: command, LinkID: linkID, Status: EXECUTE, DecidedTime: time.Now(), Jitter: 1 * time.Second}
	return &link
}
