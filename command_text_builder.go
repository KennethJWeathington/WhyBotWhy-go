package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"gorm.io/gorm"
)

func getCommandTextVariables(commandTexts []CommandText) []string {
	templateVariables := []string{}

	for _, commandText := range commandTexts {
		templateVariables = append(templateVariables, parseTemplateVariables(commandText.Text)...)
	}
	return templateVariables
}

func parseTemplateVariables(template string) []string {
	regExp, _ := regexp.Compile(`{{\.(.*?)}}`)
	templateVariables := regExp.FindAllString(template, -1)

	for i, templateVariable := range templateVariables {
		templateVariables[i] = strings.Trim(templateVariable, "{.}")
	}
	return templateVariables
}

func getCommandTextVariableValues(db *gorm.DB, templateVariables []string, commandExecutionMetadata CommandExecutionMetadata, command Command) map[string]string {
	templateVariableValues := map[string]string{}
	for _, templateVariable := range templateVariables {
		templateVariableValues[templateVariable] = getTemplateVariableValue(templateVariable, db, commandExecutionMetadata, command)
	}
	return templateVariableValues
}

func getBuiltCommandTexts(commandTexts []CommandText, templateVariableValues map[string]string) []string {
	builtCommandTexts := []string{}

	for _, commandText := range commandTexts {
		builtCommandTexts = append(builtCommandTexts, buildTemplatedString(commandText.Text, templateVariableValues))
	}
	return builtCommandTexts
}

func buildTemplatedString(templateText string, templateVariableValues map[string]string) string {
	builder := &strings.Builder{}

	template, err := template.New("").Parse(templateText)
	if err != nil {
		return ""
	}

	err = template.Execute(builder, templateVariableValues)
	if err != nil {
		return ""
	}

	return builder.String()
}

func getTemplateVariableValue(templateVariable string, db *gorm.DB, commandExecutionMetadata CommandExecutionMetadata, command Command) string { //TODO: Refactor this to use less arguments
	switch templateVariable {
	case "chatUserName":
		return commandExecutionMetadata.UserName
	case "streamName":
		return GetChannelName()
	case "commands":
		return strings.Join(getAllCommandNames(db), ", ")
	case "count":
		return strconv.Itoa(getCountFromDatabase(db, command))
	case "randomQuote":
		return getRandomQuote(db) //TODO: implement specific quotes
	default:
		return ""
	}
}

func getAllCommandNames(db *gorm.DB) []string {
	var commands []Command
	db.Find(&commands)

	commandNames := []string{}
	for _, command := range commands {
		commandNames = append(commandNames, command.Name)
	}
	return commandNames
}

func getCountFromDatabase(db *gorm.DB, command Command) int {
	var counter Counter

	if err := db.First(&counter, "id = ?", command.CounterID).Error; err != nil { //TODO: Add error catching if none found
		return 0
	}

	return counter.Count
}

func getRandomQuote(db *gorm.DB) string {
	var quote Quote

	if err := db.Order("RANDOM()").First(&quote).Error; err != nil {
		return "No quotes found."
	}

	formattedDate := quote.CreatedAt.Format("01/02/06")

	streamName := GetChannelName()

	return fmt.Sprintf("\"%s\" - %s, %s", quote.Text, streamName, formattedDate)
}
