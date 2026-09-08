package handler

import (
	"net/http"
	"testing"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
)

func TestMapErrorPreservesApplicationCode(t *testing.T) {
	info := mapError(exception.ErrSeatUnavailable)
	if info.StatusCode != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, info.StatusCode)
	}
	if info.Code != string(exception.CodeSeatUnavailable) {
		t.Fatalf("expected code %q, got %q", exception.CodeSeatUnavailable, info.Code)
	}
}

func TestMapErrorHidesInternalMessage(t *testing.T) {
	info := mapError(exception.ErrDatabase)
	if info.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, info.StatusCode)
	}
	if info.Code != string(exception.CodeDatabase) {
		t.Fatalf("expected code %q, got %q", exception.CodeDatabase, info.Code)
	}
	if info.Message == exception.ErrDatabase.Message {
		t.Fatal("internal error message must not be exposed")
	}
}
