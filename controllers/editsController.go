package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joaooliveirapro/xmlsyncgo/models"
)

func EditsGetAll(c *gin.Context) {
	// Parse job_id param
	job_id, err := strconv.Atoi(c.Param("job_id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// Parse ?page= args from request
	pageNumber, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// Get all jobs paginated and more info
	response, err := models.Paginate[models.Edit](50, pageNumber, "job_id = ?", "id DESC", job_id)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// Send data to client
	c.JSON(http.StatusOK, &response)
}
