package command

type CommandExecutionDetails struct {
	UserName    string
	IsModerator bool
	CommandName string
	Arguments   []string
}
