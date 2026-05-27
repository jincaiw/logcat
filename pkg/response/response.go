package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Response unified API response structure
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	RequestID string      `json:"requestId"`
}

// PageResponse paginated response structure
type PageResponse struct {
	Items    interface{} `json:"items"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

// ErrorResponse error response structure
type ErrorResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"requestId"`
}

// getRequestID extracts or generates a request ID from the context
func getRequestID(c *gin.Context) string {
	if rid, ok := c.Get("RequestID"); ok {
		if s, ok := rid.(string); ok {
			return s
		}
	}
	return uuid.New().String()
}

// Success responds with a successful result
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      200,
		Message:   "success",
		Data:      data,
		RequestID: getRequestID(c),
	})
}

// SuccessWithPage responds with paginated data
func SuccessWithPage(c *gin.Context, items interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:      200,
		Message:   "success",
		Data: PageResponse{
			Items:    items,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
		RequestID: getRequestID(c),
	})
}

// SuccessWithMessage responds with a custom success message
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      200,
		Message:   message,
		Data:      data,
		RequestID: getRequestID(c),
	})
}

// Created responds with 201 status
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:      201,
		Message:   "created",
		Data:      data,
		RequestID: getRequestID(c),
	})
}

// Error responds with an error
func Error(c *gin.Context, httpStatus int, code int, message string) {
	c.AbortWithStatusJSON(httpStatus, ErrorResponse{
		Code:      code,
		Message:   message,
		RequestID: getRequestID(c),
	})
}

// BadRequest responds with 400
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, http.StatusBadRequest, message)
}

// Unauthorized responds with 401
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "unauthorized"
	}
	Error(c, http.StatusUnauthorized, http.StatusUnauthorized, message)
}

// Forbidden responds with 403
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "forbidden"
	}
	Error(c, http.StatusForbidden, http.StatusForbidden, message)
}

// NotFound responds with 404
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "not found"
	}
	Error(c, http.StatusNotFound, http.StatusNotFound, message)
}

// Conflict responds with 409
func Conflict(c *gin.Context, message string) {
	if message == "" {
		message = "conflict"
	}
	Error(c, http.StatusConflict, http.StatusConflict, message)
}

// InternalError responds with 500
func InternalError(c *gin.Context, message string) {
	if message == "" {
		message = "internal server error"
	}
	Error(c, http.StatusInternalServerError, http.StatusInternalServerError, message)
}