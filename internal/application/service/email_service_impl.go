package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/pkg/logger"
	mail "github.com/mehdi1443/cinema-booking-system/pkg/wneessen-go-mail"
	"github.com/mehdi1443/cinema-booking-system/static"
)

type emailService struct {
	sender          mail.Sender
	debugModeConfig *bootstrap.DebugModeOptions
	statics         *static.StaticFiles
	appLogger       logger.Logger
}

func NewEmailService(
	sender mail.Sender,
	debugModeConfig *bootstrap.DebugModeOptions,
	statics *static.StaticFiles,
	appLogger logger.Logger,
) service_contract.EmailServiceProviderInterface {
	return &emailService{
		sender:          sender,
		debugModeConfig: debugModeConfig,
		statics:         statics,
		appLogger:       appLogger,
	}
}

func (s *emailService) SendEmail(ctx context.Context, req dto.SendEmailRequest) (*bool, error) {
	success := true

	if s.debugModeConfig != nil && s.debugModeConfig.Enable {
		debugText := fmt.Sprintf("[%s] To: %s | Subject: %s\nText: %s\n----------------------------------------\n",
			time.Now().Format(time.RFC3339), req.To, req.Subject, req.TextBody)

		emailFile, err := os.OpenFile("./SEND_EMAIL.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			s.appLogger.Error(logger.IO, logger.Startup, "failed to open email debug file", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			return nil, fmt.Errorf("failed to open email debug log file: %w", err)
		}
		defer emailFile.Close()

		if _, err := emailFile.WriteString(debugText); err != nil {
			s.appLogger.Error(logger.IO, logger.Startup, "failed to write to email debug file", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			return nil, fmt.Errorf("failed to write to email debug log file: %w", err)
		}

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
		s.appLogger.Error(logger.General, logger.ExternalService, "failed to send email via SMTP provider", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, fmt.Errorf("failed to send email: %w", err)
	}

	return &success, nil
}

func (s *emailService) SendOTPEmail(ctx context.Context, req dto.SendOTPEmailRequest) (*bool, error) {
	buffer, err := s.statics.ExecuteOTPTemplate(req.OTPCode, req.OTPExpireMinutes)
	if err != nil {
		s.appLogger.Error(logger.Internal, logger.Api, "failed to execute OTP HTML template", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, fmt.Errorf("failed to render OTP email template: %w", err)
	}

	plainTextBody := fmt.Sprintf(
		"Your account verification code: %s\nThis code is valid for %d minutes.",
		req.OTPCode,
		req.OTPExpireMinutes,
	)

	success, err := s.SendEmail(ctx, dto.SendEmailRequest{
		To:       req.To,
		Subject:  "Verification Code - Cinema Booking System",
		HTMLBody: buffer.String(),
		TextBody: plainTextBody,
	})
	if err != nil {
		return nil, err
	}

	return success, nil
}
