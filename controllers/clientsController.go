package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaooliveirapro/xmlsyncgo/initializers"
	"github.com/joaooliveirapro/xmlsyncgo/models"
)

func ClientGetAll(c *gin.Context) {
	// Get all clients from DB
	var clients []models.Client
	result := initializers.DB.Find(&clients)
	if result.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// Send data to client
	c.JSON(http.StatusOK, &clients)
}
