package broker

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	domainbroker "github.com/mehdi1443/cinema-booking-system/internal/domain/broker"
)

type MockPaymentGateway struct{}

func NewMockPaymentGateway() domainbroker.PaymentGateway { return &MockPaymentGateway{} }

func (g *MockPaymentGateway) Create(_ context.Context, paymentID, amount int64) (string, string, error) {
	authority := "mock_" + uuid.NewString()
	return authority, "mock://payments/" + authority, nil
}

func (g *MockPaymentGateway) Verify(_ context.Context, _ string, succeed bool) (*domainbroker.PaymentGatewayResult, error) {
	if !succeed {
		return &domainbroker.PaymentGatewayResult{Succeeded: false, FailureReason: "MOCK_PAYMENT_REJECTED"}, nil
	}
	return &domainbroker.PaymentGatewayResult{Succeeded: true, ReferenceID: fmt.Sprintf("mock-ref-%s", uuid.NewString())}, nil
}

func (g *MockPaymentGateway) Refund(_ context.Context, _ string, _ int64) (string, error) {
	return fmt.Sprintf("mock-refund-%s", uuid.NewString()), nil
}
