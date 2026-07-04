package response

import "github.com/gin-gonic/gin"

type ApiResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Errors  any    `json:"errors"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func Success(c *gin.Context, status int, message string, data any) {
	c.JSON(status, ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
		Errors:  nil,
	})
}

func NoContent(c *gin.Context) {
	c.Status(204)
}

func Error(c *gin.Context, status int, message string, errs any) {
	c.JSON(status, ApiResponse{
		Success: false,
		Message: message,
		Data:    nil,
		Errors:  errs,
	})
}
