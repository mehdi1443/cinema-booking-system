package entity

import (
	"encoding/json"
	"time"
)

const (
	OutboxTicketPDFGenerate     = "ticket.pdf.generate"
	OutboxEmailPaymentSucceeded = "email.payment.succeeded"
	OutboxEmailPaymentRefunded  = "email.payment.refunded"
)

type OutboxEvent struct {
	ID            int64           `gorm:"primaryKey;autoIncrement"`
	EventType     string          `gorm:"size:100;not null"`
	AggregateType string          `gorm:"size:50;not null"`
	AggregateID   int64           `gorm:"not null"`
	Payload       json.RawMessage `gorm:"type:jsonb;not null"`
	Status        string          `gorm:"size:20;not null"`
	Attempts      int             `gorm:"not null"`
	AvailableAt   time.Time       `gorm:"not null"`
	LockedAt      *time.Time
	ProcessedAt   *time.Time
	LastError     *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (OutboxEvent) TableName() string { return "outbox_events" }
