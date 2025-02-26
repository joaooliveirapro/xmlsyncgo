package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joaooliveirapro/xmlsyncgo/src/models"
)

func ClientGetAll(c *gin.Context) {
	// Parse ?page= args from request
	pageNumber, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// Get all clients paginated and more info
	args := models.PaginateArgs{
		PageSize:   50,
		PageNumber: pageNumber,
		OrderQ:     "id DESC",
	}
	response, err := models.Paginate[models.Client](args)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// Add custom headers
	c.Writer.Header().Set("X-Total-Count", fmt.Sprintf("%d", response.Total))

	// Include editsCount

	// Send data to client
	c.JSON(http.StatusOK, &response)
}
