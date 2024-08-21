package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseError struct {
	Message []string `json:"message"`
	Error   string   `json:"error"`
	Status  int      `json:"statusCode"`
}

func SuccessResponse(c *gin.Context, code int, data interface{}) {
	c.JSON(http.StatusOK, ResponseData{
		Code:    code,
		Message: msg[code],
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, message, error string, status int) {
	c.JSON(http.StatusOK, ResponseError{
		Message: []string{message},
		Error:   error,
		Status:  status,
	})
}