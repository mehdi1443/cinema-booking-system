package service

import (
	"fmt"
	"log"
	"os"
	"time"

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
	Text := fmt.Sprintf("[%s] Phone:%s - Text:%s\n", time.Now().String(), req.Phone, req.Text)

	ffile, err := os.OpenFile("./SEND_SMS.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println("Fail:Send SMS:%w", err)
		return nil, err
	}
	defer ffile.Close()
	_, err = ffile.WriteString(Text)
	if err != nil {
		return nil, err
	}

	success := true
	return &success, nil
}
