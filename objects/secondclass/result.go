package secondclass

import "time"

type Result struct {
	Out          string
	Err          string
	ExitCode     string
	Status       string
	ExecutedTime time.Time
	FinishedTime time.Time
}

func NewResult(out string, err string, exitCode string, status string, executedTime time.Time, finishedTime time.Time) Result {
	return Result{
		Out:          out,
		Err:          err,
		ExitCode:     exitCode,
		Status:       status,
		ExecutedTime: executedTime,
		FinishedTime: finishedTime,
	}
}

func (r Result) Duration() time.Duration {
	return r.FinishedTime.Sub(r.ExecutedTime)
}
