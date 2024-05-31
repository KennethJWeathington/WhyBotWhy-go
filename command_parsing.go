package main

import (
	"errors"
	"strings"
)

func ParseCommand(chatMessage string, commandSiginfier string) (command string, arguments []string, err error) {
	if len(chatMessage) == 0 || chatMessage[0:len(commandSiginfier)] != commandSignifier {
		return "", nil, errors.New("not a command")
	}

	words := strings.Fields(chatMessage)

	command = words[0][len(commandSiginfier):]
	arguments = words[1:]

	return command, arguments, nil
}
