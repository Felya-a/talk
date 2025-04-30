package http_handler_404

import (
	. "talk/internal/http/handlers"

	"github.com/gin-gonic/gin"
)

func Handle404(ctx *gin.Context) {
	response := ErrorResponse{
		Status:  "error",
		Message: "not found",
		Error:   "not found",
	}
	ctx.JSON(404, response)
	return
}
