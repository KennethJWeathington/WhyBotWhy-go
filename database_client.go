package main

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var databaseConnection *gorm.DB

func InitDatabase(databaseName string) {
	db, err := gorm.Open(sqlite.Open(databaseName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic("failed to connect database")
	}

	databaseConnection = db
}

func GetConnection() *gorm.DB {
	return databaseConnection
}

func GetAllCommands(db *gorm.DB) []Command {
	var commands []Command
	db.Preload("CommandType").Preload("CommandTexts").Preload("Counter").Find(&commands)
	return commands
}
