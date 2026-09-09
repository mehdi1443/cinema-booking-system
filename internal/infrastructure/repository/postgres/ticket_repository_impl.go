package postgres_repo

import (
	"context"
	"errors"
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ticketRepository struct {
	db database.Database
}

func NewTicketRepository(db database.Database) DB.TicketRepositoryInterface {
	return &ticketRepository{db: db}
}

func ticketDetails(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Reservation.Showtime.Hall.Cinema").
		Preload("Reservation.Showtime.Movie").
		Preload("Reservation.Items.ShowtimeSeat")
}

func (r *ticketRepository) GetTicket(ctx context.Context, id int64) (*entity.Ticket, error) {
	var ticket entity.Ticket
	err := ticketDetails(r.db.GetDB().WithContext(ctx)).First(&ticket, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrTicketNotFound
	}
	if err != nil {
		return nil, handlePaymentDBError(err)
	}
	return &ticket, nil
}

func (r *ticketRepository) GetTicketByPaymentID(ctx context.Context, paymentID int64) (*entity.Ticket, error) {
	var ticket entity.Ticket
	err := ticketDetails(r.db.GetDB().WithContext(ctx)).Where("tickets.payment_id = ?", paymentID).First(&ticket).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrTicketNotFound
	}
	if err != nil {
		return nil, handlePaymentDBError(err)
	}
	return &ticket, nil
}

func (r *ticketRepository) ListUserTickets(ctx context.Context, userID int64, limit, offset int) ([]entity.Ticket, int64, error) {
	var tickets []entity.Ticket
	query := r.db.GetDB().WithContext(ctx).Model(&entity.Ticket{}).Where("tickets.user_id = ?", userID)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, handlePaymentDBError(err)
	}
	err := ticketDetails(query).Order("tickets.created_at DESC").Limit(limit).Offset(offset).Find(&tickets).Error
	if err != nil {
		return nil, 0, handlePaymentDBError(err)
	}
	return tickets, total, nil
}

func (r *ticketRepository) SetTicketPDF(ctx context.Context, id int64, url, objectKey string, generatedAt time.Time) (*entity.Ticket, error) {
	result := r.db.GetDB().WithContext(ctx).Model(&entity.Ticket{}).Where("id = ?", id).Updates(map[string]any{
		"pdf_url": url, "pdf_object_key": objectKey, "pdf_generated_at": generatedAt, "updated_at": generatedAt,
	})
	if result.Error != nil {
		return nil, handlePaymentDBError(result.Error)
	}
	if result.RowsAffected != 1 {
		return nil, exception.ErrTicketNotFound
	}
	return r.GetTicket(ctx, id)
}

func (r *ticketRepository) ConsumeTicket(ctx context.Context, code string, validatorUserID int64, allowAny bool, now time.Time) (*entity.Ticket, error) {
	var ticket entity.Ticket
	var validationErr error
	err := r.db.WithTransaction(func(txDB database.Database) error {
		tx := txDB.GetDB().WithContext(ctx)
		query := ticketDetails(tx.Clauses(clause.Locking{Strength: "UPDATE"}))
		if err := query.Where("tickets.code = ?", code).First(&ticket).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrTicketNotFound
			}
			return handlePaymentDBError(err)
		}

		if ticket.Reservation == nil {
			return exception.ErrDatabase
		}
		showtime := ticket.Reservation.Showtime
		if showtime == nil || showtime.Hall == nil || showtime.Hall.Cinema == nil {
			return exception.ErrDatabase
		}
		if !allowAny {
			ownerID := showtime.Hall.Cinema.OwnerUserID
			if ownerID == nil || *ownerID != validatorUserID {
				return exception.ErrNotAccessible
			}
		}

		switch ticket.Status {
		case enum.TicketUsed:
			validationErr = exception.ErrTicketAlreadyUsed
			return nil
		case enum.TicketExpired:
			validationErr = exception.ErrTicketExpired
			return nil
		case enum.TicketCancelled:
			validationErr = exception.ErrTicketInvalid
			return nil
		case enum.TicketActive:
		default:
			validationErr = exception.ErrTicketInvalid
			return nil
		}
		if showtime.Status == enum.ShowtimeStatusCancelled {
			if err := tx.Model(&ticket).Updates(map[string]any{"status": enum.TicketCancelled, "updated_at": now}).Error; err != nil {
				return handlePaymentDBError(err)
			}
			ticket.Status = enum.TicketCancelled
			ticket.UpdatedAt = now
			validationErr = exception.ErrTicketInvalid
			return nil
		}
		if showtime.Status != enum.ShowtimeStatusPublished {
			validationErr = exception.ErrTicketInvalid
			return nil
		}

		if now.Before(showtime.StartsAt.Add(-time.Hour)) {
			validationErr = exception.ErrTicketNotYetValid
			return nil
		}
		if !now.Before(showtime.EndsAt) {
			if err := tx.Model(&ticket).Updates(map[string]any{"status": enum.TicketExpired, "updated_at": now}).Error; err != nil {
				return handlePaymentDBError(err)
			}
			ticket.Status = enum.TicketExpired
			ticket.UpdatedAt = now
			validationErr = exception.ErrTicketExpired
			return nil
		}

		result := tx.Model(&entity.Ticket{}).
			Where("id = ? AND status = ?", ticket.ID, enum.TicketActive).
			Updates(map[string]any{"status": enum.TicketUsed, "used_at": now, "used_by_user_id": validatorUserID, "updated_at": now})
		if result.Error != nil {
			return handlePaymentDBError(result.Error)
		}
		if result.RowsAffected != 1 {
			return exception.ErrTicketAlreadyUsed
		}
		ticket.Status = enum.TicketUsed
		ticket.UsedAt = &now
		ticket.UsedByUserID = &validatorUserID
		ticket.UpdatedAt = now
		return nil
	})
	if err != nil {
		return nil, err
	}
	if validationErr != nil {
		return nil, validationErr
	}
	return &ticket, nil
}
