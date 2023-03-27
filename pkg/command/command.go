package command

import (
	"errors"
	"strings"
)

const commandSignifier = "!"

type Command struct {
	Name      string
	Arguments []string
}

func ParseCommand(message string) (Command, error) {
	if len(message) == 0 || message[0:1] != commandSignifier {
		return Command{}, errors.New("not a command")
	}

	words := strings.Fields(message)

	commandName := words[0][1:]
	arguments := words[1:]

	return Command{Name: commandName, Arguments: arguments}, nil
}
