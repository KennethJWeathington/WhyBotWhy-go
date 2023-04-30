package database_client

import (
	"reflect"

	"gorm.io/gorm"
)

type Counter struct {
	gorm.Model
	Name           string `gorm:"unique"`
	Count          int
	CounterByUsers []CounterByUser
}

type CounterByUser struct {
	gorm.Model
	UserName  string
	CounterID uint
}

type CommandTextType struct {
	gorm.Model
	Name string `gorm:"unique"`
}

type CommandText struct {
	gorm.Model
	CommandID         uint   `gorm:"not null"`
	Text              string `gorm:"not null"`
	CustomTextQuery   string
	CommandTextTypeID uint
	CommandTextType   CommandTextType
	NeedsStreamInfo   bool `gorm:"default false"`
	NeedsUserInfo     bool `gorm:"default false"`
	NeedsCounterInfo  bool `gorm:"default false"`
}

type CommandType struct {
	gorm.Model
	Name string `gorm:"unique"`
}

type Command struct {
	gorm.Model
	Name            string `gorm:"not null;unique"`
	CommandTypeID   uint   `gorm:"not null"`
	CommandType     CommandType
	CommandTexts    []CommandText
	CounterID       uint
	Counter         Counter
	IsModeratorOnly bool `gorm:"default false"`
}

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
			{Text: `Please remember the channel rules: 
			1. Be kind 
			2. No politics or religion 
			3. No spam 
			4. Only backseat if I ask for it`,
			},
		},
	},
}

func CreateInitialDatabaseData(db *gorm.DB) error {
	db.AutoMigrate(&CounterByUser{})
	db.AutoMigrate(&Counter{})
	if err := createInitialDatabaseCountersIfNotExists(db); err != nil {
		return err
	}

	db.AutoMigrate(&CommandTextType{})
	if err := createInitialDatabaseCommandTextTypesIfNotExists(db); err != nil {
		return err
	}

	db.AutoMigrate(&CommandType{})
	if err := createInitialDatabaseCommandTypesIfNotExists(db); err != nil {
		return err
	}

	db.AutoMigrate(&Command{})
	db.AutoMigrate(&CommandText{})
	if err := createInitialDatabaseCommandsIfNotExists(db); err != nil {
		return err
	}

	return nil
}

func createInitialDatabaseCountersIfNotExists(db *gorm.DB) error {
	for _, counter := range baseCounters {
		if err := db.FirstOrCreate(&counter, counter).Error; err != nil {
			return err
		}
	}
	return nil
}

func createInitialDatabaseCommandTextTypesIfNotExists(db *gorm.DB) error {
	for _, commandTextType := range baseCommandTextTypes {
		if err := db.FirstOrCreate(&commandTextType, commandTextType).Error; err != nil {
			return err
		}
	}
	return nil
}

func createInitialDatabaseCommandTypesIfNotExists(db *gorm.DB) error {
	for _, commandType := range baseCommandTypes {
		if err := db.FirstOrCreate(&commandType, commandType).Error; err != nil {
			return err
		}
	}
	return nil
}

func createInitialDatabaseCommandsIfNotExists(db *gorm.DB) error {
	for _, command := range baseCommands {
		if err := db.First(&command.CommandType, command.CommandType).Error; err != nil {
			return err
		}

		for idx, _ := range command.CommandTexts {
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
