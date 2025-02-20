package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joaooliveirapro/xmlsyncgo/initializers"
	"github.com/joaooliveirapro/xmlsyncgo/models"
)

func EditsGetAll(c *gin.Context) {
	// Parse client_id param
	_, err := strconv.Atoi(c.Param("client_id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// Parse file_id param
	_, err = strconv.Atoi(c.Param("file_id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// Parse job_id param
	job_id, err := strconv.Atoi(c.Param("job_id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// Get edits from DB
	var edits []models.Edit
	result := initializers.DB.Where("job_id = ?", job_id).Find(&edits)
	if result.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// Send data to client
	c.JSON(http.StatusOK, &edits)
}
