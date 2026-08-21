package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum/sort"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/storage"
	"github.com/mehdi1443/cinema-booking-system/pkg/image"
)

const (
	maxPosterSize  = 5 * 1024 * 1024   // 5 MB
	maxTrailerSize = 100 * 1024 * 1024 // 100 MB
)

type movieService struct {
	movieRepo repository_contract.MovieInterFace
	storage   storage.StorageInterFace
}

func NewMovieService(
	movieRepo repository_contract.MovieInterFace,
	storage storage.StorageInterFace,
) service_contract.MovieServiceInterface {
	return &movieService{
		movieRepo: movieRepo,
		storage:   storage,
	}
}

// تابع کمکی بررسی و اعمال محدودیت حجم فایل
func checkFileSize(r io.Reader, maxSize int64) (io.Reader, error) {
	lr := io.LimitReader(r, maxSize+1)
	buf, err := io.ReadAll(lr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrInvalidFileType, err)
	}
	if int64(len(buf)) > maxSize {
		return nil, exception.ErrFileSizeTooLarge
	}
	return bytes.NewReader(buf), nil
}

func (s *movieService) Create(ctx context.Context, req *dto.CreateMovieRequest) (*dto.MovieDetailResponse, error) {
	if !req.Status.IsValid() {
		return nil, exception.ErrInvalidInput
	}

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
		if !p.Role.IsValid() {
			return nil, exception.ErrInvalidInput
		}
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

func (s *movieService) GetMovies(ctx context.Context, req dto.MovieFilterRequest) ([]*dto.MovieListResponse, int64, error) {
	if req.Status != nil && !req.Status.IsValid() {
		return nil, 0, exception.ErrInvalidInput
	}
	if req.PersonRole != nil && !req.PersonRole.IsValid() {
		return nil, 0, exception.ErrInvalidInput
	}

	filter := repository_contract.MovieFilter{
		Status:      req.Status,
		GenreIDs:    req.GenreIDs,
		PersonIDs:   req.PersonIDs,
		PersonRole:  req.PersonRole,
		Country:     req.Country,
		MinRating:   req.MinRating,
		MaxRating:   req.MaxRating,
		FromYear:    req.FromYear,
		ToYear:      req.ToYear,
		MinDuration: req.MinDuration,
		MaxDuration: req.MaxDuration,
		HasTrailer:  req.HasTrailer,
	}

	if req.Search != nil {
		filter.SearchKeyword = *req.Search
	}

	opts := repository_contract.NewQueryOptions()

	if req.SortBy != "" {
		if sortBy, ok := sort.ParseMovieSortBy(req.SortBy); ok {
			asc := false
			if req.SortDir != "" {
				if dir, dirOk := sort.ParseSortDirection(req.SortDir); dirOk && dir == sort.SortAscending {
					asc = true
				}
			}
			opts.WithSorting(sortBy.DBColumn(), asc)
		}
	}

	req.Normalize()
	opts.WithPagination(req.Limit, req.GetOffset())

	movies, total, err := s.movieRepo.GetMovies(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}

	return dto.ToMovieListResponseSlice(movies), total, nil
}

func (s *movieService) Update(ctx context.Context, id int64, req *dto.UpdateMovieRequest) (*dto.MovieDetailResponse, error) {
	movie, err := s.movieRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		movie.Title = *req.Title
	}
	if req.Description != nil {
		movie.Description = *req.Description
	}
	if req.DurationMinutes != nil {
		movie.DurationMinutes = *req.DurationMinutes
	}
	if req.ReleaseDate != nil {
		movie.ReleaseDate = *req.ReleaseDate
	}
	if req.Country != nil {
		movie.Country = *req.Country
	}
	if req.Status != nil {
		if !req.Status.IsValid() {
			return nil, exception.ErrInvalidInput
		}
		movie.Status = *req.Status
	}

	var genreIDs []int64
	if req.GenreIDs != nil {
		genreIDs = *req.GenreIDs
	} else {
		genreIDs = make([]int64, 0, len(movie.Genres))
		for _, g := range movie.Genres {
			genreIDs = append(genreIDs, g.ID)
		}
	}

	var peopleMap map[int64]enum.MoviePersonRole
	if req.People != nil {
		peopleMap = make(map[int64]enum.MoviePersonRole, len(*req.People))
		for _, p := range *req.People {
			if !p.Role.IsValid() {
				return nil, exception.ErrInvalidInput
			}
			peopleMap[p.PersonID] = p.Role
		}
	}

	updatedMovie, err := s.movieRepo.Update(ctx, movie, genreIDs, peopleMap)
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

	// ۱. بررسی حجم فایل (سقف ۵ مگابایت)
	fileReader, err := checkFileSize(file, maxPosterSize)
	if err != nil {
		return nil, err
	}

	// ۲. بررسی فرمت تصویر (دیکود شدن صحت فرمت JPEG/PNG/WebP را تضمین می‌کند)
	img, err := image.DecodeImage(fileReader)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrInvalidFileType, err)
	}

	lowImg := image.Resize(300, 450, img)
	lowBuf, err := image.EncodeToBuffer(lowImg, image.FormatWebP)
	if err != nil {
		return nil, fmt.Errorf("failed to encode low-res poster: %w", err)
	}
	lowKey := fmt.Sprintf("movies/movie_%d_poster_low.webp", id)
	lowURL, err := s.storage.UploadObject(ctx, enum.MoviePoster, lowKey, lowBuf, image.GetContentType(image.FormatWebP))
	if err != nil {
		return nil, fmt.Errorf("failed to upload low-res poster: %w", err)
	}

	highImg := image.Resize(1000, 1500, img)
	highBuf, err := image.EncodeToBuffer(highImg, image.FormatJPEG)
	if err != nil {
		return nil, fmt.Errorf("failed to encode high-res poster: %w", err)
	}
	highKey := fmt.Sprintf("movies/movie_%d_poster_high.jpg", id)
	highURL, err := s.storage.UploadObject(ctx, enum.MoviePoster, highKey, highBuf, image.GetContentType(image.FormatJPEG))
	if err != nil {
		return nil, fmt.Errorf("failed to upload high-res poster: %w", err)
	}

	movie.PosterURLLow = lowURL
	movie.PosterURLHigh = highURL

	genreIDs := make([]int64, 0, len(movie.Genres))
	for _, g := range movie.Genres {
		genreIDs = append(genreIDs, g.ID)
	}

	updatedMovie, err := s.movieRepo.Update(ctx, movie, genreIDs, nil)
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

	genreIDs := make([]int64, 0, len(movie.Genres))
	for _, g := range movie.Genres {
		genreIDs = append(genreIDs, g.ID)
	}

	_, err = s.movieRepo.Update(ctx, movie, genreIDs, nil)
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

	// ۱. بررسی حجم فایل (سقف ۱۰۰ مگابایت)
	fileReader, err := checkFileSize(file, maxTrailerSize)
	if err != nil {
		return nil, err
	}

	// ۲. اعتبارسنجی فرمت واقعی ویدیو با ۵۱۲ بایت اول فایل (Magic Bytes)
	header := make([]byte, 512)
	n, readErr := fileReader.Read(header)
	if readErr != nil && readErr != io.EOF {
		return nil, fmt.Errorf("%w: failed to read file header", exception.ErrInvalidFileType)
	}

	contentType := http.DetectContentType(header[:n])
	if contentType != "video/mp4" && contentType != "video/webm" && contentType != "application/octet-stream" {
		return nil, exception.ErrInvalidFileType
	}

	// بازسازی Reader برای آپلود کامل فایل
	fullStream := io.MultiReader(bytes.NewReader(header[:n]), fileReader)

	trailerKey := fmt.Sprintf("movies/movie_%d_trailer.mp4", id)
	trailerURL, err := s.storage.UploadObject(ctx, enum.MovieTrailer, trailerKey, fullStream, "video/mp4")
	if err != nil {
		return nil, fmt.Errorf("failed to upload trailer: %w", err)
	}

	movie.TrailerURL = trailerURL

	genreIDs := make([]int64, 0, len(movie.Genres))
	for _, g := range movie.Genres {
		genreIDs = append(genreIDs, g.ID)
	}

	updatedMovie, err := s.movieRepo.Update(ctx, movie, genreIDs, nil)
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

	genreIDs := make([]int64, 0, len(movie.Genres))
	for _, g := range movie.Genres {
		genreIDs = append(genreIDs, g.ID)
	}

	_, err = s.movieRepo.Update(ctx, movie, genreIDs, nil)
	if err != nil {
		return err
	}

	return nil
}
