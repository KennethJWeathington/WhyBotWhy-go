package database_client

import (
	"reflect"

	"gorm.io/gorm"
)

var baseCounters = []Counter{
	{Name: "deaths", Count: 0},
	{Name: "boops", Count: 0},
}

var baseCommandTypes = []CommandType{
	{Name: "text"},
	{Name: "increment_count"},
	{Name: "increment_count_by_user"},
	{Name: "set_count"},
	{Name: "add_text_command"},
	{Name: "remove_text_command"},
}

var baseCommandTextTypes = []CommandTextType{
	{Name: "success"},
	{Name: "failure"},
	{Name: "header"},
	{Name: "body"},
}

var baseCommands = []Command{
	{Name: "whyme",
		CommandType: CommandType{Name: "text"},
		CommandTexts: []CommandText{
			{Text: "WHY {{.chatUserName}} WHY!?",
				NeedsUserInfo: true,
			},
		},
	},
	{Name: "death",
		CommandType: CommandType{Name: "increment_count"},
		CommandTexts: []CommandText{
			{Text: "{{.streamName}} has died embarrassingly {{.deaths}} times on stream!",
				NeedsStreamInfo:  true,
				NeedsCounterInfo: true,
			},
		},
		Counter: Counter{Name: "deaths"},
	},
	{Name: "setdeaths",
		CommandType: CommandType{Name: "set_count"},
		CommandTexts: []CommandText{
			{Text: "Deaths set to {{.deaths}}.",
				NeedsCounterInfo: true,
			},
		},
		Counter:         Counter{Name: "deaths"},
		IsModeratorOnly: true,
	},
	{Name: "boop",
		CommandType: CommandType{Name: "increment_count"},
		CommandTexts: []CommandText{
			{Text: "{{.chatUserName}} booped the snoot! The snoot has been booped {{.boops}} times.",
				NeedsUserInfo:    true,
				NeedsCounterInfo: true,
			},
		},
		Counter: Counter{Name: "boops"},
	},
	{Name: "boopboard",
		CommandType: CommandType{Name: "text"},
		CommandTexts: []CommandText{
			{Text: "Top Boopers",
				CommandTextType: CommandTextType{Name: "header"},
			},
			{Text: "{{.row}}. @{{.chatUserName}}: ${{countByUser}} boops",
				CommandTextType:  CommandTextType{Name: "body"},
				NeedsUserInfo:    true,
				NeedsCounterInfo: true,
			},
			{Text: "{{.row}}. @{{.chatUserName}}: ${{countByUser}} boops",
				CommandTextType:  CommandTextType{Name: "body"},
				NeedsUserInfo:    true,
				NeedsCounterInfo: true,
			},
			{Text: "{{.row}}. @{{.chatUserName}}: ${{countByUser}} boops",
				CommandTextType:  CommandTextType{Name: "body"},
				NeedsUserInfo:    true,
				NeedsCounterInfo: true,
			},
		},
		Counter: Counter{Name: "boops"},
	},
	{Name: "addcommand",
		CommandType: CommandType{Name: "add_text_command"},
		CommandTexts: []CommandText{
			{Text: "Command added.",
				CommandTextType: CommandTextType{Name: "success"},
			},
			{Text: "Command already exists.",
				CommandTextType: CommandTextType{Name: "failure"},
			},
		},
		IsModeratorOnly: true,
	},
	{Name: "removecommand",
		CommandType: CommandType{Name: "remove_text_command"},
		CommandTexts: []CommandText{
			{Text: "Command removed.",
				CommandTextType: CommandTextType{Name: "success"},
			},
			{Text: "Command not found.",
				CommandTextType: CommandTextType{Name: "failure"},
			},
		},
		IsModeratorOnly: true,
	},
	{Name: "rules",
		CommandType: CommandType{Name: "text"},
		CommandTexts: []CommandText{
			{Text: "Please remember the channel rules:"},
			{Text: "1. Be kind"},
			{Text: "2. No politics or religion"},
			{Text: "3. No spam "},
			{Text: "4. Only backseat if I ask for it"},
		},
	},
	{Name: "commands",
		CommandType: CommandType{Name: "text"},
		CommandTexts: []CommandText{
			{Text: `The current commands are: {{.commands}}`},
		},
	},
}

func CreateInitialDatabaseData(db *gorm.DB) error {
	db.AutoMigrate(&CounterByUser{})
	db.AutoMigrate(&Counter{})
	if err := createInitialDatabaseCountersIfNotExists(db, baseCounters); err != nil {
		return err
	}

	db.AutoMigrate(&CommandTextType{})
	if err := createInitialDatabaseCommandTextTypesIfNotExists(db, baseCommandTextTypes); err != nil {
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

func createInitialDatabaseCommandTextTypesIfNotExists(db *gorm.DB, baseCommandTextTypes []CommandTextType) error {
	for _, commandTextType := range baseCommandTextTypes {
		if err := db.FirstOrCreate(&commandTextType, commandTextType).Error; err != nil {
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

		for idx := range command.CommandTexts {
			if err := db.First(&command.CommandTexts[idx].CommandTextType, command.CommandTexts[idx].CommandTextType).Error; err != nil {
				return err
			}
		}

		if !reflect.DeepEqual(command.Counter, Counter{}) {
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
