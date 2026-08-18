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
	person := &entity.Person{
		Name:      req.Name,
		Biography: req.Biography,
	}

	createdPerson, err := s.personRepo.Create(ctx, person)
	if err != nil {
		return nil, err
	}

	return dto.ToPersonDetailResponse(createdPerson), nil
}

func (s *personService) GetByID(ctx context.Context, id int64) (*dto.PersonDetailResponse, error) {
	person, err := s.personRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.ToPersonDetailResponse(person), nil
}

func (s *personService) GetAll(ctx context.Context) ([]*dto.PersonListResponse, error) {
	persons, err := s.personRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return dto.ToPersonListResponseSlice(persons), nil
}

func (s *personService) Update(ctx context.Context, id int64, req *dto.UpdatePersonRequest) (*dto.PersonDetailResponse, error) {
	person, err := s.personRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		person.Name = *req.Name
	}

	if req.Biography != nil {
		person.Biography = *req.Biography
	}

	updatedPerson, err := s.personRepo.Update(ctx, person)
	if err != nil {
		return nil, err
	}

	return dto.ToPersonDetailResponse(updatedPerson), nil
}

func (s *personService) Delete(ctx context.Context, id int64) error {
	person, err := s.personRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if person.ImageURLLow != "" || person.ImageURLHigh != "" {
		_ = s.DeleteImage(ctx, id)
	}

	return s.personRepo.Delete(ctx, id)
}

func (s *personService) UploadImage(ctx context.Context, id int64, file io.Reader) (*dto.PersonDetailResponse, error) {
	person, err := s.personRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	img, err := image.DecodeImage(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	lowImg := image.Resize(200, 200, img)
	lowBuf, err := image.EncodeToBuffer(lowImg, image.FormatWebP)
	if err != nil {
		return nil, fmt.Errorf("failed to encode low-res image: %w", err)
	}
	lowKey := fmt.Sprintf("persons/person_%d_low.webp", id)
	lowURL, err := s.storage.UploadObject(ctx, enum.ActorPhoto, lowKey, lowBuf, image.GetContentType(image.FormatWebP))
	if err != nil {
		return nil, fmt.Errorf("failed to upload low-res image: %w", err)
	}

	highImg := image.Resize(800, 800, img)
	highBuf, err := image.EncodeToBuffer(highImg, image.FormatJPEG)
	if err != nil {
		return nil, fmt.Errorf("failed to encode high-res image: %w", err)
	}
	highKey := fmt.Sprintf("persons/person_%d_high.jpg", id)
	highURL, err := s.storage.UploadObject(ctx, enum.ActorPhoto, highKey, highBuf, image.GetContentType(image.FormatJPEG))
	if err != nil {
		return nil, fmt.Errorf("failed to upload high-res image: %w", err)
	}

	person.ImageURLLow = lowURL
	person.ImageURLHigh = highURL

	updatedPerson, err := s.personRepo.Update(ctx, person)
	if err != nil {
		return nil, err
	}

	return dto.ToPersonDetailResponse(updatedPerson), nil
}

func (s *personService) DeleteImage(ctx context.Context, id int64) error {
	person, err := s.personRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	lowKey := fmt.Sprintf("persons/person_%d_low.webp", id)
	highKey := fmt.Sprintf("persons/person_%d_high.jpg", id)

	_ = s.storage.DeleteObject(ctx, enum.ActorPhoto, lowKey)
	_ = s.storage.DeleteObject(ctx, enum.ActorPhoto, highKey)

	person.ImageURLLow = ""
	person.ImageURLHigh = ""

	_, err = s.personRepo.Update(ctx, person)
	if err != nil {
		return err
	}

	return nil
}
