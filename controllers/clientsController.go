package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joaooliveirapro/xmlsyncgo/models"
)

func ClientGetAll(c *gin.Context) {
	// Parse ?page= args from request
	pageNumber, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// Get all clients paginated and more info
	response, err := models.Paginate[models.Client](50, pageNumber, "", "id DESC")
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// Add custom headers
	c.Writer.Header().Set("X-Total-Count", fmt.Sprintf("%d", response.Total))
	// Send data to client
	c.JSON(http.StatusOK, &response)
}
