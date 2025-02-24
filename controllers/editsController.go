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
	args := models.PaginateArgs{
		PageSize:   50,
		PageNumber: pageNumber,
		WhereQ:     "job_id = ?",
		OrderQ:     "id DESC",
		WhereA:     []interface{}{job_id},
	}
	response, err := models.Paginate[models.Edit](args)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// Send data to client
	c.JSON(http.StatusOK, &response)
}
