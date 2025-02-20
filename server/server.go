package server

import (
	"github.com/gin-gonic/gin"
	"github.com/joaooliveirapro/xmlsyncgo/controllers"
)

type WebServer struct {
	DistDir string
}

func (ws *WebServer) Start(port string) {
	// Server index.html
	r := gin.Default()

	// GET
	r.GET("/clients", controllers.ClientGetAll)
	r.GET("/clients/:client_id/files", controllers.FilesGetAll)
	r.GET("/clients/:client_id/files/:file_id/audits", controllers.AuditsGetAll)
	r.GET("/clients/:client_id/files/:file_id/stats", controllers.StatsGetAll)
	r.GET("/clients/:client_id/files/:file_id/jobs", controllers.JobsGetAll)
	r.GET("/clients/:client_id/files/:file_id/jobs/:job_id/edits", controllers.EditsGetAll)

	r.Run()
}
