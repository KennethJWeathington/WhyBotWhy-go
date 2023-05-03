package command

import (
	"reflect"

	"github.com/jake-weath/whybotwhy_go/pkg/database_client"
	"gorm.io/gorm"
)

func ExecuteCommands(db *gorm.DB, commandExecutionChannel <-chan CommandExecutionDetails, outgoingMessageChannel chan<- string) {
	for commandToExecute := range commandExecutionChannel {
		go executeCommand(db, commandToExecute, outgoingMessageChannel)
	}
}

func executeCommand(db *gorm.DB, commandToExecute CommandExecutionDetails, outgoingMessageChannel chan<- string) {
	command := getCommandFromName(db, commandToExecute.CommandName)
	if reflect.DeepEqual(command, database_client.Command{}) { //TODO: Replace reflect.DeepEqual with custom equality function
		return
	}

	switch command.CommandType.Name {
	case "text":
		executeTextCommand(db, command, commandToExecute, outgoingMessageChannel)
	// case "increment_count":
	// 	executeIncrementCountCommand(db, command, commandToExecute, outgoingMessageChannel)
	// case "increment_count_by_user":
	// 	executeIncrementCountByUserCommand(db, command, commandToExecute, outgoingMessageChannel)
	// case "set_count":
	// 	executeSetCountCommand(db, command, commandToExecute, outgoingMessageChannel)
	// case "add_text_command":
	// 	executeAddTextCommand(db, command, commandToExecute, outgoingMessageChannel)
	// case "remove_text_command":
	// 	executeRemoveTextCommand(db, command, commandToExecute, outgoingMessageChannel)
	default:
		outgoingMessageChannel <- "Command type not found." //TODO: Add logging, move to global constant

	}
}

func getCommandFromName(db *gorm.DB, commandName string) database_client.Command {
	var command database_client.Command
	if err := db.Preload("CommandType").Preload("CommandTexts").Preload("CommandTexts.CommandTextType").Preload("Counter").First(&command, "name = ?", commandName).Error; err != nil {
		return database_client.Command{}
	}
	return command
}

func executeTextCommand(db *gorm.DB, command database_client.Command, commandToExecute CommandExecutionDetails, outgoingMessageChannel chan<- string) {
	if len(command.CommandTexts) == 0 {
		outgoingMessageChannel <- "No text found for command." //TODO: Add logging, move to global constant
		return
	}

	for _, commandText := range command.CommandTexts {
		outgoingMessageChannel <- commandText.Text //TODO: Add templating
	}
}
