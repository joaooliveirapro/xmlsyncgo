package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joaooliveirapro/xmlsyncgo/models"
)

func FilesGetAll(c *gin.Context) {
	// Parse client_id param
	client_id, err := strconv.Atoi(c.Param("client_id"))
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
	// Get files from DB
	args := models.PaginateArgs{
		PageSize:   50,
		PageNumber: pageNumber,
		WhereQ:     "client_id = ?",
		OrderQ:     "id DESC",
		WhereA:     []interface{}{client_id},
	}
	response, err := models.Paginate[models.File](args)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// Send data to client
	c.JSON(http.StatusOK, &response)
}
