package database_client

import (
	"github.com/glebarez/sqlite"
	"github.com/jake-weath/whybotwhy_go/pkg/database_client/model"
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

func GetAllCommands(db *gorm.DB) []model.Command {
	var commands []model.Command
	db.Preload("CommandType").Preload("CommandTexts").Preload("Counter").Find(&commands)
	return commands
}
