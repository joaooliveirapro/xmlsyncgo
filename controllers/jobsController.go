package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joaooliveirapro/xmlsyncgo/initializers"
	"github.com/joaooliveirapro/xmlsyncgo/models"
)

func JobsGetAll(c *gin.Context) {
	// Parse client_id param
	_, err := strconv.Atoi(c.Param("client_id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// Parse file_id param
	file_id, err := strconv.Atoi(c.Param("file_id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// Get jobs from DB
	var jobs []models.Job
	result := initializers.DB.Where("file_id = ?", file_id).Find(&jobs)
	if result.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// Send data to client
	c.JSON(http.StatusOK, &jobs)
}
