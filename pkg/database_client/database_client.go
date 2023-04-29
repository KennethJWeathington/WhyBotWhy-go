package database_client

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func SetUpDatabase(databaseName string) *gorm.DB {
	db := connectToDatabase(databaseName)

	return db
}

func connectToDatabase(databaseName string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(databaseName), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	return db
}
