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

// TODO: 统一返回类组件

func SuccessResponse(data interface{}) ApiResponse {
	return NewResponse(200, data, "success")
}

func UserErrorResponse(data interface{}, message string) ApiResponse {
	return NewResponse(400, data, message)
}
func UserErrorNoMsgResponse(message string) ApiResponse {
	return NewResponse(400, nil, message)
}

func InternalErrorResponse(data interface{}, message string) ApiResponse {
	return NewResponse(500, data, message)
}

// WriteJSON 将 ApiResponse 写入 HTTP 响应
func WriteJSON(c *gin.Context, response ApiResponse) {
	c.JSON(http.StatusOK, response)
}

func WriteBadRequestJSON(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, UserErrorNoMsgResponse(message))
}
