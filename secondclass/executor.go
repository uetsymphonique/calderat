package secondclass

type Executor struct {
	Name     string   `json:"name"`
	Platform string   `json:"platform"`
	Command  string   `json:"command"`
	Code     string   `json:"code"`
	Payloads []string `json:"payloads"`
	Uploads  []string `json:"upload"`
	Timeout  int      `json:"timeout"`
	Cleanup  []string `json:"cleanup"`
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
