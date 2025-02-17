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

	// Web Server
	// r := gin.Default()

	// // Clients
	// r.GET("/clients", controllers.ClientGetAll)
	// r.POST("/clients", controllers.ClientCreate)

	// // Files
	// r.GET("/clients/:client_id/files", controllers.FileGet)
	// r.POST("/clients/:client_id/files", controllers.FileCreate)

	// // Serve on env.PORT
	// r.Run()
}
