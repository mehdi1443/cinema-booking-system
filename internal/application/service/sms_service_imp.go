package service

import (
	"fmt"
	"log"
	"os"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
)

type smsService struct {
}

func NewSMSService() service_contract.SMSServiceProviderInterface {
	return &smsService{}
}

//TODO: Send OTP via SMS provider logic here...

func (s *smsService) SendSMS(req dto.SendSMSRequest) (*bool, error) {
	Text := fmt.Sprintf("Phone:%s - Text:%s\n", req.Phone, req.Text)
	err := os.WriteFile("./SEND_SMS", []byte(Text), os.FileMode(os.O_CREATE))
	if err != nil {
		log.Println("Fail:Send SMS:%w", err)
		return nil, err
	}
	success := true
	return &success, nil
}
