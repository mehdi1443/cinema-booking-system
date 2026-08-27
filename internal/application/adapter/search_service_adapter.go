package service_contract

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
)

type SearchServiceInterface interface {
	Suggest(ctx context.Context, req dto.SuggestRequest) ([]dto.SearchResponseItem, error)
	Search(ctx context.Context, req dto.SearchQueryRequest) ([]dto.SearchResponseItem, int64, error)
}
