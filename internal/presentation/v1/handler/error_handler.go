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

	// برای Errorهای داخلی هیچ‌وقت err خام را به response نده.
	response.Error(
		c,
		info.StatusCode,
		info.Message,
	)

	return
}

func HandleBadRequest(c *gin.Context, message string) {
	response.Error(
		c,
		http.StatusBadRequest,
		message,
	)
}

func HandleNotFound(c *gin.Context, message string) {
	response.Error(
		c,
		http.StatusNotFound,
		message,
	)
}

func HandleInternal(c *gin.Context) {
	response.Error(
		c,
		http.StatusInternalServerError,
		"Internal server error occurred",
	)
}

func mapError(err error) ErrorInfo {
	switch {

	case errors.Is(err, exception.ErrInvalidInput):
		return ErrorInfo{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid input data or enum value",
		}

	case errors.Is(err, exception.ErrInvalidFilter):
		return ErrorInfo{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid filter parameter",
		}

	case errors.Is(err, exception.ErrInvalidFileType):
		return ErrorInfo{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid file type",
		}

	case errors.Is(err, exception.ErrFileSizeTooLarge):
		return ErrorInfo{
			StatusCode: http.StatusRequestEntityTooLarge,
			Message:    "File size exceeds the allowed limit",
		}

	case errors.Is(err, exception.ErrMovieNotFound):
		return ErrorInfo{
			StatusCode: http.StatusNotFound,
			Message:    "Movie not found",
		}

	case errors.Is(err, exception.ErrGenreNotFound):
		return ErrorInfo{
			StatusCode: http.StatusBadRequest,
			Message:    "One or more specified genres were not found",
		}

	case errors.Is(err, exception.ErrPersonNotFound):
		return ErrorInfo{
			StatusCode: http.StatusBadRequest,
			Message:    "One or more specified persons were not found",
		}

	case errors.Is(err, exception.ErrMovieAlreadyExists):
		return ErrorInfo{
			StatusCode: http.StatusConflict,
			Message:    "Movie already exists",
		}

	case errors.Is(err, exception.ErrDuplicateKey):
		return ErrorInfo{
			StatusCode: http.StatusConflict,
			Message:    "Resource already exists",
		}

	case errors.Is(err, exception.ErrDatabase):
		return ErrorInfo{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal database error occurred",
		}

	case errors.Is(err, exception.ErrNoteAccessible):
		return ErrorInfo{
			StatusCode: http.StatusForbidden,
			Message:    "You do not have permission to perform this action",
		}

	case errors.Is(err, exception.ErrGetTokenFromContext):
		return ErrorInfo{
			StatusCode: http.StatusUnauthorized,
			Message:    "Authentication context is invalid",
		}

	default:
		return ErrorInfo{
			StatusCode: http.StatusInternalServerError,
			Message:    "An unexpected error occurred",
		}
	}
}
