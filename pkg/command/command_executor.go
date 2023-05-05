package command

import (
	"errors"
	"reflect"
	"strings"

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

func ExecuteCommands(db *gorm.DB, commandExecutionMetadataChannel <-chan CommandExecutionMetadata, outgoingMessageChannel chan<- string) { //TODO: Replace all instances of db with a database client interface
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

	var err error

	switch command.CommandType.Name { //TODO: Change the rest of these to return an error
	case command_type.IncrementCountCommandType:
		executeIncrementCountCommand(db, command.Counter)
	case command_type.IncrementCountByUserCommandType:
		executeIncrementCountByUserCommand(db, command, commandExecutionMetadata.UserName)
	case command_type.SetCountCommandType:
		executeSetCountCommand(db, command.Counter, commandExecutionMetadata.Arguments)
	case command_type.AddTextCommandType:
		err = executeAddTextCommand(db, commandExecutionMetadata.Arguments) //BUG: If you add a command a second time it sends the success and failure message
	case command_type.RemoveTextCommandType:
		err = executeRemoveTextCommand(db, commandExecutionMetadata.Arguments) //BUG: Sends the success and failure message

	}

	if err != nil {
		sendFailureMessage(command, err, outgoingMessageChannel)
	}

	sendCommandText(command, commandExecutionMetadata, outgoingMessageChannel)
}

func getCommandFromName(db *gorm.DB, commandName string) model.Command { //TODO: remove this function and replace with a syncmap
	var command model.Command
	if err := db.Preload("CommandType").Preload("CommandTexts").Preload("CommandTexts.CommandTextType").Preload("Counter").First(&command, "name = ?", commandName).Error; err != nil {
		return model.Command{}
	}
	return command
}

func sendCommandText(command model.Command, commandExecutionMetadata CommandExecutionMetadata, outgoingMessageChannel chan<- string) {
	nonFailureCommandTexts := make([]model.CommandText, 0)

	for _, commandText := range command.CommandTexts {
		if commandText.CommandTextType.Name != "failure" {
			nonFailureCommandTexts = append(nonFailureCommandTexts, commandText)
		}
	}

	templateVariables := getCommandTextVariables(nonFailureCommandTexts)

	templateVariableValues := getCommandTextVariableValues(templateVariables, commandExecutionMetadata, command)

	builtCommandTexts := getBuiltCommandTexts(nonFailureCommandTexts, templateVariableValues)

	for _, builtCommandText := range builtCommandTexts {
		outgoingMessageChannel <- builtCommandText
	}
}

func sendFailureMessage(command model.Command, err error, outgoingMessageChannel chan<- string) {
	var failureCommandText string

	for _, commandText := range command.CommandTexts {
		if commandText.CommandTextType.Name == "failure" {
			failureCommandText = commandText.Text
		}
	}

	if failureCommandText == "" {
		failureCommandText = "Error executing command: " + err.Error()
	}

	outgoingMessageChannel <- failureCommandText
}

func executeIncrementCountCommand(db *gorm.DB, counter model.Counter) {
	if err := db.Model(&counter).Update("count", gorm.Expr("count + ?", 1)).Error; err != nil {
		return //TODO: add logging
	}
}

func executeIncrementCountByUserCommand(db *gorm.DB, command model.Command, userName string) { //TODO: refactor this to accept better arguments
	var counter model.Counter
	var counterByUser model.CounterByUser

	if err := db.First(&counter, "id = ?", command.CounterID).Error; err != nil { //TODO: Add error catching if none found
		return //TODO: add logging
	}
	if err := db.FirstOrCreate(&counterByUser, model.CounterByUser{UserName: userName, CounterID: counter.ID}).Error; err != nil { //ERROR: not creating new counter by user
		return //TODO: add logging
	}
	if err := db.Model(&counter).Update("count", gorm.Expr("count + ?", 1)).Error; err != nil {
		return //TODO: add logging
	}
	if err := db.Model(&counterByUser).Update("count", gorm.Expr("count + ?", 1)).Error; err != nil {
		return //TODO: add logging
	}
}

func executeSetCountCommand(db *gorm.DB, counter model.Counter, commandArguments []string) {
	if len(commandArguments) == 0 {
		return //TODO: add logging
	}

	newCount := commandArguments[0]

	if err := db.Model(&counter).Update("count", newCount).Error; err != nil {
		return //TODO: add logging
	}
}

func executeAddTextCommand(db *gorm.DB, commandArguments []string) error {
	if len(commandArguments) < 2 || len(strings.TrimSpace(commandArguments[0])) == 0 {
		return errors.New("invalid arguments")
	}

	commandName := commandArguments[0]
	commandText := strings.Join(commandArguments[1:], " ")

	newCommand := model.Command{Name: commandName,
		CommandType: model.CommandType{Name: command_type.UserEnteredTextCommandType},
		CommandTexts: []model.CommandText{
			{Text: commandText},
		},
	}

	if err := db.Create(&newCommand).Error; err != nil {
		return err
	}

	return nil
}

func executeRemoveTextCommand(db *gorm.DB, commandArguments []string) error {
	if len(commandArguments) < 1 || len(strings.TrimSpace(commandArguments[0])) == 0 {
		return errors.New("invalid arguments")
	}

	commandName := commandArguments[0]

	if err := db.Delete(&model.Command{}, "name = ?", commandName).Error; err != nil {
		return err
	} else if db.RowsAffected == 0 {
		return errors.New("command not found")
	}

	return nil
}
