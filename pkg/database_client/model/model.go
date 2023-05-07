package model

import (
	"time"
)

type Counter struct {
	ID             uint `gorm:"primarykey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Name           string `gorm:"unique"`
	Count          int
	CounterByUsers []CounterByUser
}

type CounterByUser struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserName  string
	CounterID uint
	Count     int `gorm:"default 0"`
}

type CommandText struct {
	ID              uint `gorm:"primarykey"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	CommandID       uint   `gorm:"not null"`
	Text            string `gorm:"not null"`
	CustomTextQuery string
	Order           int
}

type CommandType struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"unique"`
}

type Command struct {
	ID              uint `gorm:"primarykey"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Name            string `gorm:"not null;unique"`
	CommandTypeID   uint   `gorm:"not null"`
	CommandType     CommandType
	CommandTexts    []CommandText
	CounterID       uint
	Counter         Counter
	IsModeratorOnly bool `gorm:"default false"`
}

type Quote struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"not null;unique"`
	Text      string `gorm:"not null"`
}

func (command *Command) Equals(otherCommand Command) bool {
	if command.Name != otherCommand.Name {
		return false
	}
	if command.CounterID != otherCommand.CounterID {
		return false
	}
	if command.CommandTypeID != otherCommand.CommandTypeID {
		return false
	}
	if command.IsModeratorOnly != otherCommand.IsModeratorOnly {
		return false
	}
	return true
}

func (counter *Counter) Equals(otherCounter Counter) bool {
	if counter.Name != otherCounter.Name {
		return false
	}
	if counter.Count != otherCounter.Count {
		return false
	}
	return true
}
