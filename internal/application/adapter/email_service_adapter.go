package service_contract

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
)

type EmailServiceProviderInterface interface {
	SendEmail(ctx context.Context, req dto.SendEmailRequest) (*bool, error)
}
