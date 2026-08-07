package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	mail "github.com/mehdi1443/cinema-booking-system/pkg/wneessen-go-mail"
)

type emailService struct {
	emailConfig       *bootstrap.EmailConfig
	DEBUG_MODE_CONFIG *bootstrap.DebugModeOptions
	sender            mail.Sender
}

func NewEmailService(
	emailConfig *bootstrap.EmailConfig,
	DEBUG_MODE_CONFIG *bootstrap.DebugModeOptions,
) service_contract.EmailServiceProviderInterface {
	return &emailService{
		emailConfig:       emailConfig,
		DEBUG_MODE_CONFIG: DEBUG_MODE_CONFIG,
	}
}

func (s *emailService) SendEmail(ctx context.Context, req dto.SendEmailRequest) (*bool, error) {
	var success bool = false
	if s.DEBUG_MODE_CONFIG.Enable { // write to SEND_EMAIL.txt file
		Text := fmt.Sprintf("[%s] Email:%s - Text:%s\n", time.Now().String(), req.To, req.TextBody)

		emailFile, err := os.OpenFile("./SEND_EMAIL.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Println("Fail:Send Email:%w", err)
			return nil, err
		}
		defer emailFile.Close()
		_, err = emailFile.WriteString(Text)
		if err != nil {
			return nil, err
		}

		success = true
		return &success, nil
	}
	err := s.sender.Send(
		ctx,
		mail.SendOptions{
			To:          req.To,
			Subject:     req.Subject,
			PlainBody:   req.TextBody,
			HTMLBody:    req.HTMLBody,
			EmbedFiles:  req.EmbedFiles,
			Attachments: req.Attachments,
		},
	)
	if err != nil {
		return nil, err
	}
	success = true
	return &success, nil
}
