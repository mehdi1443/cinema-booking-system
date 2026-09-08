package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StandardResponse struct {
	Success    bool     `json:"success"`
	StatusCode int      `json:"statusCode,omitempty"`
	Code       string   `json:"code,omitempty"`
	Message    string   `json:"message"`
	Data       any      `json:"data,omitempty"`
	Meta       *Meta    `json:"meta,omitempty"`
	Errors     []string `json:"errors,omitempty"`
}

type Meta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

func (sr *StandardResponse) SendResponse(c *gin.Context) {
	c.JSON(sr.StatusCode, sr)
}

func Error(c *gin.Context, statusCode int, message string, errorsList ...error) {
	ErrorWithCode(c, statusCode, defaultErrorCode(statusCode), message, errorsList...)
}

func ErrorWithCode(c *gin.Context, statusCode int, code, message string, errorsList ...error) {
	if statusCode == 0 {
		statusCode = http.StatusBadRequest
	}
	res := &StandardResponse{
		Success:    false,
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Data:       nil,
		Errors:     make([]string, 0),
	}

	for _, err := range errorsList {
		if err != nil {
			res.Errors = append(res.Errors, err.Error())
		}
	}

	res.SendResponse(c)
}

func defaultErrorCode(statusCode int) string {
	switch statusCode {
	case http.StatusBadRequest:
		return "INVALID_INPUT"
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusConflict:
		return "CONFLICT"
	case http.StatusRequestEntityTooLarge:
		return "FILE_SIZE_TOO_LARGE"
	case http.StatusTooManyRequests:
		return "RATE_LIMIT_EXCEEDED"
	default:
		if statusCode >= http.StatusInternalServerError {
			return "INTERNAL_ERROR"
		}
		return "REQUEST_FAILED"
	}
}

func Success(c *gin.Context, statusCode int, message string, data any) {
	res := &StandardResponse{
		Success:    true,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}

	if res.StatusCode == 0 {
		res.StatusCode = http.StatusOK
	}

	res.SendResponse(c)
}

func SuccessWithMeta(c *gin.Context, statusCode int, message string, data any, meta Meta) {
	res := &StandardResponse{
		Success:    true,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Meta:       &meta,
	}

	if res.StatusCode == 0 {
		res.StatusCode = http.StatusOK
	}

	res.SendResponse(c)
}
