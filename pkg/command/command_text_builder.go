package command

import (
	"regexp"
	"strings"
	"text/template"

	"github.com/jake-weath/whybotwhy_go/pkg/database_client/model"
)

func getCommandTextVariables(commandTexts []model.CommandText) []string {
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

func getCommandTextVariableValues(templateVariables []string, commandExecutionMetadata CommandExecutionMetadata, command model.Command) map[string]string {
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

func getBuiltCommandTexts(commandTexts []model.CommandText, templateVariableValues map[string]string) []string {
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
		return "" //TODO: add logging
	}

	err = template.Execute(builder, templateVariableValues)
	if err != nil {
		return "" //TODO: add logging
	}

	return builder.String()
}
