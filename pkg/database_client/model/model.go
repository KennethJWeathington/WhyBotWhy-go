package model

import (
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
	Count     int `gorm:"default 0"`
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
