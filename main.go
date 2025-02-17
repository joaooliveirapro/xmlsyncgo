package main

import (
	"github.com/joaooliveirapro/xmlsyncgo/initializers"
	"github.com/joaooliveirapro/xmlsyncgo/models"
	"github.com/joaooliveirapro/xmlsyncgo/parser"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnecToDB()
}

func main() {
	// Get all clients from DB
	var clients []models.Client
	initializers.DB.Preload("Files").Find(&clients)

	// Create parser manager
	pm := parser.ParserManager{}
	pm.Run(&clients)

}
