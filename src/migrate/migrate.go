package main

import (
	"github.com/joaooliveirapro/xmlsyncgo/src/initializers"
	"github.com/joaooliveirapro/xmlsyncgo/src/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnecToDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.Client{})
	initializers.DB.AutoMigrate(&models.File{})
	initializers.DB.AutoMigrate(&models.Edit{})
	initializers.DB.AutoMigrate(&models.AuditLog{})
	initializers.DB.AutoMigrate(&models.Job{})
	initializers.DB.AutoMigrate(&models.Stat{})
}
