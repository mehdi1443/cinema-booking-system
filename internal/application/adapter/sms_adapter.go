package service_contract

import "github.com/mehdi1443/cinema-booking-system/internal/application/dto"

type SMSServiceProviderInterface interface {
	SendSMS(req dto.SendSMSRequest) (*bool, error)
}
