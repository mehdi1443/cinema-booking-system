package elastic

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type SearchRepository interface {
	Index(ctx context.Context, doc entity.SearchDocument) error
	Delete(ctx context.Context, id string) error
	Suggest(ctx context.Context, keyword string, limit int) ([]entity.SearchDocument, error)
	Search(ctx context.Context, keyword string, types []enum.DocType, limit int) ([]entity.SearchDocument, int64, error)
}
