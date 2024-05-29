package main

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectToDatabase(databaseName string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(databaseName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func GetAllCommands(db *gorm.DB) []Command {
	var commands []Command
	db.Preload("CommandType").Preload("CommandTexts").Preload("Counter").Find(&commands)
	return commands
}
