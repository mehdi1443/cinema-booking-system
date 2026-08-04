package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StandardResponse struct {
	Success    bool     `json:"success"`
	StatusCode int      `json:"statusCode,omitempty"`
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
	res := &StandardResponse{
		Success:    false,
		StatusCode: statusCode,
		Message:    message,
		Data:       nil,
		Errors:     make([]string, 0),
	}

	for _, err := range errorsList {
		if err != nil {
			res.Errors = append(res.Errors, err.Error())
		}
	}

	if res.StatusCode == 0 {
		res.StatusCode = http.StatusBadRequest
	}

	res.SendResponse(c)
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
