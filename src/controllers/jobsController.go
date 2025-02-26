package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joaooliveirapro/xmlsyncgo/src/models"
)

func JobsGetAll(c *gin.Context) {
	// Parse file_id param
	file_id, err := strconv.Atoi(c.Param("file_id"))
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
	args := models.PaginateArgs{
		PageSize:   50,
		PageNumber: pageNumber,
		WhereQ:     "file_id = ?",
		OrderQ:     "id DESC",
		Preload:    true,
		PreloadQ:   "Edits",
		WhereA:     []interface{}{file_id},
	}
	response, err := models.Paginate[models.Job](args)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// Send data to client
	c.JSON(http.StatusOK, &response)
}
