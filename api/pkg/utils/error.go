package utils

import (
	"fmt"
	"net/http"

	"github.com/MAJIAXIT/projname/api/pkg/logger"
	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewConflict(format string, v ...any) *AppError {
	return &AppError{
		Code:    http.StatusConflict,
		Message: fmt.Sprintf(format, v...),
	}
}

func NewBadRequest(format string, v ...any) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf(format, v...),
	}
}

func NewNotFound(format string, v ...any) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf(format, v...),
	}
}

func NewInternal(format string, v ...any) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: fmt.Sprintf(format, v...),
	}
}

func NewUnauthorized(format string, v ...any) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: fmt.Sprintf(format, v...),
	}
}

func NewTooManyRequests(format string, v ...any) *AppError {
	return &AppError{
		Code:    http.StatusTooManyRequests,
		Message: fmt.Sprintf(format, v...),
	}
}

func HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		c.JSON(appErr.Code, gin.H{"error": appErr.Message})
	} else {
		if err != nil {
			logger.Error(err.Error())
		}
		c.Status(http.StatusInternalServerError)
	}
}

func WrapError(err error) error {
	if _, ok := err.(*AppError); ok {
		return err
	} else {
		return logger.WrapError(err)
	}
}
