package service

import (
	"fmt"
	"log"
	"os"
	"time"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
)

type emailService struct {
}

func NewEmailService() service_contract.EmailServiceProviderInterface {
	return &emailService{}
}

//TODO: Send OTP via SMS provider logic here...

func (s *emailService) SendEmail(req dto.SendEmailRequest) (*bool, error) {
	Text := fmt.Sprintf("[%s] Email:%s - Text:%s\n", time.Now().String(), req.Email, req.Text)

	ffile, err := os.OpenFile("./SEND_EMAIL.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println("Fail:Send Email:%w", err)
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
