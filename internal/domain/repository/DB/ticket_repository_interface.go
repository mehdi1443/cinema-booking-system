package DB

import (
	"context"
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
)

type TicketRepositoryInterface interface {
	GetTicket(ctx context.Context, id int64) (*entity.Ticket, error)
	ListUserTickets(ctx context.Context, userID int64) ([]entity.Ticket, error)
	ConsumeTicket(ctx context.Context, code string, validatorUserID int64, allowAny bool, now time.Time) (*entity.Ticket, error)
}
