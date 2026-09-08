package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

func Handle(c *gin.Context, err error, log domainlogger.Logger) {
	if err == nil {
		return
	}

	info := mapError(err)

	ctx := c.Request.Context()
	fields := []domainlogger.Field{
		domainlogger.Error(err),
		domainlogger.Int("status_code", info.StatusCode),
		domainlogger.String("error_code", info.Code),
	}

	if log == nil {
		response.ErrorWithCode(c, info.StatusCode, info.Code, info.Message)
		return
	}

	if info.StatusCode >= 500 {
		log.Error(ctx, "Internal server error occurred", fields...)
	} else if info.StatusCode >= 400 {
		log.Warn(ctx, "Client validation or logic error", fields...)
	}

	response.ErrorWithCode(
		c,
		info.StatusCode,
		info.Code,
		info.Message,
	)
}

func HandleBadRequest(c *gin.Context, message string) {
	response.Error(c, http.StatusBadRequest, message)
}

func HandleNotFound(c *gin.Context, message string) {
	response.Error(c, http.StatusNotFound, message)
}

func HandleInternal(c *gin.Context, log domainlogger.Logger, err error) {
	if log != nil && err != nil {
		log.Error(c.Request.Context(), "Unhandled internal error", domainlogger.Error(err))
	}
	response.Error(c, http.StatusInternalServerError, "Internal server error occurred")
}

type ErrorInfo struct {
	StatusCode int
	Code       string
	Message    string
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
				Code:       string(appErr.Code),
				Message:    "An unexpected internal error occurred",
			}
		}

		return ErrorInfo{
			StatusCode: status,
			Code:       string(appErr.Code),
			Message:    appErr.Message,
		}
	}

	return ErrorInfo{
		StatusCode: http.StatusInternalServerError,
		Code:       string(exception.CodeInternal),
		Message:    "An unexpected error occurred",
	}
}
