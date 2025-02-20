package secondclass

type Executor struct {
	Name     string
	Platform string
	Command  string
	Code     string
	Payloads []string
	Uploads  []string
	Timeout  int
	Cleanup  []string
}

func NewExecutor(name string, platform string, command string, code string, payloads []string, uploads []string, timeout int, cleanup []string) *Executor {
	return &Executor{
		Name:     name,
		Platform: platform,
		Command:  command,
		Code:     code,
		Payloads: payloads,
		Uploads:  uploads,
		Timeout:  timeout,
		Cleanup:  cleanup,
	}
}
