package service

import (
	"context"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
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
	genre, err := s.genreRepo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return dto.ToGenreResponse(genre), nil
}

func (s *genreService) GetAllGenres(ctx context.Context) ([]*dto.GenreResponse, error) {
	genres, err := s.genreRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return dto.ToGenreResponseSlice(genres), nil
}

func (s *genreService) UpdateGenre(ctx context.Context, id int64, req *dto.UpdateGenreRequest) (*dto.GenreResponse, error) {
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
