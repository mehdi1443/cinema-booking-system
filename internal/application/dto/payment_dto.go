package dto

import (
	"fmt"
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/application/helper"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type VerifyMockPaymentRequest struct {
	Success       *bool  `json:"success" binding:"required"`
	FailureReason string `json:"failure_reason" binding:"omitempty,max=255"`
}

type PaymentResponse struct {
	ID                int64                `json:"id"`
	ReservationID     int64                `json:"reservation_id"`
	UserID            int64                `json:"user_id"`
	Amount            int64                `json:"amount"`
	Provider          enum.PaymentProvider `json:"provider"`
	Status            enum.PaymentStatus   `json:"status"`
	Authority         string               `json:"authority"`
	PaymentURL        string               `json:"payment_url,omitempty"`
	VerifyEndpoint    string               `json:"verify_endpoint,omitempty"`
	ReferenceID       *string              `json:"reference_id,omitempty"`
	FailureReason     *string              `json:"failure_reason,omitempty"`
	PaidAt            *time.Time           `json:"paid_at,omitempty"`
	RefundReferenceID *string              `json:"refund_reference_id,omitempty"`
	RefundedAt        *time.Time           `json:"refunded_at,omitempty"`
	CreatedAt         *time.Time           `json:"created_at,omitempty"`
	UpdatedAt         *time.Time           `json:"updated_at,omitempty"`
}

func ToPaymentResponse(payment *entity.Payment) *PaymentResponse {
	if payment == nil {
		return nil
	}
	result := &PaymentResponse{
		ID: payment.ID, ReservationID: payment.ReservationID, UserID: payment.UserID,
		Amount: payment.Amount, Provider: payment.Provider, Status: payment.Status,
		Authority: payment.Authority, ReferenceID: payment.ReferenceID,
		FailureReason: payment.FailureReason, PaidAt: payment.PaidAt,
		RefundReferenceID: payment.RefundReferenceID, RefundedAt: payment.RefundedAt,
		CreatedAt: helper.TimeToPtr(payment.CreatedAt), UpdatedAt: helper.TimeToPtr(payment.UpdatedAt),
	}
	if payment.Status == enum.PaymentPending {
		result.PaymentURL = "mock://payments/" + payment.Authority
		result.VerifyEndpoint = fmt.Sprintf("/v1/payments/%d/verify", payment.ID)
	}
	return result
}
