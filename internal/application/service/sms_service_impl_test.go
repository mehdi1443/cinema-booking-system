package service

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
)

func TestMockSMSWritesOTPToConfiguredFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sms.log")
	provider := NewMockSMSService(path)
	sent, err := provider.SendSMS(dto.SendSMSRequest{Phone: "09123456789", Text: "OTP:123456"})
	if err != nil || sent == nil || !*sent {
		t.Fatalf("send mock SMS: sent=%v err=%v", sent, err)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "Phone:09123456789") || !strings.Contains(string(content), "OTP:123456") {
		t.Fatalf("unexpected mock SMS content: %q", content)
	}
}
