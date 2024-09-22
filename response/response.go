// response.go
package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ApiResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"msg"`
}

// NewResponse 创建一个新的 ApiResponse 实例
func NewResponse(code int, data interface{}, message string) ApiResponse {
	return ApiResponse{
		Code:    code,
		Data:    data,
		Message: message,
	}
}

// WriteJSON 将 ApiResponse 写入 HTTP 响应
func WriteJSON(c *gin.Context, response ApiResponse) {
	c.JSON(http.StatusOK, response)
}
