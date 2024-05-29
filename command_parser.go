package main

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
	outgoingCommandsChannel chan<- CommandExecutionMetadata) {

	for messageDetails := range incomingMessagesChannel {
		if command, err := ParseCommand(messageDetails); err == nil {
			outgoingCommandsChannel <- command
		}
	}
}

func ParseCommand(messageDetails struct {
	UserName    string
	Message     string
	IsModerator bool
}) (CommandExecutionMetadata, error) {
	userName, message, isModerator := messageDetails.UserName, messageDetails.Message, messageDetails.IsModerator
	if len(message) == 0 || message[0:1] != commandSignifier {
		return CommandExecutionMetadata{}, errors.New("not a command")
	}

	words := strings.Fields(message)

	commandName := words[0][1:]
	arguments := words[1:]

	return CommandExecutionMetadata{UserName: userName, IsModerator: isModerator, CommandName: commandName, Arguments: arguments}, nil
}
