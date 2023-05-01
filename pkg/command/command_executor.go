package command

import (
	"github.com/jake-weath/whybotwhy_go/pkg/database_client"
)

func ExecuteCommands(commands []database_client.Command, commandExecutionChannel <-chan CommandExecutionDetails, outgoingMessageChannel chan<- string) {
	for commandToExecute := range commandExecutionChannel {
		go executeCommand(commandToExecute, commands, outgoingMessageChannel)
	}
}

func executeCommand(commandToExecute CommandExecutionDetails, commands []database_client.Command, outgoingMessageChannel chan<- string) {
	for _, command := range commands {
		if command.Name == commandToExecute.Name {
			for _, commandText := range command.CommandTexts {
				outgoingMessageChannel <- commandText.Text //TODO: Replace with templating
			}
		}
	}
}
