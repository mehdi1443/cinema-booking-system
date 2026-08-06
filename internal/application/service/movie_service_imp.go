package service

import (
	"context"
	"errors"
	"fmt"
	"io"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/storage"
	"github.com/mehdi1443/cinema-booking-system/pkg/imageutil"
)

const maxTrailerSize = 100 * 1024 * 1024

type movieService struct {
	movieRepo repository_contract.MovieInterface
	storage   storage.StorageInterFace
}

func NewMovieService(
	movieRepo repository_contract.MovieInterface,
	storage storage.StorageInterFace,
) service_contract.MovieServiceInterface {
	return &movieService{
		movieRepo: movieRepo,
		storage:   storage,
	}
}

func (s *movieService) Create(ctx context.Context, req *dto.CreateMovieRequest) (*dto.MovieDetailResponse, error) {
	movie := &entity.Movie{
		Title:           req.Title,
		Description:     req.Description,
		DurationMinutes: req.DurationMinutes,
		ReleaseDate:     req.ReleaseDate,
		Country:         req.Country,
		Status:          req.Status,
	}

	peopleMap := make(map[int64]enum.MoviePersonRole, len(req.People))
	for _, p := range req.People {
		peopleMap[p.PersonID] = p.Role
	}

	createdMovie, err := s.movieRepo.Create(ctx, movie, req.GenreIDs, peopleMap)
	if err != nil {
		return nil, err
	}

	return dto.ToMovieDetailResponse(createdMovie), nil
}

func (s *movieService) GetByID(ctx context.Context, id int64) (*dto.MovieDetailResponse, error) {
	movie, err := s.movieRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.ToMovieDetailResponse(movie), nil
}

func (s *movieService) GetByStatus(ctx context.Context, status enum.MovieStatus) ([]*dto.MovieListResponse, error) {
	movies, err := s.movieRepo.GetByStatus(ctx, status)
	if err != nil {
		return nil, err
	}

	return dto.ToMovieListResponseSlice(movies), nil
}

func (s *movieService) GetByGenreID(ctx context.Context, genreID int64) ([]*dto.MovieListResponse, error) {
	movies, err := s.movieRepo.GetByGenreID(ctx, genreID)
	if err != nil {
		return nil, err
	}

	return dto.ToMovieListResponseSlice(movies), nil
}

func (s *movieService) GetAll(ctx context.Context) ([]*dto.MovieListResponse, error) {
	movies, err := s.movieRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return dto.ToMovieListResponseSlice(movies), nil
}

func (s *movieService) Update(ctx context.Context, id int64, req *dto.UpdateMovieRequest) (*dto.MovieDetailResponse, error) {
	movie, err := s.movieRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	movie.Title = req.Title
	movie.Description = req.Description
	movie.DurationMinutes = req.DurationMinutes
	movie.ReleaseDate = req.ReleaseDate
	movie.Country = req.Country
	movie.Status = req.Status

	peopleMap := make(map[int64]enum.MoviePersonRole, len(req.People))
	for _, p := range req.People {
		peopleMap[p.PersonID] = p.Role
	}

	updatedMovie, err := s.movieRepo.Update(ctx, movie, req.GenreIDs, peopleMap)
	if err != nil {
		return nil, err
	}

	return dto.ToMovieDetailResponse(updatedMovie), nil
}

func (s *movieService) Delete(ctx context.Context, id int64) error {
	movie, err := s.movieRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if movie.PosterURLLow != "" || movie.PosterURLHigh != "" {
		_ = s.DeletePoster(ctx, id)
	}

	if movie.TrailerURL != "" {
		_ = s.DeleteTrailer(ctx, id)
	}

	return s.movieRepo.Delete(ctx, id)
}

func (s *movieService) UploadPoster(ctx context.Context, id int64, file io.Reader) (*dto.MovieDetailResponse, error) {
	movie, err := s.movieRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	img, err := imageutil.DecodeImage(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode poster image: %w", err)
	}

	lowImg := imageutil.Resize(300, 450, img)
	lowBuf, err := imageutil.EncodeToBuffer(lowImg, imageutil.FormatWebP)
	if err != nil {
		return nil, fmt.Errorf("failed to encode low-res poster: %w", err)
	}
	lowKey := fmt.Sprintf("movies/movie_%d_poster_low.webp", id)
	lowURL, err := s.storage.UploadObject(ctx, enum.MoviePoster, lowKey, lowBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to upload low-res poster: %w", err)
	}

	highImg := imageutil.Resize(1000, 1500, img)
	highBuf, err := imageutil.EncodeToBuffer(highImg, imageutil.FormatJPEG)
	if err != nil {
		return nil, fmt.Errorf("failed to encode high-res poster: %w", err)
	}
	highKey := fmt.Sprintf("movies/movie_%d_poster_high.jpg", id)
	highURL, err := s.storage.UploadObject(ctx, enum.MoviePoster, highKey, highBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to upload high-res poster: %w", err)
	}

	movie.PosterURLLow = lowURL
	movie.PosterURLHigh = highURL

	updatedMovie, err := s.movieRepo.Update(ctx, movie, nil, nil)
	if err != nil {
		return nil, err
	}

	return dto.ToMovieDetailResponse(updatedMovie), nil
}

func (s *movieService) DeletePoster(ctx context.Context, id int64) error {
	movie, err := s.movieRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	lowKey := fmt.Sprintf("movies/movie_%d_poster_low.webp", id)
	highKey := fmt.Sprintf("movies/movie_%d_poster_high.jpg", id)

	_ = s.storage.DeleteObject(ctx, enum.MoviePoster, lowKey)
	_ = s.storage.DeleteObject(ctx, enum.MoviePoster, highKey)

	movie.PosterURLLow = ""
	movie.PosterURLHigh = ""

	_, err = s.movieRepo.Update(ctx, movie, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *movieService) UploadTrailer(ctx context.Context, id int64, file io.Reader) (*dto.MovieDetailResponse, error) {
	movie, err := s.movieRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	limitedReader := io.LimitReader(file, maxTrailerSize+1)
	trailerKey := fmt.Sprintf("movies/movie_%d_trailer.mp4", id)

	trailerURL, err := s.storage.UploadObject(ctx, enum.MovieTrailer, trailerKey, limitedReader)
	if err != nil {
		return nil, fmt.Errorf("failed to upload trailer: %w", err)
	}

	movie.TrailerURL = trailerURL

	updatedMovie, err := s.movieRepo.Update(ctx, movie, nil, nil)
	if err != nil {
		return nil, err
	}

	return dto.ToMovieDetailResponse(updatedMovie), nil
}

func (s *movieService) DeleteTrailer(ctx context.Context, id int64) error {
	movie, err := s.movieRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	trailerKey := fmt.Sprintf("movies/movie_%d_trailer.mp4", id)

	_ = s.storage.DeleteObject(ctx, enum.MovieTrailer, trailerKey)

	movie.TrailerURL = ""

	_, err = s.movieRepo.Update(ctx, movie, nil, nil)
	if err != nil {
		return err
	}

	return nil
}
