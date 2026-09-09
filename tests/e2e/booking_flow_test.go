//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

type envelope struct {
	Success bool            `json:"success"`
	Code    string          `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
	Meta    *struct {
		Total int `json:"total"`
	} `json:"meta"`
}

type client struct {
	baseURL, token string
	http           *http.Client
	t              *testing.T
}

func (c *client) call(method, path string, body any, expected int, headers map[string]string) envelope {
	c.t.Helper()
	var reader io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			c.t.Fatal(err)
		}
		reader = bytes.NewReader(encoded)
	}
	req, err := http.NewRequest(method, c.baseURL+path, reader)
	if err != nil {
		c.t.Fatal(err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	res, err := c.http.Do(req)
	if err != nil {
		c.t.Fatalf("%s %s: %v", method, path, err)
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		c.t.Fatal(err)
	}
	if res.StatusCode != expected {
		c.t.Fatalf("%s %s: status=%d want=%d body=%s", method, path, res.StatusCode, expected, raw)
	}
	var result envelope
	if err := json.Unmarshal(raw, &result); err != nil {
		c.t.Fatalf("decode response: %v body=%s", err, raw)
	}
	if expected < 400 && !result.Success {
		c.t.Fatalf("%s %s returned unsuccessful response: %s", method, path, raw)
	}
	return result
}

func idFrom(t *testing.T, raw json.RawMessage) int64 {
	t.Helper()
	var item struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(raw, &item); err != nil || item.ID <= 0 {
		t.Fatalf("invalid resource response: %s (%v)", raw, err)
	}
	return item.ID
}

func TestHTTPBookingLifecycle(t *testing.T) {
	baseURL := os.Getenv("E2E_BASE_URL")
	if baseURL == "" {
		baseURL = "http://127.0.0.1:8080"
	}
	c := &client{baseURL: baseURL, http: &http.Client{Timeout: 15 * time.Second}, t: t}
	c.call(http.MethodGet, "/health/ready", nil, http.StatusOK, nil)

	phone := os.Getenv("E2E_ADMIN_PHONE")
	if phone == "" {
		phone = "09120000000"
	}
	sent := c.call(http.MethodPost, "/v1/auth/otp/send", map[string]any{"phone": phone}, http.StatusOK, nil)
	var otp struct {
		Code string `json:"otp_code"`
	}
	if err := json.Unmarshal(sent.Data, &otp); err != nil || otp.Code == "" {
		t.Fatalf("debug OTP missing: %s", sent.Data)
	}
	verified := c.call(http.MethodPost, "/v1/auth/otp/verify", map[string]any{"phone": phone, "code": otp.Code}, http.StatusOK, nil)
	var auth struct {
		AccessToken string `json:"access_token"`
		UserID      int64  `json:"user_id"`
	}
	if err := json.Unmarshal(verified.Data, &auth); err != nil || auth.AccessToken == "" || auth.UserID <= 0 {
		t.Fatalf("invalid auth response: %s", verified.Data)
	}
	c.token = auth.AccessToken
	emailOTPResponse := c.call(http.MethodPost, "/v1/auth/email/otp/send", nil, http.StatusOK, nil)
	var emailOTP struct {
		Code string `json:"otp_code"`
	}
	if err := json.Unmarshal(emailOTPResponse.Data, &emailOTP); err != nil || emailOTP.Code == "" {
		t.Fatalf("debug email OTP missing: %s", emailOTPResponse.Data)
	}
	c.call(http.MethodPost, "/v1/auth/email/otp/verify", map[string]any{"code": emailOTP.Code}, http.StatusOK, nil)

	stamp := time.Now().UnixNano()
	provinceID := idFrom(t, c.call(http.MethodPost, "/v1/provinces", map[string]any{"name": fmt.Sprintf("E2E Province %d", stamp)}, http.StatusCreated, nil).Data)
	cityID := idFrom(t, c.call(http.MethodPost, "/v1/cities", map[string]any{"province_id": provinceID, "name": fmt.Sprintf("E2E City %d", stamp)}, http.StatusCreated, nil).Data)
	movieID := idFrom(t, c.call(http.MethodPost, "/v1/movies", map[string]any{"title": fmt.Sprintf("E2E Movie %d", stamp), "description": "automated booking flow", "duration_minutes": 90, "release_date": time.Now().UTC(), "country": "IR", "status": "NOW_SHOWING", "genre_ids": []int64{}, "people": []any{}}, http.StatusCreated, nil).Data)
	cinemaID := idFrom(t, c.call(http.MethodPost, "/v1/cinemas", map[string]any{"owner_user_id": auth.UserID, "city_id": cityID, "name": fmt.Sprintf("E2E Cinema %d", stamp), "description": "automated", "address": "E2E address", "status": "ACTIVE", "facility_ids": []int64{}}, http.StatusCreated, nil).Data)
	hallID := idFrom(t, c.call(http.MethodPost, fmt.Sprintf("/v1/cinemas/%d/halls", cinemaID), map[string]any{"name": "Hall A", "status": "ACTIVE"}, http.StatusCreated, nil).Data)
	c.call(http.MethodPost, fmt.Sprintf("/v1/halls/%d/seats", hallID), map[string]any{"seats": []map[string]any{{"row_label": "A", "number": 1, "type": "STANDARD", "status": "ACTIVE"}}}, http.StatusCreated, nil)
	startsAt := time.Now().UTC().Add(4 * time.Hour)
	showtimeID := idFrom(t, c.call(http.MethodPost, fmt.Sprintf("/v1/halls/%d/showtimes", hallID), map[string]any{"movie_id": movieID, "starts_at": startsAt, "status": "PUBLISHED", "base_price": 250000, "type_prices": []any{}}, http.StatusCreated, nil).Data)
	seatsResponse := c.call(http.MethodGet, fmt.Sprintf("/v1/showtimes/%d/seats", showtimeID), nil, http.StatusOK, nil)
	var seats []struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(seatsResponse.Data, &seats); err != nil || len(seats) != 1 {
		t.Fatalf("invalid showtime seats: %s", seatsResponse.Data)
	}
	reservationID := idFrom(t, c.call(http.MethodPost, "/v1/reservations", map[string]any{"showtime_id": showtimeID, "showtime_seat_ids": []int64{seats[0].ID}}, http.StatusCreated, nil).Data)
	paymentID := idFrom(t, c.call(http.MethodPost, fmt.Sprintf("/v1/reservations/%d/payments", reservationID), nil, http.StatusCreated, map[string]string{"Idempotency-Key": fmt.Sprintf("e2e-%d", stamp)}).Data)
	c.call(http.MethodPost, fmt.Sprintf("/v1/payments/%d/verify", paymentID), map[string]any{"success": true}, http.StatusOK, nil)

	tickets := c.call(http.MethodGet, "/v1/tickets?page=1&limit=1", nil, http.StatusOK, nil)
	if tickets.Meta == nil || tickets.Meta.Total < 1 {
		t.Fatalf("ticket pagination metadata missing: %+v", tickets.Meta)
	}
	var ticketItems []struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(tickets.Data, &ticketItems); err != nil || len(ticketItems) == 0 {
		t.Fatalf("ticket not issued: %s", tickets.Data)
	}
	deadline := time.Now().Add(20 * time.Second)
	for {
		res := c.call(http.MethodGet, fmt.Sprintf("/v1/tickets/%d/pdf", ticketItems[0].ID), nil, http.StatusOK, nil)
		var pdf struct {
			URL string `json:"url"`
		}
		_ = json.Unmarshal(res.Data, &pdf)
		if pdf.URL != "" {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("ticket PDF was not generated by outbox worker")
		}
		time.Sleep(500 * time.Millisecond)
	}
	c.call(http.MethodGet, fmt.Sprintf("/v1/cinemas/%d/reports/sales", cinemaID), nil, http.StatusOK, nil)
	c.call(http.MethodPost, fmt.Sprintf("/v1/payments/%d/refund", paymentID), nil, http.StatusOK, nil)
	reservation := c.call(http.MethodGet, fmt.Sprintf("/v1/reservations/%d", reservationID), nil, http.StatusOK, nil)
	var reservationState struct {
		Status string `json:"status"`
	}
	_ = json.Unmarshal(reservation.Data, &reservationState)
	if reservationState.Status != "CANCELLED" {
		t.Fatalf("refund did not cancel reservation: %s", reservation.Data)
	}
	audit := c.call(http.MethodGet, "/v1/audit-logs?page=1&limit=10", nil, http.StatusOK, nil)
	if audit.Meta == nil || audit.Meta.Total < 1 {
		t.Fatalf("audit log is empty: %+v", audit.Meta)
	}
}
