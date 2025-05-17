 package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 通用响应结构
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// SuccessResponse 返回成功响应
func SuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// CreatedResponse 返回创建成功响应
func CreatedResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// ErrorResponse 返回错误响应
func ErrorResponse(c *gin.Context, statusCode int, message string, errors interface{}) {
	c.JSON(statusCode, Response{
		Status:  "error",
		Message: message,
		Errors:  errors,
	})
}

// BadRequestResponse 返回400错误响应
func BadRequestResponse(c *gin.Context, message string, errors interface{}) {
	ErrorResponse(c, http.StatusBadRequest, message, errors)
}

// UnauthorizedResponse 返回401错误响应
func UnauthorizedResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, message, nil)
}

// ForbiddenResponse 返回403错误响应
func ForbiddenResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, message, nil)
}

// NotFoundResponse 返回404错误响应
func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message, nil)
}

// InternalServerErrorResponse 返回500错误响应
func InternalServerErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, message, nil)
}

// PaginationResponse 分页响应数据
type PaginationResponse struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
}

// GetPagination 获取分页参数
func GetPagination(c *gin.Context) (page int, limit int) {
	page = 1
	limit = 10

	// 从查询参数中获取
	pageParam := c.DefaultQuery("page", "1")
	limitParam := c.DefaultQuery("limit", "10")

	// 尝试转换为整数
	_, err := fmt.Sscanf(pageParam, "%d", &page)
	if err != nil || page < 1 {
		page = 1
	}

	_, err = fmt.Sscanf(limitParam, "%d", &limit)
	if err != nil || limit < 1 {
		limit = 10
	}

	// 限制最大结果数
	if limit > 100 {
		limit = 100
	}

	return page, limit
} 