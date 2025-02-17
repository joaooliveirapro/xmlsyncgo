package initializers

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnecToDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open(os.Getenv("DB_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to DB")
	}
	DB.Logger.LogMode(logger.Silent)
}
