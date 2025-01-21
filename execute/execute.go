package execute

type ExecutingService interface {
	Execute(command string, timeout int)
}
