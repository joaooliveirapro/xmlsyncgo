package initializers

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func NewProductionFileLogger(filePath string) (logger.Interface, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return logger.New(
		log.New(file, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error, // Only log errors
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error
			Colorful:                  false,        // Disable color
		},
	), nil
}

func ConnecToDB() {
	var err error
	myLogger, err := NewProductionFileLogger("./temp/db.log")
	if err != nil {
		log.Fatal("Error opening db log")
	}
	DB, err = gorm.Open(sqlite.Open(os.Getenv("DB_URL")), &gorm.Config{
		Logger: myLogger,
	})
	if err != nil {
		log.Fatal("Error connecting to DB")
	}
	DB.Logger.LogMode(logger.Error)
}
