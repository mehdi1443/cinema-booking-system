package ticketpdf

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func sampleData() Data {
	start := time.Date(2026, time.September, 12, 19, 30, 0, 0, time.FixedZone("IRST", 3*60*60+30*60))
	return Data{
		TicketCode: "ct_76f52e783bde4f69a570c08c966a7337", TicketID: 42,
		ReservationID: 1204, PaymentID: 9881, MovieTitle: "آخرین اکران",
		CinemaName: "سینما آزادی", HallName: "سالن ۳", StartsAt: start, EndsAt: start.Add(2*time.Hour + 10*time.Minute),
		Seats: []Seat{{Row: "B", Number: 7, Type: "VIP"}, {Row: "B", Number: 8, Type: "VIP"}},
	}
}

func TestGenerate(t *testing.T) {
	pdf, err := Generate(sampleData())
	if err != nil {
		t.Fatal(err)
	}
	if len(pdf) < 1_000 || !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		t.Fatalf("invalid PDF output: %d bytes", len(pdf))
	}

	if output := os.Getenv("TICKET_PDF_SAMPLE_PATH"); output != "" {
		if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(output, pdf, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}
