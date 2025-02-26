package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaooliveirapro/xmlsyncgo/src/controllers"
)

type Transport interface {
	ServerHTTP()
}

type HTTPTransport struct {
	Engine *gin.Engine
}

func NewTransport() *HTTPTransport {
	return &HTTPTransport{
		Engine: gin.Default(),
	}
}

func (h *HTTPTransport) ServerHTTP() {
	// Grab root path
	rootPath := h.Engine.Group("/")

	// Use middleware
	rootPath.Use(h.CORSMiddleware())

	// Set API Routes
	h.APIRoutes(rootPath)

	// Start server in a goroutine
	go func() {
		// PORT is set in .env
		if err := h.Engine.Run(); err != nil {
			log.Fatal(err)
		}
	}()
}

func (h *HTTPTransport) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,PUT,PATCH,POST,DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "range")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// Register the API routes
func (h *HTTPTransport) APIRoutes(r *gin.RouterGroup, middleware ...gin.HandlerFunc) {
	// GET
	v1 := r.Group("/v1")
	v1.GET("/clients", controllers.ClientGetAll)
	v1.GET("/clients/:client_id/files", controllers.FilesGetAll)
	v1.GET("/clients/:client_id/files/:file_id/audits", controllers.AuditsGetAll)
	v1.GET("/clients/:client_id/files/:file_id/stats", controllers.StatsGetAll)
	v1.GET("/clients/:client_id/files/:file_id/jobs", controllers.JobsGetAll)
	v1.GET("/clients/:client_id/files/:file_id/jobs/:job_id/edits", controllers.EditsGetAll)
}
