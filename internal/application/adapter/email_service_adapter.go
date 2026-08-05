package service_contract

import "github.com/mehdi1443/cinema-booking-system/internal/application/dto"

type EmailServiceProviderInterface interface {
	SendEmail(req dto.SendEmailRequest) (*bool, error)
}
