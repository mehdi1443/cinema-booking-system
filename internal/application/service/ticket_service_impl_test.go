package service

import (
	"bytes"
	"image/png"
	"testing"
)

func TestGenerateTicketQRCode(t *testing.T) {
	t.Parallel()
	data, err := generateTicketQRCode("ct_test_ticket_code")
	if err != nil {
		t.Fatalf("generate QR code: %v", err)
	}
	if len(data) < 8 || !bytes.Equal(data[:8], []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a}) {
		t.Fatal("QR response is not a PNG")
	}
	image, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("decode QR PNG: %v", err)
	}
	if image.Bounds().Dx() != 256 || image.Bounds().Dy() != 256 {
		t.Fatalf("unexpected QR dimensions: %v", image.Bounds())
	}
}
