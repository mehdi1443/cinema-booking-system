package broker

import "context"

type PaymentGatewayResult struct {
	Succeeded     bool
	ReferenceID   string
	FailureReason string
}

type PaymentGateway interface {
	Create(ctx context.Context, paymentID, amount int64) (authority, paymentURL string, err error)
	Verify(ctx context.Context, authority string, succeed bool) (*PaymentGatewayResult, error)
}
