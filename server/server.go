package server

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joaooliveirapro/xmlsyncgo/controllers"
)

type WebServer struct {
	DistDir string
}

func (ws *WebServer) Start(port string) {
	// Server index.html
	r := gin.Default()

	// Configure CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8080"}                                     // Replace with your frontend's domain
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}  // Adjust as needed
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"} // Adjust as needed
	config.AllowCredentials = true                                                              // If you need to send cookies
	config.MaxAge = 12 * time.Hour                                                              // Optional: Set max age for preflight requests

	r.Use(cors.New(config))

	// GET
	r.GET("/clients", controllers.ClientGetAll)
	r.GET("/clients/:client_id/files", controllers.FilesGetAll)
	r.GET("/clients/:client_id/files/:file_id/audits", controllers.AuditsGetAll)
	r.GET("/clients/:client_id/files/:file_id/stats", controllers.StatsGetAll)
	r.GET("/clients/:client_id/files/:file_id/jobs", controllers.JobsGetAll)
	r.GET("/clients/:client_id/files/:file_id/jobs/:job_id/edits", controllers.EditsGetAll)

	r.Run()
}
