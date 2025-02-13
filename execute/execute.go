package execute

import "time"

type ExecutingService interface {
	Execute(string, time.Duration) (string, error)
	ShortName() string
}
