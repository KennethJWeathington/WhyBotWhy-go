package main

import (
	"gorm.io/gorm"
)

var baseCounters = []Counter{
	{Name: "deaths", Count: 0},
	{Name: "boops", Count: 0},
}

var baseCommandTypes = []CommandType{
	{Name: TextCommandType},
	{Name: IncrementCountCommandType},
	{Name: IncrementCountByUserCommandType},
	{Name: SetCountCommandType},
	{Name: AddTextCommandType},
	{Name: RemoveTextCommandType},
	{Name: UserEnteredTextCommandType},
	{Name: AddQuoteCommandType},
}

var baseCommands = []Command{
	{Name: "whyme",
		CommandType: CommandType{Name: TextCommandType},
		CommandTexts: []CommandText{
			{Text: "WHY {{.chatUserName}} WHY!?"},
		},
	},
	{Name: "death",
		CommandType: CommandType{Name: IncrementCountCommandType},
		CommandTexts: []CommandText{
			{Text: "{{.streamName}} has died embarrassingly {{.count}} times on stream!"},
		},
		Counter: Counter{Name: "deaths"},
	},
	{Name: "setdeaths",
		CommandType: CommandType{Name: SetCountCommandType},
		CommandTexts: []CommandText{
			{Text: "Deaths set to {{.count}}."},
		},
		Counter:         Counter{Name: "deaths"},
		IsModeratorOnly: true,
	},
	{Name: "boop",
		CommandType: CommandType{Name: IncrementCountByUserCommandType},
		CommandTexts: []CommandText{
			{Text: "{{.chatUserName}} booped the snoot! The snoot has been booped {{.count}} times."},
		},
		Counter: Counter{Name: "boops"},
	},
	{Name: "addcommand",
		CommandType: CommandType{Name: AddTextCommandType},
		CommandTexts: []CommandText{
			{Text: "Command added."},
		},
		IsModeratorOnly: true,
	},
	{Name: "removecommand",
		CommandType: CommandType{Name: RemoveTextCommandType},
		CommandTexts: []CommandText{
			{Text: "Command removed."},
		},
		IsModeratorOnly: true,
	},
	{Name: "rules",
		CommandType: CommandType{Name: TextCommandType},
		CommandTexts: []CommandText{
			{Text: "Please remember the channel rules:"}, //TODO: Come up with rules timer
			{Text: "1. Be kind"},
			{Text: "2. No politics or religion"},
			{Text: "3. No spam "},
			{Text: "4. Only backseat if I ask for it"},
		},
	},
	{Name: "commands",
		CommandType: CommandType{Name: TextCommandType},
		CommandTexts: []CommandText{
			{Text: `The current commands are: {{.commands}}`},
		},
	},
	{Name: "addquote",
		CommandType: CommandType{Name: AddQuoteCommandType},
		CommandTexts: []CommandText{
			{Text: `Quote added`},
		},
	},
	{Name: "quote", //TODO implement specific quotes
		CommandType: CommandType{Name: TextCommandType},
		CommandTexts: []CommandText{
			{Text: `{{.randomQuote}}`},
		},
	},
}

func CreateInitialDatabaseData(db *gorm.DB) error {
	db.AutoMigrate(&CounterByUser{})
	db.AutoMigrate(&Counter{})
	if err := createInitialDatabaseCountersIfNotExists(db, baseCounters); err != nil {
		return err
	}

	db.AutoMigrate(&CommandType{})
	if err := createInitialDatabaseCommandTypesIfNotExists(db, baseCommandTypes); err != nil {
		return err
	}

	db.AutoMigrate(&Command{})
	db.AutoMigrate(&CommandText{})
	if err := createInitialDatabaseCommandsIfNotExists(db, baseCommands); err != nil {
		return err
	}

	db.AutoMigrate(&Quote{})

	return nil
}

func createInitialDatabaseCountersIfNotExists(db *gorm.DB, baseCounters []Counter) error {
	for _, counter := range baseCounters {
		if err := db.FirstOrCreate(&counter, counter).Error; err != nil {
			return err
		}
	}
	return nil
}

func createInitialDatabaseCommandTypesIfNotExists(db *gorm.DB, baseCommandTypes []CommandType) error {
	for _, commandType := range baseCommandTypes {
		if err := db.FirstOrCreate(&commandType, commandType).Error; err != nil {
			return err
		}
	}
	return nil
}

func createInitialDatabaseCommandsIfNotExists(db *gorm.DB, baseCommands []Command) error {
	for _, command := range baseCommands {
		if err := db.First(&command.CommandType, command.CommandType).Error; err != nil {
			return err
		}

		if !command.Counter.Equals(Counter{}) {
			if err := db.First(&command.Counter, command.Counter).Error; err != nil {
				return err
			}
		}

		if err := db.FirstOrCreate(&command, command).Error; err != nil {
			return err
		}
	}
	return nil
}
