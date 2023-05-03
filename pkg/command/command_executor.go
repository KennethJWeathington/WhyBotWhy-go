package command

import (
	"reflect"
	"regexp"
	"strings"
	"text/template"

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

	//TODO: Move template logic here and combine with other data to pass down to handler methods, figure out counters

	switch command.CommandType.Name {
	case command_type.TextCommandType:
		executeTextCommand(db, command, commandExecutionMetadata, outgoingMessageChannel)
	// case "increment_count":
	// 	executeIncrementCountCommand(db, command, commandExecutionMetadata, outgoingMessageChannel)
	// case "increment_count_by_user":
	// 	executeIncrementCountByUserCommand(db, command, commandExecutionMetadata, outgoingMessageChannel)
	// case "set_count":
	// 	executeSetCountCommand(db, command, commandExecutionMetadata, outgoingMessageChannel)
	// case "add_text_command":
	// 	executeAddTextCommand(db, command, commandExecutionMetadata, outgoingMessageChannel)
	// case "remove_text_command":
	// 	executeRemoveTextCommand(db, command, commandExecutionMetadata, outgoingMessageChannel)
	default:
		outgoingMessageChannel <- "Command type not found." //TODO: Add logging, move to global constant

	}
}

func getCommandFromName(db *gorm.DB, commandName string) model.Command {
	var command model.Command
	if err := db.Preload("CommandType").Preload("CommandTexts").Preload("CommandTexts.CommandTextType").Preload("Counter").First(&command, "name = ?", commandName).Error; err != nil {
		return model.Command{}
	}
	return command
}

func executeTextCommand(db *gorm.DB, command model.Command, commandExecutionMetadata CommandExecutionMetadata, outgoingMessageChannel chan<- string) {
	if len(command.CommandTexts) == 0 {
		outgoingMessageChannel <- "No text found for command." //TODO: Add logging, move to global constant
		return
	}

	templateVariables := []string{}

	for _, commandText := range command.CommandTexts {
		templateVariables = append(templateVariables, getTemplateVariables(commandText.Text)...)
	}

	templateVariableValues := getTemplateVariableValues(templateVariables, commandExecutionMetadata, command)

	for _, commandText := range command.CommandTexts {
		fullCommandText, err := buildTemplatedString(commandText.Text, templateVariableValues)
		if err != nil {
			outgoingMessageChannel <- "Error building templated string." //TODO: Add logging, move to global constant
			return
		}
		outgoingMessageChannel <- fullCommandText
	}
}

func getTemplateVariables(template string) []string {
	regExp, _ := regexp.Compile(`{{\.(.*?)}}`)
	templateVariables := regExp.FindAllString(template, -1)

	for i, templateVariable := range templateVariables {
		templateVariables[i] = strings.Trim(templateVariable, "{.}")
	}
	return templateVariables
}

func getTemplateVariableValues(templateVariables []string, commandExecutionMetadata CommandExecutionMetadata, command model.Command) map[string]string {
	templateVariableValues := map[string]string{}
	for _, templateVariable := range templateVariables {
		templateVariableValues[templateVariable] = getTemplateVariableValue(templateVariable, commandExecutionMetadata)
	}
	return templateVariableValues
}

func getTemplateVariableValue(templateVariable string, commandExecutionMetadata CommandExecutionMetadata) string {
	switch templateVariable {
	case "chatUserName":
		return commandExecutionMetadata.UserName
	default:
		return ""
	}
}

func buildTemplatedString(templateText string, templateVariableValues map[string]string) (string, error) {
	builder := &strings.Builder{}

	template, err := template.New("").Parse(templateText)
	if err != nil {
		return "", err //TODO: add logging
	}

	err = template.Execute(builder, templateVariableValues)
	if err != nil {
		return "", err //TODO: add logging
	}

	return builder.String(), nil
}
