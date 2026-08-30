package service

import (
	"context"
	"fmt"
	"io"
	"strings"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/broker"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/storage"
	"github.com/mehdi1443/cinema-booking-system/pkg/image"
)

const maxPersonImageSize = 5 * 1024 * 1024

type personService struct {
	personRepo repository_contract.PersonInterface
	storage    storage.StorageInterFace
	log        domainlogger.Logger
	searchPub  broker.SearchEventPublisher
}

func NewPersonService(
	personRepo repository_contract.PersonInterface,
	storage storage.StorageInterFace,
	log domainlogger.Logger,
	searchPub broker.SearchEventPublisher,
) service_contract.PersonServiceInterface {
	return &personService{
		personRepo: personRepo,
		storage:    storage,
		log:        log,
		searchPub:  searchPub,
	}
}

func (s *personService) publishIndexEvent(ctx context.Context, person *entity.Person) {
	if s.searchPub == nil || person == nil {
		return
	}

	doc := entity.SearchDocument{
		ID:       fmt.Sprintf("person_%d", person.ID),
		EntityID: person.ID,
		Type:     enum.DocTypePerson,
		Title:    person.Name,
		Subtitle: person.Biography,
		ImageURL: person.ImageURLLow,
	}

	if err := s.searchPub.PublishSyncEvent(ctx, doc); err != nil {
		s.log.Error(ctx, "Failed to publish search sync event", domainlogger.Error(err), domainlogger.Int64("person_id", person.ID))
	}
}

func (s *personService) publishDeleteEvent(ctx context.Context, personID int64) {
	docID := fmt.Sprintf("person_%d", personID)
	if err := s.searchPub.PublishDeleteEvent(ctx, docID); err != nil {
		s.log.Error(ctx, "Failed to publish search delete event", domainlogger.Error(err), domainlogger.Int64("person_id", personID))
	}
}

func (s *personService) SyncAllToSearch(ctx context.Context) error {
	limit := 100
	page := 1

	for {
		opts := repository_contract.NewQueryOptions().WithPagination(limit, (page-1)*limit)
		persons, _, err := s.personRepo.GetPersons(ctx, repository_contract.PersonFilter{}, opts)

		if err != nil {
			return fmt.Errorf("failed to fetch persons for sync: %w", err)
		}

		if len(persons) == 0 {
			break
		}

		for i := range persons {
			s.publishIndexEvent(ctx, &persons[i])
		}

		page++
	}

	s.log.Info(ctx, "All existing persons sent to RabbitMQ for search sync")
	return nil
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

	s.publishIndexEvent(ctx, createdPerson)

	s.log.Info(ctx, "Person created successfully", domainlogger.Int64("person_id", createdPerson.ID))

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

	s.publishIndexEvent(ctx, updatedPerson)

	s.log.Info(ctx, "Person updated successfully", domainlogger.Int64("person_id", updatedPerson.ID))

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

	s.publishDeleteEvent(ctx, id)

	s.log.Info(ctx, "Person deleted successfully", domainlogger.Int64("person_id", id))

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

	s.publishIndexEvent(ctx, updatedPerson)

	s.log.Info(ctx, "Person image uploaded successfully", domainlogger.Int64("person_id", id))

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

	if err := s.storage.DeleteObject(ctx, enum.ActorPhoto, lowKey); err != nil {
		s.log.Warn(ctx, "Failed to delete low-res person image from storage", domainlogger.Int64("person_id", id), domainlogger.Error(err))
	}

	if err := s.storage.DeleteObject(ctx, enum.ActorPhoto, highKey); err != nil {
		s.log.Warn(ctx, "Failed to delete high-res person image from storage", domainlogger.Int64("person_id", id), domainlogger.Error(err))
	}

	person.ImageURLLow = ""
	person.ImageURLHigh = ""

	updatedPerson, err := s.personRepo.Update(ctx, person)
	if err != nil {
		return fmt.Errorf("remove image URLs for person %d: %w", id, err)
	}

	s.publishIndexEvent(ctx, updatedPerson)

	s.log.Info(ctx, "Person image deleted successfully", domainlogger.Int64("person_id", id))

	return nil
}
