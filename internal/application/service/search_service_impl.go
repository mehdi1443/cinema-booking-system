package service

import (
	"context"
	"fmt"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/elastic"
)

type searchService struct {
	searchRepo elastic.SearchRepository
	log        domainlogger.Logger
}

func NewSearchService(
	searchRepo elastic.SearchRepository,
	log domainlogger.Logger,
) service_contract.SearchServiceInterface {
	return &searchService{
		searchRepo: searchRepo,
		log:        log,
	}
}

func (s *searchService) Suggest(
	ctx context.Context,
	req dto.SuggestRequest,
) ([]dto.SearchResponseItem, error) {

	if req.Keyword == "" {
		return nil, exception.ErrInvalidInput
	}

	docs, err := s.searchRepo.Suggest(
		ctx,
		req.Keyword,
		req.Limit,
	)

	if err != nil {
		s.log.Error(ctx, "Failed to execute search suggestion", domainlogger.Error(err), domainlogger.String("keyword", req.Keyword))
		return nil, fmt.Errorf("suggestion failed: %w", err)
	}

	return dto.ToSearchResponseSlice(docs), nil
}

func (s *searchService) Search(
	ctx context.Context,
	req dto.SearchQueryRequest,
) ([]dto.SearchResponseItem, int64, error) {

	if req.Keyword == "" {
		return nil, 0, exception.ErrInvalidInput
	}

	types := req.GetDocTypes()

	docs, total, err := s.searchRepo.Search(
		ctx,
		req.Keyword,
		types,
		req.Limit,
	)

	if err != nil {
		s.log.Error(ctx, "Failed to execute full search", domainlogger.Error(err), domainlogger.String("keyword", req.Keyword))
		return nil, 0, fmt.Errorf("search failed: %w", err)
	}

	return dto.ToSearchResponseSlice(docs), total, nil
}
