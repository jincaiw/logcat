// Package errors 定义业务错误类型，便于上层按类型处理而非字符串匹配。
package errors

import "fmt"

// ErrorType 错误类型
type ErrorType int

const (
	ErrorTypeNotFound ErrorType = iota
	ErrorTypeValidation
	ErrorTypeConflict
	ErrorTypeInternal
	ErrorTypeExternal
)

// BusinessError 业务错误
type BusinessError struct {
	Type    ErrorType
	Message string
	Cause   error
}

func (e *BusinessError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *BusinessError) Unwrap() error {
	return e.Cause
}

// New 创建业务错误
func New(t ErrorType, message string) *BusinessError {
	return &BusinessError{Type: t, Message: message}
}

// Wrap 包装底层错误
func Wrap(t ErrorType, message string, cause error) *BusinessError {
	return &BusinessError{Type: t, Message: message, Cause: cause}
}

// NotFound 资源不存在
func NotFound(format string, args ...interface{}) *BusinessError {
	return &BusinessError{Type: ErrorTypeNotFound, Message: fmt.Sprintf(format, args...)}
}

// Validation 校验失败
func Validation(format string, args ...interface{}) *BusinessError {
	return &BusinessError{Type: ErrorTypeValidation, Message: fmt.Sprintf(format, args...)}
}

// Conflict 资源冲突（如唯一键重复）
func Conflict(format string, args ...interface{}) *BusinessError {
	return &BusinessError{Type: ErrorTypeConflict, Message: fmt.Sprintf(format, args...)}
}

// Internal 内部错误
func Internal(format string, args ...interface{}) *BusinessError {
	return &BusinessError{Type: ErrorTypeInternal, Message: fmt.Sprintf(format, args...)}
}

// External 外部服务错误
func External(format string, args ...interface{}) *BusinessError {
	return &BusinessError{Type: ErrorTypeExternal, Message: fmt.Sprintf(format, args...)}
}
