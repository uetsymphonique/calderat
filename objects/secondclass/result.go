package secondclass

import "time"

type Result struct {
	out          string
	err          string
	exitCode     string
	status       string
	executedTime time.Time
	finishedTime time.Time
}

func NewResult(stdout string, stderr string, exitCode string, status string, executedTime time.Time, finishedTime time.Time) Result {
	return Result{
		out:          stdout,
		err:          stderr,
		exitCode:     exitCode,
		status:       status,
		executedTime: executedTime,
		finishedTime: finishedTime,
	}
}

func (r Result) Duration() time.Duration {
	return r.finishedTime.Sub(r.executedTime)
}
