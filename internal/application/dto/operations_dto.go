package dto

import "time"

type CreateSettlementRequest struct {
	From time.Time `json:"from" binding:"required"`
	To   time.Time `json:"to" binding:"required"`
}

type OperationsFilterRequest struct {
	From        *time.Time
	To          *time.Time
	ActorUserID *int64
	Page        int
	Limit       int
}
