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
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/storage"
	"github.com/mehdi1443/cinema-booking-system/pkg/image"
)

const maxPersonImageSize = 5 * 1024 * 1024 // 5 MB

type personService struct {
	personRepo repository_contract.PersonInterface
	storage    storage.StorageInterFace
}

func NewPersonService(
	personRepo repository_contract.PersonInterface,
	storage storage.StorageInterFace,
) service_contract.PersonServiceInterface {
	return &personService{
		personRepo: personRepo,
		storage:    storage,
	}
}

func (s *personService) Create(ctx context.Context, req *dto.CreatePersonRequest) (*dto.PersonDetailResponse, error) {
	if req == nil {
		return nil, exception.ErrInvalidInput
	}

	trimmedName := strings.TrimSpace(req.Name)
	if trimmedName == "" {
		return nil, exception.ErrInvalidInput
	}

	person := &entity.Person{
		Name:      trimmedName,
		Biography: req.Biography,
	}

	createdPerson, err := s.personRepo.Create(ctx, person)
	if err != nil {
		return nil, fmt.Errorf("create person: %w", err)
	}

	return dto.ToPersonDetailResponse(createdPerson), nil
}

func (s *personService) GetByID(ctx context.Context, id int64) (*dto.PersonDetailResponse, error) {
	if id <= 0 {
		return nil, exception.ErrInvalidInput
	}

	person, err := s.personRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get person by id %d: %w", id, err)
	}

	return dto.ToPersonDetailResponse(person), nil
}

func (s *personService) GetPersons(ctx context.Context, req dto.PersonFilterRequest) ([]*dto.PersonListResponse, int64, error) {
	filter := repository_contract.PersonFilter{}
	if req.Search != nil {
		filter.SearchKeyword = strings.TrimSpace(*req.Search)
	}

	req.Normalize()

	opts := repository_contract.NewQueryOptions()
	opts.WithPagination(req.Limit, req.GetOffset())

	persons, total, err := s.personRepo.GetPersons(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("get persons: %w", err)
	}

	return dto.ToPersonListResponseSlice(persons), total, nil
}

func (s *personService) Update(ctx context.Context, id int64, req *dto.UpdatePersonRequest) (*dto.PersonDetailResponse, error) {
	if id <= 0 || req == nil {
		return nil, exception.ErrInvalidInput
	}

	person, err := s.personRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get person before update %d: %w", id, err)
	}

	if req.Name != nil {
		trimmedName := strings.TrimSpace(*req.Name)
		if trimmedName == "" {
			return nil, exception.ErrInvalidInput
		}
		person.Name = trimmedName
	}

	if req.Biography != nil {
		person.Biography = *req.Biography
	}

	updatedPerson, err := s.personRepo.Update(ctx, person)
	if err != nil {
		return nil, fmt.Errorf("update person %d: %w", id, err)
	}

	return dto.ToPersonDetailResponse(updatedPerson), nil
}

func (s *personService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return exception.ErrInvalidInput
	}

	person, err := s.personRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get person before delete %d: %w", id, err)
	}

	if person.ImageURLLow != "" || person.ImageURLHigh != "" {
		_ = s.DeleteImage(ctx, id)
	}

	if err := s.personRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete person %d: %w", id, err)
	}

	return nil
}

func (s *personService) UploadImage(ctx context.Context, id int64, file io.Reader) (*dto.PersonDetailResponse, error) {
	if id <= 0 {
		return nil, exception.ErrInvalidInput
	}

	if file == nil {
		return nil, exception.ErrInvalidFileType
	}

	person, err := s.personRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get person before image upload %d: %w", id, err)
	}

	fileReader, err := checkFileSize(file, maxPersonImageSize)
	if err != nil {
		return nil, err
	}

	img, err := image.DecodeImage(fileReader)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid person image", exception.ErrInvalidFileType)
	}

	lowImg := image.Resize(200, 200, img)
	lowBuf, err := image.EncodeToBuffer(lowImg, image.FormatWebP)
	if err != nil {
		return nil, fmt.Errorf("encode low resolution person image: %w", exception.ErrInternal)
	}

	lowKey := fmt.Sprintf("persons/person_%d_low.webp", id)
	lowURL, err := s.storage.UploadObject(ctx, enum.ActorPhoto, lowKey, lowBuf, image.GetContentType(image.FormatWebP))
	if err != nil {
		return nil, fmt.Errorf("upload low resolution person image: %w", exception.ErrInternal)
	}

	highImg := image.Resize(800, 800, img)
	highBuf, err := image.EncodeToBuffer(highImg, image.FormatJPEG)
	if err != nil {
		return nil, fmt.Errorf("encode high resolution person image: %w", exception.ErrInternal)
	}

	highKey := fmt.Sprintf("persons/person_%d_high.jpg", id)
	highURL, err := s.storage.UploadObject(ctx, enum.ActorPhoto, highKey, highBuf, image.GetContentType(image.FormatJPEG))
	if err != nil {
		return nil, fmt.Errorf("upload high resolution person image: %w", exception.ErrInternal)
	}

	person.ImageURLLow = lowURL
	person.ImageURLHigh = highURL

	updatedPerson, err := s.personRepo.Update(ctx, person)
	if err != nil {
		return nil, fmt.Errorf("save image URLs for person %d: %w", id, err)
	}

	return dto.ToPersonDetailResponse(updatedPerson), nil
}

func (s *personService) DeleteImage(ctx context.Context, id int64) error {
	if id <= 0 {
		return exception.ErrInvalidInput
	}

	person, err := s.personRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get person before image deletion %d: %w", id, err)
	}

	lowKey := fmt.Sprintf("persons/person_%d_low.webp", id)
	highKey := fmt.Sprintf("persons/person_%d_high.jpg", id)

	// Keep existing behavior:
	// storage deletion errors are intentionally ignored.
	_ = s.storage.DeleteObject(ctx, enum.ActorPhoto, lowKey)
	_ = s.storage.DeleteObject(ctx, enum.ActorPhoto, highKey)

	person.ImageURLLow = ""
	person.ImageURLHigh = ""

	_, err = s.personRepo.Update(ctx, person)
	if err != nil {
		return fmt.Errorf("remove image URLs for person %d: %w", id, err)
	}

	return nil
}
