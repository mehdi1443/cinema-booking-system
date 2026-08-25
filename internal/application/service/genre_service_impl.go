package service

import (
	"context"
	"fmt"
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
	if req == nil {
		return nil, exception.ErrInvalidInput
	}

	trimmedName := strings.TrimSpace(req.Name)
	if trimmedName == "" {
		return nil, exception.ErrInvalidInput
	}

	genre := &entity.Genre{
		Name: trimmedName,
	}

	createdGenre, err := s.genreRepo.Create(ctx, genre)
	if err != nil {
		return nil, fmt.Errorf("create genre: %w", err)
	}

	return dto.ToGenreResponse(createdGenre), nil
}

func (s *genreService) GetGenreByID(ctx context.Context, id int64) (*dto.GenreResponse, error) {
	if id <= 0 {
		return nil, exception.ErrInvalidInput
	}

	genre, err := s.genreRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get genre by id %d: %w", id, err)
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
		return nil, fmt.Errorf("get genre by name %s: %w", trimmedName, err)
	}

	return dto.ToGenreResponse(genre), nil
}

func (s *genreService) GetGenres(ctx context.Context, req dto.GenreFilterRequest) ([]*dto.GenreResponse, int64, error) {
	filter := repository_contract.GenreFilter{}
	if req.Search != nil {
		filter.SearchKeyword = strings.TrimSpace(*req.Search)
	}

	req.Normalize()

	opts := repository_contract.NewQueryOptions()
	opts.WithPagination(req.Limit, req.GetOffset())

	genres, total, err := s.genreRepo.GetGenres(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("get genres: %w", err)
	}

	return dto.ToGenreResponseSlice(genres), total, nil
}

func (s *genreService) UpdateGenre(ctx context.Context, id int64, req *dto.UpdateGenreRequest) (*dto.GenreResponse, error) {
	if id <= 0 || req == nil {
		return nil, exception.ErrInvalidInput
	}

	trimmedName := strings.TrimSpace(req.Name)
	if trimmedName == "" {
		return nil, exception.ErrInvalidInput
	}

	genre, err := s.genreRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get genre before update %d: %w", id, err)
	}

	genre.Name = trimmedName

	updatedGenre, err := s.genreRepo.Update(ctx, genre)
	if err != nil {
		return nil, fmt.Errorf("update genre %d: %w", id, err)
	}

	return dto.ToGenreResponse(updatedGenre), nil
}

func (s *genreService) DeleteGenre(ctx context.Context, id int64) error {
	if id <= 0 {
		return exception.ErrInvalidInput
	}

	_, err := s.genreRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get genre before delete %d: %w", id, err)
	}

	if err := s.genreRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete genre %d: %w", id, err)
	}

	return nil
}
