package database_client

import "gorm.io/gorm"

type Counter struct {
	gorm.Model
	Name           string
	Count          int
	CounterByUsers []CounterByUser
}

type CounterByUser struct {
	gorm.Model
	UserName  string
	CounterID int
}

type CommandTextType struct {
	gorm.Model
	Name string
}

type CommandText struct {
	gorm.Model
	CommandID         int    `gorm:"not null"`
	Text              string `gorm:"not null"`
	CustomTextQuery   string
	CommandTextTypeID int
	CommandTextType   CommandTextType
	NeedsStreamInfo   bool `gorm:"default false"`
	NeedsUserInfo     bool `gorm:"default false"`
	NeedsCounterInfo  bool `gorm:"default false"`
}

type CommandType struct {
	gorm.Model
	Name string
}

type Command struct {
	gorm.Model
	Name            string `gorm:"not null;unique"`
	CommandTypeID   int    `gorm:"not null"`
	CommandType     CommandType
	CommandTexts    []CommandText
	CounterID       int
	Counter         Counter
	IsModeratorOnly bool `gorm:"default false"`
}
