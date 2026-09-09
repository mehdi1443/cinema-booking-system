package service

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	qrcode "github.com/skip2/go-qrcode"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	storagerepo "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/storage"
	"github.com/mehdi1443/cinema-booking-system/pkg/ticketpdf"
)

const ticketQRPrefix = "cinema-ticket:"
const ticketPDFURLLifetime = 15 * time.Minute

type ticketService struct {
	repo    DB.TicketRepositoryInterface
	storage storagerepo.StorageInterFace
	log     domainlogger.Logger
}

func NewTicketService(repo DB.TicketRepositoryInterface, storage storagerepo.StorageInterFace, log domainlogger.Logger) service_contract.TicketServiceInterface {
	return &ticketService{repo: repo, storage: storage, log: log}
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
	png, err := generateTicketQRCode(ticket.Code)
	if err != nil {
		return nil, exception.Wrap(exception.ErrInternal, "failed to generate ticket QR code", err)
	}
	return png, nil
}

func generateTicketQRCode(code string) ([]byte, error) {
	return qrcode.Encode(ticketQRPrefix+code, qrcode.Medium, 256)
}

func (s *ticketService) PDF(ctx context.Context, id, userID int64, allowAny bool) (*dto.TicketPDFResponse, error) {
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
	if ticket.PDFObjectKey == nil || strings.TrimSpace(*ticket.PDFObjectKey) == "" {
		if err := s.ensurePDF(ctx, ticket); err != nil {
			return nil, err
		}
		ticket, err = s.repo.GetTicket(ctx, id)
		if err != nil {
			return nil, err
		}
	}
	expiresAt := time.Now().Add(ticketPDFURLLifetime)
	url, err := s.storage.GetPresignedURL(ctx, enum.TicketPDF, *ticket.PDFObjectKey, ticketPDFURLLifetime)
	if err != nil {
		return nil, exception.Wrap(exception.ErrInternal, "failed to create ticket PDF URL", err)
	}
	return &dto.TicketPDFResponse{URL: url, ExpiresAt: expiresAt}, nil
}

func (s *ticketService) EnsurePDFForPayment(ctx context.Context, paymentID int64) error {
	if paymentID <= 0 {
		return exception.ErrInvalidInput
	}
	ticket, err := s.repo.GetTicketByPaymentID(ctx, paymentID)
	if err != nil {
		return err
	}
	return s.ensurePDF(ctx, ticket)
}

func (s *ticketService) ensurePDF(ctx context.Context, ticket *entity.Ticket) error {
	if ticket == nil {
		return exception.ErrTicketNotFound
	}
	if ticket.PDFURL != nil && ticket.PDFObjectKey != nil && ticket.PDFGeneratedAt != nil {
		return nil
	}
	reservation := ticket.Reservation
	if reservation == nil || reservation.Showtime == nil {
		return exception.ErrDatabase
	}
	showtime := reservation.Showtime
	data := ticketpdf.Data{
		TicketCode: ticket.Code, TicketID: ticket.ID, ReservationID: ticket.ReservationID,
		PaymentID: ticket.PaymentID, StartsAt: showtime.StartsAt, EndsAt: showtime.EndsAt,
	}
	if showtime.Movie != nil {
		data.MovieTitle = showtime.Movie.Title
	}
	if showtime.Hall != nil {
		data.HallName = showtime.Hall.Name
		if showtime.Hall.Cinema != nil {
			data.CinemaName = showtime.Hall.Cinema.Name
		}
	}
	data.Seats = make([]ticketpdf.Seat, 0, len(reservation.Items))
	for _, item := range reservation.Items {
		if item.ShowtimeSeat != nil {
			data.Seats = append(data.Seats, ticketpdf.Seat{
				Row: item.ShowtimeSeat.RowLabel, Number: item.ShowtimeSeat.SeatNumber,
				Type: string(item.ShowtimeSeat.SeatType),
			})
		}
	}
	pdf, err := ticketpdf.Generate(data)
	if err != nil {
		return exception.Wrap(exception.ErrInternal, "failed to generate ticket PDF", err)
	}
	objectKey := fmt.Sprintf("tickets/user-%d/ticket-%d.pdf", ticket.UserID, ticket.ID)
	url, err := s.storage.UploadObject(ctx, enum.TicketPDF, objectKey, bytes.NewReader(pdf), "application/pdf")
	if err != nil {
		return exception.Wrap(exception.ErrInternal, "failed to upload ticket PDF", err)
	}
	generatedAt := time.Now()
	if _, err := s.repo.SetTicketPDF(ctx, ticket.ID, url, objectKey, generatedAt); err != nil {
		return err
	}
	s.log.Info(ctx, "Ticket PDF generated successfully",
		domainlogger.Int64("ticket_id", ticket.ID),
		domainlogger.Int64("payment_id", ticket.PaymentID),
	)
	return nil
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
