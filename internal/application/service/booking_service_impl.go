package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
)

type hallService struct {
	repo       DB.HallRepositoryInterface
	cinemaRepo DB.CinemaInterface
	log        domainlogger.Logger
}

func NewHallService(repo DB.HallRepositoryInterface, cinemaRepo DB.CinemaInterface, log domainlogger.Logger) service_contract.HallServiceInterface {
	return &hallService{repo: repo, cinemaRepo: cinemaRepo, log: log}
}

func (s *hallService) CreateHall(ctx context.Context, cinemaID int64, req *dto.CreateHallRequest) (*dto.HallResponse, error) {
	if req == nil || cinemaID <= 0 || !req.Status.IsValid() {
		return nil, exception.ErrInvalidInput
	}
	if _, err := s.cinemaRepo.GetByID(ctx, cinemaID); err != nil {
		return nil, err
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, exception.ErrInvalidInput
	}
	hall, err := s.repo.CreateHall(ctx, &entity.Hall{CinemaID: cinemaID, Name: name, Status: req.Status})
	if err != nil {
		return nil, err
	}
	s.log.Info(ctx, "Hall created successfully", domainlogger.Int64("hall_id", hall.ID), domainlogger.Int64("cinema_id", cinemaID))
	return dto.ToHallResponse(hall), nil
}

func (s *hallService) GetHall(ctx context.Context, id int64) (*dto.HallResponse, error) {
	hall, err := s.repo.GetHall(ctx, id)
	if err != nil {
		return nil, err
	}
	return dto.ToHallResponse(hall), nil
}

func (s *hallService) ListHalls(ctx context.Context, cinemaID int64) ([]*dto.HallResponse, error) {
	if _, err := s.cinemaRepo.GetByID(ctx, cinemaID); err != nil {
		return nil, err
	}
	halls, err := s.repo.ListHalls(ctx, cinemaID)
	if err != nil {
		return nil, err
	}
	result := make([]*dto.HallResponse, 0, len(halls))
	for i := range halls {
		result = append(result, dto.ToHallResponse(&halls[i]))
	}
	return result, nil
}

func (s *hallService) UpdateHall(ctx context.Context, id int64, req *dto.UpdateHallRequest) (*dto.HallResponse, error) {
	if req == nil {
		return nil, exception.ErrInvalidInput
	}
	hall, err := s.repo.GetHall(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, exception.ErrInvalidInput
		}
		hall.Name = name
	}
	if req.Status != nil {
		if !req.Status.IsValid() {
			return nil, exception.ErrInvalidStatus
		}
		hall.Status = *req.Status
	}
	return s.updateHall(ctx, hall)
}

func (s *hallService) updateHall(ctx context.Context, hall *entity.Hall) (*dto.HallResponse, error) {
	updated, err := s.repo.UpdateHall(ctx, hall)
	if err != nil {
		return nil, err
	}
	s.log.Info(ctx, "Hall updated successfully", domainlogger.Int64("hall_id", updated.ID))
	return dto.ToHallResponse(updated), nil
}

func (s *hallService) DeleteHall(ctx context.Context, id int64) error {
	if err := s.repo.DeleteHall(ctx, id); err != nil {
		return err
	}
	s.log.Info(ctx, "Hall deleted successfully", domainlogger.Int64("hall_id", id))
	return nil
}

func (s *hallService) CreateSeats(ctx context.Context, hallID int64, req *dto.CreateSeatsRequest) ([]dto.SeatResponse, error) {
	if req == nil || len(req.Seats) == 0 || len(req.Seats) > 500 {
		return nil, exception.ErrInvalidSeatSelection
	}
	if _, err := s.repo.GetHall(ctx, hallID); err != nil {
		return nil, err
	}
	seen := make(map[string]struct{}, len(req.Seats))
	seats := make([]entity.Seat, 0, len(req.Seats))
	for _, input := range req.Seats {
		row := strings.ToUpper(strings.TrimSpace(input.RowLabel))
		if row == "" || input.Number <= 0 || !input.Type.IsValid() || !input.Status.IsValid() {
			return nil, exception.ErrInvalidSeatSelection
		}
		key := row + ":" + strconv.Itoa(input.Number)
		if _, exists := seen[key]; exists {
			return nil, exception.ErrInvalidSeatSelection
		}
		seen[key] = struct{}{}
		seats = append(seats, entity.Seat{HallID: hallID, RowLabel: row, Number: input.Number, Type: input.Type, Status: input.Status})
	}
	created, err := s.repo.CreateSeats(ctx, hallID, seats)
	if err != nil {
		return nil, err
	}
	s.log.Info(ctx, "Seats created successfully", domainlogger.Int64("hall_id", hallID), domainlogger.Int("seat_count", len(created)))
	return dto.ToSeatResponses(created), nil
}

func (s *hallService) ListSeats(ctx context.Context, hallID int64) ([]dto.SeatResponse, error) {
	if _, err := s.repo.GetHall(ctx, hallID); err != nil {
		return nil, err
	}
	seats, err := s.repo.ListSeats(ctx, hallID)
	if err != nil {
		return nil, err
	}
	return dto.ToSeatResponses(seats), nil
}

func (s *hallService) UpdateSeat(ctx context.Context, id int64, req *dto.UpdateSeatRequest) (*dto.SeatResponse, error) {
	if req == nil {
		return nil, exception.ErrInvalidInput
	}
	seat, err := s.repo.GetSeat(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.RowLabel != nil {
		seat.RowLabel = strings.ToUpper(strings.TrimSpace(*req.RowLabel))
		if seat.RowLabel == "" {
			return nil, exception.ErrInvalidInput
		}
	}
	if req.Number != nil {
		if *req.Number <= 0 || *req.Number > 1000 {
			return nil, exception.ErrInvalidInput
		}
		seat.Number = *req.Number
	}
	if req.Type != nil {
		if !req.Type.IsValid() {
			return nil, exception.ErrInvalidInput
		}
		seat.Type = *req.Type
	}
	if req.Status != nil {
		if !req.Status.IsValid() {
			return nil, exception.ErrInvalidStatus
		}
		seat.Status = *req.Status
	}
	updated, err := s.repo.UpdateSeat(ctx, seat)
	if err != nil {
		return nil, err
	}
	s.log.Info(ctx, "Seat updated successfully", domainlogger.Int64("seat_id", updated.ID), domainlogger.Int64("hall_id", updated.HallID))
	response := dto.ToSeatResponse(updated)
	return &response, nil
}

func (s *hallService) DeleteSeat(ctx context.Context, id int64) error {
	if err := s.repo.DeleteSeat(ctx, id); err != nil {
		return err
	}
	s.log.Info(ctx, "Seat deleted successfully", domainlogger.Int64("seat_id", id))
	return nil
}

type showtimeService struct {
	repo      DB.ShowtimeRepositoryInterface
	hallRepo  DB.HallRepositoryInterface
	movieRepo DB.MovieInterFace
	expirer   DB.ReservationRepositoryInterface
	log       domainlogger.Logger
}

func NewShowtimeService(repo DB.ShowtimeRepositoryInterface, hallRepo DB.HallRepositoryInterface, movieRepo DB.MovieInterFace, expirer DB.ReservationRepositoryInterface, log domainlogger.Logger) service_contract.ShowtimeServiceInterface {
	return &showtimeService{repo: repo, hallRepo: hallRepo, movieRepo: movieRepo, expirer: expirer, log: log}
}

func (s *showtimeService) Create(ctx context.Context, hallID int64, req *dto.CreateShowtimeRequest) (*dto.ShowtimeResponse, error) {
	if req == nil || req.BasePrice <= 0 || (req.Status != enum.ShowtimeStatusDraft && req.Status != enum.ShowtimeStatusPublished) || !req.StartsAt.After(time.Now()) {
		return nil, exception.ErrInvalidInput
	}
	hall, err := s.hallRepo.GetHall(ctx, hallID)
	if err != nil {
		return nil, err
	}
	if hall.Status != enum.HallStatusActive {
		return nil, exception.ErrShowtimeNotBookable
	}
	movie, err := s.movieRepo.GetByID(ctx, req.MovieID)
	if err != nil {
		return nil, err
	}
	if movie.DurationMinutes <= 0 || movie.Status == enum.MovieInactive || movie.Status == enum.MovieEnded {
		return nil, exception.ErrShowtimeNotBookable
	}
	prices, err := buildSeatTypePrices(req.BasePrice, req.TypePrices)
	if err != nil {
		return nil, err
	}
	showtime := &entity.Showtime{HallID: hallID, MovieID: req.MovieID, StartsAt: req.StartsAt.UTC(), EndsAt: req.StartsAt.UTC().Add(time.Duration(movie.DurationMinutes) * time.Minute), Status: req.Status}
	created, err := s.repo.CreateShowtime(ctx, showtime, prices)
	if err != nil {
		return nil, err
	}
	s.log.Info(ctx, "Showtime created successfully", domainlogger.Int64("showtime_id", created.ID), domainlogger.Int64("hall_id", hallID), domainlogger.Int64("movie_id", req.MovieID))
	return dto.ToShowtimeResponse(created), nil
}

func buildSeatTypePrices(basePrice int64, overrides []dto.SeatTypePriceRequest) (map[enum.SeatType]int64, error) {
	if basePrice <= 0 {
		return nil, exception.ErrInvalidAmount
	}
	prices := map[enum.SeatType]int64{enum.SeatTypeStandard: basePrice, enum.SeatTypeVIP: basePrice, enum.SeatTypeWheelchair: basePrice}
	seen := make(map[enum.SeatType]struct{}, len(overrides))
	for _, item := range overrides {
		if !item.SeatType.IsValid() || item.Price <= 0 {
			return nil, exception.ErrInvalidAmount
		}
		if _, exists := seen[item.SeatType]; exists {
			return nil, exception.ErrInvalidInput
		}
		seen[item.SeatType] = struct{}{}
		prices[item.SeatType] = item.Price
	}
	return prices, nil
}

func (s *showtimeService) Get(ctx context.Context, id int64) (*dto.ShowtimeResponse, error) {
	showtime, err := s.repo.GetShowtime(ctx, id)
	if err != nil {
		return nil, err
	}
	if !isPublicShowtime(showtime, time.Now()) {
		return nil, exception.ErrShowtimeNotFound
	}
	return dto.ToShowtimeResponse(showtime), nil
}

func isPublicShowtime(showtime *entity.Showtime, now time.Time) bool {
	return showtime != nil &&
		showtime.Status == enum.ShowtimeStatusPublished &&
		showtime.StartsAt.After(now) &&
		showtime.Hall != nil && showtime.Hall.Status == enum.HallStatusActive &&
		showtime.Hall.Cinema != nil && showtime.Hall.Cinema.Status == enum.CinemaStatusActive &&
		showtime.Movie != nil && showtime.Movie.Status != enum.MovieInactive && showtime.Movie.Status != enum.MovieEnded
}

func (s *showtimeService) List(ctx context.Context, req *dto.ShowtimeFilterRequest) ([]*dto.ShowtimeResponse, int64, error) {
	if req == nil {
		return nil, 0, exception.ErrInvalidFilter
	}
	req.Normalize()
	if req.From != nil && req.To != nil && !req.To.After(*req.From) {
		return nil, 0, exception.ErrInvalidFilter
	}
	now := time.Now()
	from := req.From
	if from == nil || from.Before(now) {
		from = &now
	}
	published := enum.ShowtimeStatusPublished
	filter := DB.ShowtimeFilter{
		CinemaID:  req.CinemaID,
		HallID:    req.HallID,
		MovieID:   req.MovieID,
		Status:    &published,
		From:      from,
		To:        req.To,
		PublicOnly: true,
	}
	items, total, err := s.repo.ListShowtimes(ctx, filter, DB.NewQueryOptions().WithPagination(req.Limit, req.GetOffset()))
	if err != nil {
		return nil, 0, err
	}
	result := make([]*dto.ShowtimeResponse, 0, len(items))
	for i := range items {
		result = append(result, dto.ToShowtimeResponse(&items[i]))
	}
	return result, total, nil
}

func (s *showtimeService) Update(ctx context.Context, id int64, req *dto.UpdateShowtimeRequest) (*dto.ShowtimeResponse, error) {
	if req == nil {
		return nil, exception.ErrInvalidInput
	}
	showtime, err := s.repo.GetShowtime(ctx, id)
	if err != nil {
		return nil, err
	}
	if showtime.Status == enum.ShowtimeStatusCancelled || showtime.Status == enum.ShowtimeStatusCompleted || !showtime.StartsAt.After(time.Now()) {
		return nil, exception.ErrResourceInUse
	}
	if req.StartsAt != nil {
		if !req.StartsAt.After(time.Now()) {
			return nil, exception.ErrInvalidInput
		}
		movie, err := s.movieRepo.GetByID(ctx, showtime.MovieID)
		if err != nil {
			return nil, err
		}
		showtime.StartsAt = req.StartsAt.UTC()
		showtime.EndsAt = showtime.StartsAt.Add(time.Duration(movie.DurationMinutes) * time.Minute)
	}
	if req.Status != nil {
		if !req.Status.IsValid() {
			return nil, exception.ErrInvalidStatus
		}
		showtime.Status = *req.Status
	}
	updated, err := s.repo.UpdateShowtime(ctx, showtime)
	if err != nil {
		return nil, err
	}
	s.log.Info(ctx, "Showtime updated successfully", domainlogger.Int64("showtime_id", updated.ID))
	return dto.ToShowtimeResponse(updated), nil
}

func (s *showtimeService) Delete(ctx context.Context, id int64) error {
	if err := s.repo.DeleteShowtime(ctx, id); err != nil {
		return err
	}
	s.log.Info(ctx, "Showtime deleted successfully", domainlogger.Int64("showtime_id", id))
	return nil
}

func (s *showtimeService) UpdatePrices(ctx context.Context, id int64, req *dto.UpdateShowtimePricesRequest) error {
	if req == nil || len(req.Seats) == 0 {
		return exception.ErrInvalidInput
	}
	prices := make(map[int64]int64, len(req.Seats))
	for _, item := range req.Seats {
		if item.ShowtimeSeatID <= 0 || item.Price <= 0 {
			return exception.ErrInvalidAmount
		}
		if _, exists := prices[item.ShowtimeSeatID]; exists {
			return exception.ErrInvalidSeatSelection
		}
		prices[item.ShowtimeSeatID] = item.Price
	}
	if err := s.repo.UpdateShowtimePrices(ctx, id, prices); err != nil {
		return err
	}
	s.log.Info(ctx, "Showtime prices updated successfully", domainlogger.Int64("showtime_id", id), domainlogger.Int("seat_count", len(prices)))
	return nil
}

func (s *showtimeService) ListSeats(ctx context.Context, id int64) ([]dto.ShowtimeSeatResponse, error) {
	if err := s.expirer.ExpireReservations(ctx, time.Now()); err != nil {
		return nil, err
	}
	showtime, err := s.repo.GetShowtime(ctx, id)
	if err != nil {
		return nil, err
	}
	if !isPublicShowtime(showtime, time.Now()) {
		return nil, exception.ErrShowtimeNotFound
	}
	seats, err := s.repo.ListShowtimeSeats(ctx, id)
	if err != nil {
		return nil, err
	}
	return dto.ToShowtimeSeatResponses(seats), nil
}

type reservationService struct {
	repo   DB.ReservationRepositoryInterface
	config *bootstrap.ReservationConfig
	log    domainlogger.Logger
}

func NewReservationService(repo DB.ReservationRepositoryInterface, config *bootstrap.ReservationConfig, log domainlogger.Logger) service_contract.ReservationServiceInterface {
	return &reservationService{repo: repo, config: config, log: log}
}

func (s *reservationService) Create(ctx context.Context, userID int64, req *dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	if req == nil || userID <= 0 || len(req.ShowtimeSeatIDs) == 0 || len(req.ShowtimeSeatIDs) > 10 {
		return nil, exception.ErrInvalidSeatSelection
	}
	ids, ok := uniquePositiveIDs(req.ShowtimeSeatIDs)
	if !ok {
		return nil, exception.ErrInvalidSeatSelection
	}
	if err := s.repo.ExpireReservations(ctx, time.Now()); err != nil {
		return nil, err
	}
	reservation, err := s.repo.CreateReservation(ctx, userID, req.ShowtimeID, ids, time.Now().Add(s.config.HoldDuration()))
	if err != nil {
		return nil, err
	}
	s.log.Info(ctx, "Reservation created successfully", domainlogger.Int64("reservation_id", reservation.ID), domainlogger.Int64("user_id", userID), domainlogger.Int64("showtime_id", req.ShowtimeID), domainlogger.Int("seat_count", len(ids)))
	return dto.ToReservationResponse(reservation), nil
}

func uniquePositiveIDs(input []int64) ([]int64, bool) {
	seen := make(map[int64]struct{}, len(input))
	result := make([]int64, 0, len(input))
	for _, id := range input {
		if id <= 0 {
			return nil, false
		}
		if _, exists := seen[id]; exists {
			return nil, false
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result, true
}

func (s *reservationService) Get(ctx context.Context, id, userID int64, allowAny bool) (*dto.ReservationResponse, error) {
	if err := s.repo.ExpireReservations(ctx, time.Now()); err != nil {
		return nil, err
	}
	reservation, err := s.repo.GetReservation(ctx, id)
	if err != nil {
		return nil, err
	}
	if !allowAny && reservation.UserID != userID {
		return nil, exception.ErrNoteAccessible
	}
	return dto.ToReservationResponse(reservation), nil
}

func (s *reservationService) ListMine(ctx context.Context, userID int64) ([]*dto.ReservationResponse, error) {
	if err := s.repo.ExpireReservations(ctx, time.Now()); err != nil {
		return nil, err
	}
	reservations, err := s.repo.ListUserReservations(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]*dto.ReservationResponse, 0, len(reservations))
	for i := range reservations {
		result = append(result, dto.ToReservationResponse(&reservations[i]))
	}
	return result, nil
}

func (s *reservationService) Cancel(ctx context.Context, id, userID int64, allowAny bool) (*dto.ReservationResponse, error) {
	reservation, err := s.repo.GetReservation(ctx, id)
	if err != nil {
		return nil, err
	}
	if !allowAny && reservation.UserID != userID {
		return nil, exception.ErrNoteAccessible
	}
	cancelled, err := s.repo.CancelReservation(ctx, id)
	if err != nil {
		return nil, err
	}
	s.log.Info(ctx, "Reservation cancelled successfully", domainlogger.Int64("reservation_id", cancelled.ID), domainlogger.Int64("user_id", userID), domainlogger.String("status", string(cancelled.Status)))
	return dto.ToReservationResponse(cancelled), nil
}

func (s *reservationService) ExpirePending(ctx context.Context, now time.Time) error {
	if err := s.repo.ExpireReservations(ctx, now); err != nil {
		return err
	}
	s.log.Info(ctx, "Pending reservations expiration completed")
	return nil
}
