package service

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	qrcode "github.com/skip2/go-qrcode"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
)

const ticketQRPrefix = "cinema-ticket:"

type ticketService struct {
	repo DB.TicketRepositoryInterface
	log  domainlogger.Logger
}

func NewTicketService(repo DB.TicketRepositoryInterface, log domainlogger.Logger) service_contract.TicketServiceInterface {
	return &ticketService{repo: repo, log: log}
}

func (s *ticketService) Get(ctx context.Context, id, userID int64, allowAny bool) (*dto.TicketResponse, error) {
	if id <= 0 || userID <= 0 {
		return nil, exception.ErrInvalidInput
	}
	ticket, err := s.repo.GetTicket(ctx, id)
	if err != nil {
		return nil, err
	}
	if !allowAny && ticket.UserID != userID {
		return nil, exception.ErrNotAccessible
	}
	return dto.ToTicketResponse(ticket), nil
}

func (s *ticketService) ListMine(ctx context.Context, userID int64) ([]*dto.TicketResponse, error) {
	if userID <= 0 {
		return nil, exception.ErrInvalidInput
	}
	tickets, err := s.repo.ListUserTickets(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]*dto.TicketResponse, 0, len(tickets))
	for i := range tickets {
		result = append(result, dto.ToTicketResponse(&tickets[i]))
	}
	return result, nil
}

func (s *ticketService) QRCode(ctx context.Context, id, userID int64, allowAny bool) ([]byte, error) {
	ticket, err := s.repo.GetTicket(ctx, id)
	if err != nil {
		return nil, err
	}
	if !allowAny && ticket.UserID != userID {
		return nil, exception.ErrNotAccessible
	}
	png, err := generateTicketQRCode(ticket.Code)
	if err != nil {
		return nil, exception.Wrap(exception.ErrInternal, "failed to generate ticket QR code", err)
	}
	return png, nil
}

func generateTicketQRCode(code string) ([]byte, error) {
	return qrcode.Encode(ticketQRPrefix+code, qrcode.Medium, 256)
}

func (s *ticketService) Validate(ctx context.Context, code string, validatorUserID int64, allowAny bool) (*dto.TicketResponse, error) {
	code = strings.TrimSpace(code)
	code = strings.TrimPrefix(code, ticketQRPrefix)
	if code == "" || utf8.RuneCountInString(code) > 100 || validatorUserID <= 0 {
		return nil, exception.ErrInvalidInput
	}
	ticket, err := s.repo.ConsumeTicket(ctx, code, validatorUserID, allowAny, time.Now())
	if err != nil {
		return nil, err
	}
	s.log.Info(ctx, "Ticket validated successfully",
		domainlogger.Int64("ticket_id", ticket.ID),
		domainlogger.Int64("reservation_id", ticket.ReservationID),
		domainlogger.Int64("validator_user_id", validatorUserID),
	)
	return dto.ToTicketResponse(ticket), nil
}
