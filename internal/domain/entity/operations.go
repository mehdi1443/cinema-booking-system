package entity

import "time"

type AuditLog struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	ActorUserID *int64    `json:"actor_user_id,omitempty"`
	Method      string    `json:"method"`
	Path        string    `json:"path"`
	Route       *string   `json:"route,omitempty"`
	StatusCode  int       `json:"status_code"`
	IPAddress   *string   `json:"ip_address,omitempty"`
	UserAgent   *string   `json:"user_agent,omitempty"`
	RequestID   *string   `json:"request_id,omitempty"`
	DurationMS  int64     `json:"duration_ms"`
	CreatedAt   time.Time `json:"created_at"`
}

type CinemaSettlement struct {
	ID              int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	CinemaID        int64      `json:"cinema_id"`
	PeriodFrom      time.Time  `json:"period_from"`
	PeriodTo        time.Time  `json:"period_to"`
	GrossAmount     int64      `json:"gross_amount"`
	Status          string     `json:"status"`
	CreatedByUserID int64      `json:"created_by_user_id"`
	PaidByUserID    *int64     `json:"paid_by_user_id,omitempty"`
	PaidAt          *time.Time `json:"paid_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
