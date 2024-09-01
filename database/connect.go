package database

import (
	"fmt"
	"log"
	"os"

	"github.com/samuel-fonseca/task-manager-api/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connecting to a SQLite database
func ConnectDatabase() {
	var err error
	databasePath := os.Getenv("DB_FILE_PATH")

	ensureDatabaseFileExists(databasePath)

	DB, err = gorm.Open(sqlite.Open(databasePath), &gorm.Config{})

	if err != nil {
		log.Fatal("Could not connect to database!")
	}

	DB.AutoMigrate(&model.User{})
	DB.AutoMigrate(&model.Task{})
}

// ensure the database file already exists
// creates a new empty file if not there
func ensureDatabaseFileExists(path string) {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		fmt.Println("Database file not found. Creating new file now")
		f, err := os.Create(path)

		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()
	}
}
