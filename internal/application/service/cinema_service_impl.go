package service

import (
	"context"
	"fmt"
	"io"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/storage"
	"github.com/mehdi1443/cinema-booking-system/pkg/image"
)

type cinemaService struct {
	cinemaRepo repository_contract.CinemaRepository
	storage    storage.StorageInterFace
}

func NewCinemaService(
	cinemaRepo repository_contract.CinemaRepository,
	storage storage.StorageInterFace,
) service_contract.CinemaServiceInterface {
	return &cinemaService{
		cinemaRepo: cinemaRepo,
		storage:    storage,
	}
}

func (s *cinemaService) Create(ctx context.Context, req *dto.CreateCinemaRequest) (*dto.CinemaDetailResponse, error) {
	cinema := &entity.Cinema{
		CityID:    req.CityID,
		Name:      req.Name,
		Address:   req.Address,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Status:    req.Status,
	}

	if req.Description != "" {
		cinema.Description = &req.Description
	}
	if req.Phone != "" {
		cinema.Phone = &req.Phone
	}

	createdCinema, err := s.cinemaRepo.Create(ctx, cinema)
	if err != nil {
		return nil, fmt.Errorf("failed to create cinema: %w", err)
	}

	if len(req.FacilityIDs) > 0 {
		if err := s.cinemaRepo.UpdateFacilities(ctx, createdCinema.ID, req.FacilityIDs); err != nil {
			return nil, fmt.Errorf("failed to assign facilities to cinema: %w", err)
		}
		// Refresh entity to reflect relations
		createdCinema, _ = s.cinemaRepo.GetByID(ctx, createdCinema.ID)
	}

	return dto.ToCinemaDetailResponse(createdCinema), nil
}

func (s *cinemaService) GetByID(ctx context.Context, id int64) (*dto.CinemaDetailResponse, error) {
	cinema, err := s.cinemaRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.ToCinemaDetailResponse(cinema), nil
}

func (s *cinemaService) GetWithFilters(ctx context.Context, req *dto.CinemaFilterRequest) (*dto.CinemaListWithCountResponse, error) {
	filter := repository_contract.CinemaFilter{
		ProvinceID:  req.ProvinceID,
		CityID:      req.CityID,
		FacilityIDs: req.FacilityIDs,
		Status:      req.Status,
		Search:      req.Search,
		Limit:       req.Limit,
		Offset:      req.Offset,
	}

	cinemas, err := s.cinemaRepo.GetWithFilters(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cinemas: %w", err)
	}

	total, err := s.cinemaRepo.CountWithFilters(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count cinemas: %w", err)
	}

	return &dto.CinemaListWithCountResponse{
		Cinemas: dto.ToCinemaListResponseSlice(cinemas),
		Total:   total,
	}, nil
}

func (s *cinemaService) Update(ctx context.Context, id int64, req *dto.UpdateCinemaRequest) (*dto.CinemaDetailResponse, error) {
	cinema, err := s.cinemaRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.CityID != nil {
		cinema.CityID = *req.CityID
	}
	if req.Name != nil {
		cinema.Name = *req.Name
	}
	if req.Description != nil {
		cinema.Description = req.Description
	}
	if req.Address != nil {
		cinema.Address = *req.Address
	}
	if req.Latitude != nil {
		cinema.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		cinema.Longitude = req.Longitude
	}
	if req.Phone != nil {
		cinema.Phone = req.Phone
	}
	if req.Status != nil {
		cinema.Status = *req.Status
	}

	updatedCinema, err := s.cinemaRepo.Update(ctx, cinema)
	if err != nil {
		return nil, fmt.Errorf("failed to update cinema: %w", err)
	}

	if req.FacilityIDs != nil {
		if err := s.cinemaRepo.UpdateFacilities(ctx, id, *req.FacilityIDs); err != nil {
			return nil, fmt.Errorf("failed to update facilities: %w", err)
		}
		updatedCinema, _ = s.cinemaRepo.GetByID(ctx, id)
	}

	return dto.ToCinemaDetailResponse(updatedCinema), nil
}

func (s *cinemaService) UploadImage(ctx context.Context, id int64, file io.Reader, filename string) (*dto.CinemaDetailResponse, error) {
	cinema, err := s.cinemaRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	img, err := image.DecodeImage(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	lowImg := image.Resize(400, 300, img)
	lowBuf, err := image.EncodeToBuffer(lowImg, image.FormatWebP)
	if err != nil {
		return nil, fmt.Errorf("failed to encode low-res image: %w", err)
	}
	lowKey := fmt.Sprintf("cinemas/cinema_%d_image_low.webp", id)
	lowURL, err := s.storage.UploadObject(ctx, enum.CinemaPhoto, lowKey, lowBuf, image.GetContentType(image.FormatWebP))
	if err != nil {
		return nil, fmt.Errorf("failed to upload low-res image: %w", err)
	}

	highImg := image.Resize(1200, 900, img)
	highBuf, err := image.EncodeToBuffer(highImg, image.FormatJPEG)
	if err != nil {
		return nil, fmt.Errorf("failed to encode high-res image: %w", err)
	}
	highKey := fmt.Sprintf("cinemas/cinema_%d_image_high.jpg", id)
	highURL, err := s.storage.UploadObject(ctx, enum.CinemaPhoto, highKey, highBuf, image.GetContentType(image.FormatJPEG))
	if err != nil {
		return nil, fmt.Errorf("failed to upload high-res image: %w", err)
	}

	cinema.ImageURLLow = &lowURL
	cinema.ImageURLHigh = &highURL

	updatedCinema, err := s.cinemaRepo.Update(ctx, cinema)
	if err != nil {
		return nil, fmt.Errorf("failed to update cinema image URL: %w", err)
	}

	return dto.ToCinemaDetailResponse(updatedCinema), nil
}

func (s *cinemaService) DeleteImage(ctx context.Context, id int64) error {
	cinema, err := s.cinemaRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	lowKey := fmt.Sprintf("cinemas/cinema_%d_image_low.webp", id)
	highKey := fmt.Sprintf("cinemas/cinema_%d_image_high.jpg", id)

	_ = s.storage.DeleteObject(ctx, enum.CinemaPhoto, lowKey)
	_ = s.storage.DeleteObject(ctx, enum.CinemaPhoto, highKey)

	cinema.ImageURLLow = nil
	cinema.ImageURLHigh = nil

	_, err = s.cinemaRepo.Update(ctx, cinema)
	if err != nil {
		return fmt.Errorf("failed to remove cinema image URL: %w", err)
	}

	return nil
}

func (s *cinemaService) AddFacility(ctx context.Context, cinemaID int64, facilityID int64) error {
	return s.cinemaRepo.AddFacility(ctx, cinemaID, facilityID)
}

func (s *cinemaService) RemoveFacility(ctx context.Context, cinemaID int64, facilityID int64) error {
	return s.cinemaRepo.RemoveFacility(ctx, cinemaID, facilityID)
}

func (s *cinemaService) UpdateFacilities(ctx context.Context, cinemaID int64, facilityIDs []int64) error {
	return s.cinemaRepo.UpdateFacilities(ctx, cinemaID, facilityIDs)
}
