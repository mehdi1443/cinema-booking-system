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
	maxPosterSize  = 5 * 1024 * 1024
	maxTrailerSize = 100 * 1024 * 1024
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

// checkFileSize validates maximum file size without changing
// the business logic of the service.
func checkFileSize(r io.Reader, maxSize int64) (io.Reader, error) {
	lr := io.LimitReader(r, maxSize+1)

	buf, err := io.ReadAll(lr)
	if err != nil {
		return nil, fmt.Errorf(
			"%w: failed to read file",
			exception.ErrInvalidFileType,
		)
	}

	if int64(len(buf)) > maxSize {
		return nil, exception.ErrFileSizeTooLarge
	}

	return bytes.NewReader(buf), nil
}

// Create creates a movie and returns its detail response.
func (s *movieService) Create(
	ctx context.Context,
	req *dto.CreateMovieRequest,
) (*dto.MovieDetailResponse, error) {

	if req == nil {
		return nil, exception.ErrInvalidInput
	}

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

	peopleMap := make(
		map[int64]enum.MoviePersonRole,
		len(req.People),
	)

	for _, p := range req.People {

		if !p.Role.IsValid() {
			return nil, exception.ErrInvalidInput
		}

		peopleMap[p.PersonID] = p.Role
	}

	createdMovie, err := s.movieRepo.Create(
		ctx,
		movie,
		req.GenreIDs,
		peopleMap,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"create movie: %w",
			err,
		)
	}

	return dto.ToMovieDetailResponse(createdMovie), nil
}

// GetByID retrieves a movie by ID.
func (s *movieService) GetByID(
	ctx context.Context,
	id int64,
) (*dto.MovieDetailResponse, error) {

	if id <= 0 {
		return nil, exception.ErrInvalidInput
	}

	movie, err := s.movieRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf(
			"get movie by id %d: %w",
			id,
			err,
		)
	}

	return dto.ToMovieDetailResponse(movie), nil
}

// GetMovies retrieves movies based on filters.
func (s *movieService) GetMovies(
	ctx context.Context,
	req dto.MovieFilterRequest,
) ([]*dto.MovieListResponse, int64, error) {

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

				if dir, dirOK := sort.ParseSortDirection(req.SortDir); dirOK &&
					dir == sort.SortAscending {
					asc = true
				}
			}

			opts.WithSorting(
				sortBy.DBColumn(),
				asc,
			)
		}
	}

	req.Normalize()

	opts.WithPagination(
		req.Limit,
		req.GetOffset(),
	)

	movies, total, err := s.movieRepo.GetMovies(
		ctx,
		filter,
		opts,
	)

	if err != nil {
		return nil, 0, fmt.Errorf(
			"get movies: %w",
			err,
		)
	}

	return dto.ToMovieListResponseSlice(movies), total, nil
}

// Update updates a movie.
func (s *movieService) Update(
	ctx context.Context,
	id int64,
	req *dto.UpdateMovieRequest,
) (*dto.MovieDetailResponse, error) {

	if id <= 0 {
		return nil, exception.ErrInvalidInput
	}

	if req == nil {
		return nil, exception.ErrInvalidInput
	}

	movie, err := s.movieRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf(
			"get movie before update %d: %w",
			id,
			err,
		)
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

		genreIDs = make(
			[]int64,
			0,
			len(movie.Genres),
		)

		for _, g := range movie.Genres {
			genreIDs = append(
				genreIDs,
				g.ID,
			)
		}
	}

	var peopleMap map[int64]enum.MoviePersonRole

	if req.People != nil {

		peopleMap = make(
			map[int64]enum.MoviePersonRole,
			len(*req.People),
		)

		for _, p := range *req.People {

			if !p.Role.IsValid() {
				return nil, exception.ErrInvalidInput
			}

			peopleMap[p.PersonID] = p.Role
		}
	}

	updatedMovie, err := s.movieRepo.Update(
		ctx,
		movie,
		genreIDs,
		peopleMap,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"update movie %d: %w",
			id,
			err,
		)
	}

	return dto.ToMovieDetailResponse(updatedMovie), nil
}

// Delete deletes a movie.
func (s *movieService) Delete(
	ctx context.Context,
	id int64,
) error {

	if id <= 0 {
		return exception.ErrInvalidInput
	}

	movie, err := s.movieRepo.GetByID(
		ctx,
		id,
	)

	if err != nil {
		return fmt.Errorf(
			"get movie before delete %d: %w",
			id,
			err,
		)
	}

	if movie.PosterURLLow != "" ||
		movie.PosterURLHigh != "" {
		_ = s.DeletePoster(ctx, id)
	}

	if movie.TrailerURL != "" {
		_ = s.DeleteTrailer(ctx, id)
	}

	if err := s.movieRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf(
			"delete movie %d: %w",
			id,
			err,
		)
	}

	return nil
}

// UploadPoster uploads and stores a movie poster.
func (s *movieService) UploadPoster(
	ctx context.Context,
	id int64,
	file io.Reader,
) (*dto.MovieDetailResponse, error) {

	if id <= 0 {
		return nil, exception.ErrInvalidInput
	}

	if file == nil {
		return nil, exception.ErrInvalidFileType
	}

	movie, err := s.movieRepo.GetByID(
		ctx,
		id,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"get movie before poster upload %d: %w",
			id,
			err,
		)
	}

	fileReader, err := checkFileSize(
		file,
		maxPosterSize,
	)

	if err != nil {
		return nil, err
	}

	img, err := image.DecodeImage(
		fileReader,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"%w: invalid poster image",
			exception.ErrInvalidFileType,
		)
	}

	lowImg := image.Resize(
		300,
		450,
		img,
	)

	lowBuf, err := image.EncodeToBuffer(
		lowImg,
		image.FormatWebP,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"encode low resolution poster: %w",
			exception.ErrInternal,
		)
	}

	lowKey := fmt.Sprintf(
		"movies/movie_%d_poster_low.webp",
		id,
	)

	lowURL, err := s.storage.UploadObject(
		ctx,
		enum.MoviePoster,
		lowKey,
		lowBuf,
		image.GetContentType(image.FormatWebP),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"upload low resolution poster: %w",
			exception.ErrInternal,
		)
	}

	highImg := image.Resize(
		1000,
		1500,
		img,
	)

	highBuf, err := image.EncodeToBuffer(
		highImg,
		image.FormatJPEG,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"encode high resolution poster: %w",
			exception.ErrInternal,
		)
	}

	highKey := fmt.Sprintf(
		"movies/movie_%d_poster_high.jpg",
		id,
	)

	highURL, err := s.storage.UploadObject(
		ctx,
		enum.MoviePoster,
		highKey,
		highBuf,
		image.GetContentType(image.FormatJPEG),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"upload high resolution poster: %w",
			exception.ErrInternal,
		)
	}

	movie.PosterURLLow = lowURL
	movie.PosterURLHigh = highURL

	genreIDs := make(
		[]int64,
		0,
		len(movie.Genres),
	)

	for _, g := range movie.Genres {
		genreIDs = append(
			genreIDs,
			g.ID,
		)
	}

	updatedMovie, err := s.movieRepo.Update(
		ctx,
		movie,
		genreIDs,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"save poster URLs for movie %d: %w",
			id,
			err,
		)
	}

	return dto.ToMovieDetailResponse(updatedMovie), nil
}

// DeletePoster removes movie poster files and URLs.
func (s *movieService) DeletePoster(
	ctx context.Context,
	id int64,
) error {

	if id <= 0 {
		return exception.ErrInvalidInput
	}

	movie, err := s.movieRepo.GetByID(
		ctx,
		id,
	)

	if err != nil {
		return fmt.Errorf(
			"get movie before poster deletion %d: %w",
			id,
			err,
		)
	}

	lowKey := fmt.Sprintf(
		"movies/movie_%d_poster_low.webp",
		id,
	)

	highKey := fmt.Sprintf(
		"movies/movie_%d_poster_high.jpg",
		id,
	)

	// Keep existing behavior:
	// storage deletion errors are intentionally ignored.
	_ = s.storage.DeleteObject(
		ctx,
		enum.MoviePoster,
		lowKey,
	)

	_ = s.storage.DeleteObject(
		ctx,
		enum.MoviePoster,
		highKey,
	)

	movie.PosterURLLow = ""
	movie.PosterURLHigh = ""

	genreIDs := make(
		[]int64,
		0,
		len(movie.Genres),
	)

	for _, g := range movie.Genres {
		genreIDs = append(
			genreIDs,
			g.ID,
		)
	}

	_, err = s.movieRepo.Update(
		ctx,
		movie,
		genreIDs,
		nil,
	)

	if err != nil {
		return fmt.Errorf(
			"remove poster URLs for movie %d: %w",
			id,
			err,
		)
	}

	return nil
}

// UploadTrailer uploads a movie trailer.
func (s *movieService) UploadTrailer(
	ctx context.Context,
	id int64,
	file io.Reader,
) (*dto.MovieDetailResponse, error) {

	if id <= 0 {
		return nil, exception.ErrInvalidInput
	}

	if file == nil {
		return nil, exception.ErrInvalidFileType
	}

	movie, err := s.movieRepo.GetByID(
		ctx,
		id,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"get movie before trailer upload %d: %w",
			id,
			err,
		)
	}

	fileReader, err := checkFileSize(
		file,
		maxTrailerSize,
	)

	if err != nil {
		return nil, err
	}

	header := make([]byte, 512)

	n, readErr := fileReader.Read(header)

	if readErr != nil &&
		readErr != io.EOF {

		return nil, fmt.Errorf(
			"%w: failed to read trailer header",
			exception.ErrInvalidFileType,
		)
	}

	contentType := http.DetectContentType(
		header[:n],
	)

	if contentType != "video/mp4" &&
		contentType != "video/webm" &&
		contentType != "application/octet-stream" {

		return nil, exception.ErrInvalidFileType
	}

	fullStream := io.MultiReader(
		bytes.NewReader(header[:n]),
		fileReader,
	)

	trailerKey := fmt.Sprintf(
		"movies/movie_%d_trailer.mp4",
		id,
	)

	trailerURL, err := s.storage.UploadObject(
		ctx,
		enum.MovieTrailer,
		trailerKey,
		fullStream,
		"video/mp4",
	)

	if err != nil {
		return nil, fmt.Errorf(
			"upload trailer for movie %d: %w",
			id,
			exception.ErrInternal,
		)
	}

	movie.TrailerURL = trailerURL

	genreIDs := make(
		[]int64,
		0,
		len(movie.Genres),
	)

	for _, g := range movie.Genres {
		genreIDs = append(
			genreIDs,
			g.ID,
		)
	}

	updatedMovie, err := s.movieRepo.Update(
		ctx,
		movie,
		genreIDs,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"save trailer URL for movie %d: %w",
			id,
			err,
		)
	}

	return dto.ToMovieDetailResponse(updatedMovie), nil
}

// DeleteTrailer removes movie trailer.
func (s *movieService) DeleteTrailer(
	ctx context.Context,
	id int64,
) error {

	if id <= 0 {
		return exception.ErrInvalidInput
	}

	movie, err := s.movieRepo.GetByID(
		ctx,
		id,
	)

	if err != nil {
		return fmt.Errorf(
			"get movie before trailer deletion %d: %w",
			id,
			err,
		)
	}

	trailerKey := fmt.Sprintf(
		"movies/movie_%d_trailer.mp4",
		id,
	)

	// Keep existing behavior:
	// storage deletion errors are intentionally ignored.
	_ = s.storage.DeleteObject(
		ctx,
		enum.MovieTrailer,
		trailerKey,
	)

	movie.TrailerURL = ""

	genreIDs := make(
		[]int64,
		0,
		len(movie.Genres),
	)

	for _, g := range movie.Genres {
		genreIDs = append(
			genreIDs,
			g.ID,
		)
	}

	_, err = s.movieRepo.Update(
		ctx,
		movie,
		genreIDs,
		nil,
	)

	if err != nil {
		return fmt.Errorf(
			"remove trailer URL for movie %d: %w",
			id,
			err,
		)
	}

	return nil
}
