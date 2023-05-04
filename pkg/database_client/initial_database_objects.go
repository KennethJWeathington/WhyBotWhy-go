package database_client

import (
	"reflect"

	"github.com/jake-weath/whybotwhy_go/pkg/command/command_type"
	"github.com/jake-weath/whybotwhy_go/pkg/database_client/model"
	"gorm.io/gorm"
)

var baseCounters = []model.Counter{
	{Name: "deaths", Count: 0},
	{Name: "boops", Count: 0},
}

var baseCommandTypes = []model.CommandType{
	{Name: command_type.TextCommandType},
	{Name: command_type.IncrementCountCommandType},
	{Name: command_type.IncrementCountByUserCommandType},
	{Name: command_type.SetCountCommandType},
	{Name: command_type.AddTextCommandType},
	{Name: command_type.RemoveTextCommandType},
}

var baseCommandTextTypes = []model.CommandTextType{
	{Name: "success"},
	{Name: "failure"},
	{Name: "header"},
	{Name: "body"},
}

var baseCommands = []model.Command{
	{Name: "whyme",
		CommandType: model.CommandType{Name: command_type.TextCommandType},
		CommandTexts: []model.CommandText{
			{Text: "WHY {{.chatUserName}} WHY!?",
				NeedsUserInfo: true,
			},
		},
	},
	{Name: "death",
		CommandType: model.CommandType{Name: command_type.IncrementCountCommandType},
		CommandTexts: []model.CommandText{
			{Text: "{{.streamName}} has died embarrassingly {{.deaths}} times on stream!",
				NeedsStreamInfo:  true,
				NeedsCounterInfo: true,
			},
		},
		Counter: model.Counter{Name: "deaths"},
	},
	{Name: "setdeaths",
		CommandType: model.CommandType{Name: command_type.SetCountCommandType},
		CommandTexts: []model.CommandText{
			{Text: "Deaths set to {{.deaths}}.",
				NeedsCounterInfo: true,
			},
		},
		Counter:         model.Counter{Name: "deaths"},
		IsModeratorOnly: true,
	},
	{Name: "boop",
		CommandType: model.CommandType{Name: command_type.IncrementCountByUserCommandType},
		CommandTexts: []model.CommandText{
			{Text: "{{.chatUserName}} booped the snoot! The snoot has been booped {{.boops}} times.",
				NeedsUserInfo:    true,
				NeedsCounterInfo: true,
			},
		},
		Counter: model.Counter{Name: "boops"},
	},
	{Name: "boopboard",
		CommandType: model.CommandType{Name: command_type.TextCommandType},
		CommandTexts: []model.CommandText{
			{Text: "Top Boopers",
				CommandTextType: model.CommandTextType{Name: "header"}, //TODO replace CommandTextType Name with constants
			},
			{Text: "{{.row}}. @{{.chatUserName}}: ${{countByUser}} boops",
				CommandTextType:  model.CommandTextType{Name: "body"},
				NeedsUserInfo:    true,
				NeedsCounterInfo: true,
			},
			{Text: "{{.row}}. @{{.chatUserName}}: ${{countByUser}} boops",
				CommandTextType:  model.CommandTextType{Name: "body"},
				NeedsUserInfo:    true,
				NeedsCounterInfo: true,
			},
			{Text: "{{.row}}. @{{.chatUserName}}: ${{countByUser}} boops",
				CommandTextType:  model.CommandTextType{Name: "body"},
				NeedsUserInfo:    true,
				NeedsCounterInfo: true,
			},
		},
		Counter: model.Counter{Name: "boops"},
	},
	{Name: "addcommand",
		CommandType: model.CommandType{Name: command_type.AddTextCommandType},
		CommandTexts: []model.CommandText{
			{Text: "Command added.",
				CommandTextType: model.CommandTextType{Name: "success"},
			},
			{Text: "Command already exists.",
				CommandTextType: model.CommandTextType{Name: "failure"},
			},
		},
		IsModeratorOnly: true,
	},
	{Name: "removecommand",
		CommandType: model.CommandType{Name: command_type.RemoveTextCommandType},
		CommandTexts: []model.CommandText{
			{Text: "Command removed.",
				CommandTextType: model.CommandTextType{Name: "success"},
			},
			{Text: "Command not found.",
				CommandTextType: model.CommandTextType{Name: "failure"},
			},
		},
		IsModeratorOnly: true,
	},
	{Name: "rules",
		CommandType: model.CommandType{Name: command_type.TextCommandType},
		CommandTexts: []model.CommandText{
			{Text: "Please remember the channel rules:"},
			{Text: "1. Be kind"},
			{Text: "2. No politics or religion"},
			{Text: "3. No spam "},
			{Text: "4. Only backseat if I ask for it"},
		},
	},
	{Name: "commands",
		CommandType: model.CommandType{Name: command_type.TextCommandType},
		CommandTexts: []model.CommandText{
			{Text: `The current commands are: {{.commands}}`},
		},
	},
}

func CreateInitialDatabaseData(db *gorm.DB) error {
	db.AutoMigrate(&model.CounterByUser{})
	db.AutoMigrate(&model.Counter{})
	if err := createInitialDatabaseCountersIfNotExists(db, baseCounters); err != nil {
		return err
	}

	db.AutoMigrate(&model.CommandTextType{})
	if err := createInitialDatabaseCommandTextTypesIfNotExists(db, baseCommandTextTypes); err != nil {
		return err
	}

	db.AutoMigrate(&model.CommandType{})
	if err := createInitialDatabaseCommandTypesIfNotExists(db, baseCommandTypes); err != nil {
		return err
	}

	db.AutoMigrate(&model.Command{})
	db.AutoMigrate(&model.CommandText{})
	if err := createInitialDatabaseCommandsIfNotExists(db, baseCommands); err != nil {
		return err
	}

	return nil
}

func createInitialDatabaseCountersIfNotExists(db *gorm.DB, baseCounters []model.Counter) error {
	for _, counter := range baseCounters {
		if err := db.FirstOrCreate(&counter, counter).Error; err != nil {
			return err
		}
	}
	return nil
}

func createInitialDatabaseCommandTextTypesIfNotExists(db *gorm.DB, baseCommandTextTypes []model.CommandTextType) error {
	for _, commandTextType := range baseCommandTextTypes {
		if err := db.FirstOrCreate(&commandTextType, commandTextType).Error; err != nil {
			return err
		}
	}
	return nil
}

func createInitialDatabaseCommandTypesIfNotExists(db *gorm.DB, baseCommandTypes []model.CommandType) error {
	for _, commandType := range baseCommandTypes {
		if err := db.FirstOrCreate(&commandType, commandType).Error; err != nil {
			return err
		}
	}
	return nil
}

func createInitialDatabaseCommandsIfNotExists(db *gorm.DB, baseCommands []model.Command) error {
	for _, command := range baseCommands {
		if err := db.First(&command.CommandType, command.CommandType).Error; err != nil {
			return err
		}

		for idx := range command.CommandTexts {
			if err := db.First(&command.CommandTexts[idx].CommandTextType, command.CommandTexts[idx].CommandTextType).Error; err != nil {
				return err
			}
		}

		if !reflect.DeepEqual(command.Counter, model.Counter{}) {
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
