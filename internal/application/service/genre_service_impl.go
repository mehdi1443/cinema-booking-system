package service

import (
	"context"
	"strings"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
)

type genreService struct {
	genreRepo repository_contract.GenreInterface
}

func NewGenreService(genreRepo repository_contract.GenreInterface) service_contract.GenreServiceInterface {
	return &genreService{
		genreRepo: genreRepo,
	}
}

func (s *genreService) CreateGenre(ctx context.Context, req *dto.CreateGenreRequest) (*dto.GenreResponse, error) {
	req.Sanitize()
	if len(req.Name) < 2 || len(req.Name) > 50 {
		return nil, exception.ErrInvalidInput
	}

	genre := &entity.Genre{
		Name: req.Name,
	}

	createdGenre, err := s.genreRepo.Create(ctx, genre)
	if err != nil {
		return nil, err
	}

	return dto.ToGenreResponse(createdGenre), nil
}

func (s *genreService) GetGenreByID(ctx context.Context, id int64) (*dto.GenreResponse, error) {
	genre, err := s.genreRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.ToGenreResponse(genre), nil
}

func (s *genreService) GetGenreByName(ctx context.Context, name string) (*dto.GenreResponse, error) {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return nil, exception.ErrInvalidInput
	}

	genre, err := s.genreRepo.GetByName(ctx, trimmedName)
	if err != nil {
		return nil, err
	}

	return dto.ToGenreResponse(genre), nil
}

func (s *genreService) GetGenres(ctx context.Context, req dto.GenreFilterRequest) ([]*dto.GenreResponse, int64, error) {
	filter := repository_contract.GenreFilter{}
	if req.Search != nil {
		trimmedSearch := strings.TrimSpace(*req.Search)
		if trimmedSearch != "" {
			filter.SearchKeyword = trimmedSearch
		}
	}

	opts := repository_contract.NewQueryOptions()

	const (
		maxLimit = 100
		maxPage  = 1000
	)

	page := req.Page
	if page < 1 {
		page = 1
	} else if page > maxPage {
		page = maxPage
	}

	limit := req.Limit
	if limit < 1 {
		limit = 10
	} else if limit > maxLimit {
		limit = maxLimit
	}

	offset := (page - 1) * limit
	opts.WithPagination(limit, offset)

	genres, total, err := s.genreRepo.GetGenres(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}

	return dto.ToGenreResponseSlice(genres), total, nil
}

func (s *genreService) UpdateGenre(ctx context.Context, id int64, req *dto.UpdateGenreRequest) (*dto.GenreResponse, error) {
	req.Sanitize()
	if len(req.Name) < 2 || len(req.Name) > 50 {
		return nil, exception.ErrInvalidInput
	}

	genre, err := s.genreRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	genre.Name = req.Name

	updatedGenre, err := s.genreRepo.Update(ctx, genre)
	if err != nil {
		return nil, err
	}

	return dto.ToGenreResponse(updatedGenre), nil
}

func (s *genreService) DeleteGenre(ctx context.Context, id int64) error {
	_, err := s.genreRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.genreRepo.Delete(ctx, id)
}
