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

type CommandTextType struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"unique"`
}

type CommandText struct {
	ID                uint `gorm:"primarykey"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	CommandID         uint   `gorm:"not null"`
	Text              string `gorm:"not null"`
	CustomTextQuery   string
	CommandTextTypeID uint
	CommandTextType   CommandTextType
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
