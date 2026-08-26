package service

import (
	"context"
	"fmt"
	"io"
	"strings"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum/sort"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/storage"
	"github.com/mehdi1443/cinema-booking-system/pkg/image"
)

const maxCinemaImageSize = 5 * 1024 * 1024

type cinemaService struct {
	cinemaRepo repository_contract.CinemaInterface
	storage    storage.StorageInterFace
	log        domainlogger.Logger
}

func NewCinemaService(
	cinemaRepo repository_contract.CinemaInterface,
	storage storage.StorageInterFace,
	log domainlogger.Logger,
) service_contract.CinemaServiceInterface {
	return &cinemaService{
		cinemaRepo: cinemaRepo,
		storage:    storage,
		log:        log,
	}
}

func (s *cinemaService) Create(ctx context.Context, req *dto.CreateCinemaRequest) (*dto.CinemaDetailResponse, error) {
	if req == nil {
		return nil, exception.ErrInvalidInput
	}

	if !req.Status.IsValid() {
		return nil, exception.ErrInvalidInput
	}

	trimmedName := strings.TrimSpace(req.Name)
	trimmedAddress := strings.TrimSpace(req.Address)

	if trimmedName == "" || trimmedAddress == "" {
		return nil, exception.ErrInvalidInput
	}

	cinema := &entity.Cinema{
		CityID:    req.CityID,
		Name:      trimmedName,
		Address:   trimmedAddress,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Status:    req.Status,
	}

	if req.Description != "" {
		trimmedDesc := strings.TrimSpace(req.Description)
		cinema.Description = &trimmedDesc
	}
	if req.Phone != "" {
		trimmedPhone := strings.TrimSpace(req.Phone)
		cinema.Phone = &trimmedPhone
	}

	createdCinema, err := s.cinemaRepo.Create(ctx, cinema, req.FacilityIDs)
	if err != nil {
		return nil, fmt.Errorf("create cinema: %w", err)
	}

	s.log.Info(ctx, "Cinema created successfully", domainlogger.Int64("cinema_id", createdCinema.ID))

	return dto.ToCinemaDetailResponse(createdCinema), nil
}

func (s *cinemaService) GetByID(ctx context.Context, id int64) (*dto.CinemaDetailResponse, error) {
	if id <= 0 {
		return nil, exception.ErrInvalidInput
	}

	cinema, err := s.cinemaRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get cinema by id %d: %w", id, err)
	}

	return dto.ToCinemaDetailResponse(cinema), nil
}

func (s *cinemaService) GetCinemas(ctx context.Context, req dto.CinemaFilterRequest) ([]*dto.CinemaListResponse, int64, error) {
	if req.Status != nil && !req.Status.IsValid() {
		return nil, 0, exception.ErrInvalidInput
	}

	filter := repository_contract.CinemaFilter{
		ProvinceID:  req.ProvinceID,
		CityID:      req.CityID,
		FacilityIDs: req.FacilityIDs,
		Status:      req.Status,
		MinRating:   req.MinRating,
		MaxRating:   req.MaxRating,
	}

	if req.Search != nil {
		filter.SearchKeyword = strings.TrimSpace(*req.Search)
	}

	opts := repository_contract.NewQueryOptions()

	if req.SortBy != "" {
		if sortBy, ok := sort.ParseCinemaSortBy(req.SortBy); ok {
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

	cinemas, total, err := s.cinemaRepo.GetCinemas(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("get cinemas: %w", err)
	}

	return dto.ToCinemaListResponseSlice(cinemas), total, nil
}

func (s *cinemaService) Update(ctx context.Context, id int64, req *dto.UpdateCinemaRequest) (*dto.CinemaDetailResponse, error) {
	if id <= 0 || req == nil {
		return nil, exception.ErrInvalidInput
	}

	cinema, err := s.cinemaRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get cinema before update %d: %w", id, err)
	}

	if req.CityID != nil {
		if *req.CityID <= 0 {
			return nil, exception.ErrInvalidInput
		}
		cinema.CityID = *req.CityID
	}
	if req.Name != nil {
		trimmedName := strings.TrimSpace(*req.Name)
		if trimmedName == "" {
			return nil, exception.ErrInvalidInput
		}
		cinema.Name = trimmedName
	}
	if req.Description != nil {
		trimmedDesc := strings.TrimSpace(*req.Description)
		cinema.Description = &trimmedDesc
	}
	if req.Address != nil {
		trimmedAddress := strings.TrimSpace(*req.Address)
		if trimmedAddress == "" {
			return nil, exception.ErrInvalidInput
		}
		cinema.Address = trimmedAddress
	}
	if req.Latitude != nil {
		cinema.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		cinema.Longitude = req.Longitude
	}
	if req.Phone != nil {
		trimmedPhone := strings.TrimSpace(*req.Phone)
		cinema.Phone = &trimmedPhone
	}
	if req.Status != nil {
		if !req.Status.IsValid() {
			return nil, exception.ErrInvalidInput
		}
		cinema.Status = *req.Status
	}

	var facilityIDs []int64
	if req.FacilityIDs != nil {
		facilityIDs = *req.FacilityIDs
	} else {
		facilityIDs = make([]int64, 0, len(cinema.Facilities))
		for _, f := range cinema.Facilities {
			facilityIDs = append(facilityIDs, f.ID)
		}
	}

	updatedCinema, err := s.cinemaRepo.Update(ctx, cinema, facilityIDs)
	if err != nil {
		return nil, fmt.Errorf("update cinema %d: %w", id, err)
	}

	s.log.Info(ctx, "Cinema updated successfully", domainlogger.Int64("cinema_id", updatedCinema.ID))

	return dto.ToCinemaDetailResponse(updatedCinema), nil
}

func (s *cinemaService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return exception.ErrInvalidInput
	}

	cinema, err := s.cinemaRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get cinema before delete %d: %w", id, err)
	}

	if cinema.ImageURLLow != nil || cinema.ImageURLHigh != nil {
		_ = s.DeleteImage(ctx, id)
	}

	if err := s.cinemaRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete cinema %d: %w", id, err)
	}

	s.log.Info(ctx, "Cinema deleted successfully", domainlogger.Int64("cinema_id", id))

	return nil
}

func (s *cinemaService) UploadImage(ctx context.Context, id int64, file io.Reader, filename string) (*dto.CinemaDetailResponse, error) {
	if id <= 0 {
		return nil, exception.ErrInvalidInput
	}

	if file == nil {
		return nil, exception.ErrInvalidFileType
	}

	cinema, err := s.cinemaRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get cinema before image upload %d: %w", id, err)
	}

	fileReader, err := checkFileSize(file, maxCinemaImageSize)
	if err != nil {
		return nil, err
	}

	img, err := image.DecodeImage(fileReader)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrInvalidFileType, err)
	}

	lowImg := image.Resize(400, 300, img)
	lowBuf, err := image.EncodeToBuffer(lowImg, image.FormatWebP)
	if err != nil {
		return nil, fmt.Errorf("encode low resolution cinema image: %w", exception.ErrInternal)
	}

	lowKey := fmt.Sprintf("cinemas/cinema_%d_image_low.webp", id)
	lowURL, err := s.storage.UploadObject(ctx, enum.CinemaPhoto, lowKey, lowBuf, image.GetContentType(image.FormatWebP))
	if err != nil {
		return nil, fmt.Errorf("upload low resolution cinema image: %w", exception.ErrInternal)
	}

	highImg := image.Resize(1200, 900, img)
	highBuf, err := image.EncodeToBuffer(highImg, image.FormatJPEG)
	if err != nil {
		return nil, fmt.Errorf("encode high resolution cinema image: %w", exception.ErrInternal)
	}

	highKey := fmt.Sprintf("cinemas/cinema_%d_image_high.jpg", id)
	highURL, err := s.storage.UploadObject(ctx, enum.CinemaPhoto, highKey, highBuf, image.GetContentType(image.FormatJPEG))
	if err != nil {
		return nil, fmt.Errorf("upload high resolution cinema image: %w", exception.ErrInternal)
	}

	cinema.ImageURLLow = &lowURL
	cinema.ImageURLHigh = &highURL

	facilityIDs := make([]int64, 0, len(cinema.Facilities))
	for _, f := range cinema.Facilities {
		facilityIDs = append(facilityIDs, f.ID)
	}

	updatedCinema, err := s.cinemaRepo.Update(ctx, cinema, facilityIDs)
	if err != nil {
		return nil, fmt.Errorf("save image URLs for cinema %d: %w", id, err)
	}

	s.log.Info(ctx, "Cinema image uploaded successfully", domainlogger.Int64("cinema_id", id))

	return dto.ToCinemaDetailResponse(updatedCinema), nil
}

func (s *cinemaService) DeleteImage(ctx context.Context, id int64) error {
	if id <= 0 {
		return exception.ErrInvalidInput
	}

	cinema, err := s.cinemaRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get cinema before image deletion %d: %w", id, err)
	}

	lowKey := fmt.Sprintf("cinemas/cinema_%d_image_low.webp", id)
	highKey := fmt.Sprintf("cinemas/cinema_%d_image_high.jpg", id)

	if err := s.storage.DeleteObject(ctx, enum.CinemaPhoto, lowKey); err != nil {
		s.log.Warn(ctx, "Failed to delete low-res cinema image from storage", domainlogger.Int64("cinema_id", id), domainlogger.Error(err))
	}
	if err := s.storage.DeleteObject(ctx, enum.CinemaPhoto, highKey); err != nil {
		s.log.Warn(ctx, "Failed to delete high-res cinema image from storage", domainlogger.Int64("cinema_id", id), domainlogger.Error(err))
	}

	cinema.ImageURLLow = nil
	cinema.ImageURLHigh = nil

	facilityIDs := make([]int64, 0, len(cinema.Facilities))
	for _, f := range cinema.Facilities {
		facilityIDs = append(facilityIDs, f.ID)
	}

	_, err = s.cinemaRepo.Update(ctx, cinema, facilityIDs)
	if err != nil {
		return fmt.Errorf("remove image URLs for cinema %d: %w", id, err)
	}

	s.log.Info(ctx, "Cinema image deleted successfully", domainlogger.Int64("cinema_id", id))

	return nil
}
