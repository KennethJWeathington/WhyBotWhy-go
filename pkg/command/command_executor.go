package command

import (
	"reflect"

	"github.com/jake-weath/whybotwhy_go/pkg/command/command_type"
	"github.com/jake-weath/whybotwhy_go/pkg/database_client/model"
	"gorm.io/gorm"
)

type CommandExecutionMetadata struct {
	UserName    string
	IsModerator bool
	CommandName string
	Arguments   []string
}

func ExecuteCommands(db *gorm.DB, commandExecutionMetadataChannel <-chan CommandExecutionMetadata, outgoingMessageChannel chan<- string) {
	for commandExecutionMetadata := range commandExecutionMetadataChannel {
		go executeCommand(db, commandExecutionMetadata, outgoingMessageChannel)
	}
}

func executeCommand(db *gorm.DB, commandExecutionMetadata CommandExecutionMetadata, outgoingMessageChannel chan<- string) {
	command := getCommandFromName(db, commandExecutionMetadata.CommandName)
	if reflect.DeepEqual(command, model.Command{}) { //TODO: Replace reflect.DeepEqual with custom equality function
		return
	}

	if command.IsModeratorOnly && !commandExecutionMetadata.IsModerator {
		return
	}

	// if(command.CommandType.Name ) //TODO: possibly check that command type exists

	switch command.CommandType.Name {
	case command_type.IncrementCountCommandType:
		executeIncrementCountCommand(db, command)
	case "increment_count_by_user":
		executeIncrementCountByUserCommand(db, command, commandExecutionMetadata.UserName)
		// case "set_count":
		// 	executeSetCountCommand(db, command, commandExecutionMetadata)
		// case "add_text_command":
		// 	executeAddTextCommand(db, command, commandExecutionMetadata)
		// case "remove_text_command":
		// 	executeRemoveTextCommand(db, command, commandExecutionMetadata)

	}

	templateVariables := getCommandTextVariables(command.CommandTexts)

	templateVariableValues := getCommandTextVariableValues(templateVariables, commandExecutionMetadata, command)

	builtCommandTexts := getBuiltCommandTexts(command.CommandTexts, templateVariableValues)

	for _, builtCommandText := range builtCommandTexts {
		outgoingMessageChannel <- builtCommandText
	}
}

func getCommandFromName(db *gorm.DB, commandName string) model.Command { //TODO: remove this function and replace with a syncmap
	var command model.Command
	if err := db.Preload("CommandType").Preload("CommandTexts").Preload("CommandTexts.CommandTextType").Preload("Counter").First(&command, "name = ?", commandName).Error; err != nil {
		return model.Command{}
	}
	return command
}

func executeIncrementCountCommand(db *gorm.DB, command model.Command) {
	if err := db.Model(&command.Counter).Update("count", gorm.Expr("count + ?", 1)).Error; err != nil {
		return //TODO: add logging
	}
}

func executeIncrementCountByUserCommand(db *gorm.DB, command model.Command, userName string) {
	var counter model.Counter
	var counterByUser model.CounterByUser

	if err := db.First(&counter, "id = ?", command.CounterID).Error; err != nil {
		return //TODO: add logging
	}
	if err := db.FirstOrCreate(&counterByUser, model.CounterByUser{UserName: userName, CounterID: counter.ID, Count: 0}).Error; err != nil { //ERROR: not creating new counter by user
		return //TODO: add logging
	}
	if err := db.Model(&counter).Update("count", gorm.Expr("count + ?", 1)).Error; err != nil {
		return //TODO: add logging
	}
	if err := db.Model(&counterByUser).Update("count", gorm.Expr("count + ?", 1)).Error; err != nil {
		return //TODO: add logging
	}
}
