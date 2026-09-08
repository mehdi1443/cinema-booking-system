package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestErrorIncludesStableCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ErrorWithCode(ctx, http.StatusConflict, "SEAT_UNAVAILABLE", "seat is unavailable")

	var body StandardResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if recorder.Code != http.StatusConflict || body.Code != "SEAT_UNAVAILABLE" || body.Success {
		t.Fatalf("unexpected error response: %#v", body)
	}
}

func TestErrorUsesDefaultStatusAndCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	Error(ctx, 0, "invalid request")

	var body StandardResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if recorder.Code != http.StatusBadRequest || body.Code != "INVALID_INPUT" {
		t.Fatalf("unexpected default error response: %#v", body)
	}
}
