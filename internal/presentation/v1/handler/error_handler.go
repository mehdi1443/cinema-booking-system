package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

type ErrorInfo struct {
	StatusCode int
	Message    string
}

func Handle(c *gin.Context, err error) {
	if err == nil {
		return
	}

	info := mapError(err)

	response.Error(
		c,
		info.StatusCode,
		info.Message,
	)
}

func HandleBadRequest(c *gin.Context, message string) {
	response.Error(c, http.StatusBadRequest, message)
}

func HandleNotFound(c *gin.Context, message string) {
	response.Error(c, http.StatusNotFound, message)
}

func HandleInternal(c *gin.Context) {
	response.Error(c, http.StatusInternalServerError, "Internal server error occurred")
}

func mapError(err error) ErrorInfo {
	var appErr *exception.AppError

	if errors.As(err, &appErr) {
		status := http.StatusInternalServerError

		switch appErr.Kind {
		case exception.KindValidation:
			if appErr.Code == exception.CodeFileSizeTooLarge {
				status = http.StatusRequestEntityTooLarge
			} else {
				status = http.StatusBadRequest
			}
		case exception.KindNotFound:
			status = http.StatusNotFound
		case exception.KindConflict:
			status = http.StatusConflict
		case exception.KindForbidden:
			status = http.StatusForbidden
		case exception.KindAuth:
			status = http.StatusUnauthorized
		case exception.KindDatabase, exception.KindInternal:
			status = http.StatusInternalServerError
		}

		if status == http.StatusInternalServerError {
			return ErrorInfo{
				StatusCode: status,
				Message:    "An unexpected internal error occurred",
			}
		}

		return ErrorInfo{
			StatusCode: status,
			Message:    appErr.Message,
		}
	}

	return ErrorInfo{
		StatusCode: http.StatusInternalServerError,
		Message:    "An unexpected error occurred",
	}
}
