package command

import (
	"errors"
	"strings"
)

const commandSignifier = "!"

func ParseIncomingMessagesToCommands(incomingMessagesChannel <-chan struct {
	UserName    string
	Message     string
	IsModerator bool
},
	outgoingCommandsChannel chan<- CommandExecutionDetails) {

	for messageDetails := range incomingMessagesChannel {
		if command, err := ParseCommand(messageDetails.Message); err == nil {
			outgoingCommandsChannel <- command
		}
	}
}

func ParseCommand(message string) (CommandExecutionDetails, error) {
	if len(message) == 0 || message[0:1] != commandSignifier {
		return CommandExecutionDetails{}, errors.New("not a command")
	}

	words := strings.Fields(message)

	commandName := words[0][1:]
	arguments := words[1:]

	return CommandExecutionDetails{Name: commandName, Arguments: arguments}, nil
}
