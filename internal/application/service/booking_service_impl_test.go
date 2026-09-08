package service

import (
	"testing"
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

func TestBuildSeatTypePrices(t *testing.T) {
	t.Parallel()
	prices, err := buildSeatTypePrices(100_000, []dto.SeatTypePriceRequest{{SeatType: enum.SeatTypeVIP, Price: 175_000}})
	if err != nil {
		t.Fatalf("build prices: %v", err)
	}
	if prices[enum.SeatTypeStandard] != 100_000 || prices[enum.SeatTypeWheelchair] != 100_000 || prices[enum.SeatTypeVIP] != 175_000 {
		t.Fatalf("unexpected price map: %#v", prices)
	}
	if _, err := buildSeatTypePrices(100_000, []dto.SeatTypePriceRequest{{SeatType: enum.SeatTypeVIP, Price: 150_000}, {SeatType: enum.SeatTypeVIP, Price: 160_000}}); err == nil {
		t.Fatal("duplicate seat type overrides must fail")
	}
}

func TestUniquePositiveIDs(t *testing.T) {
	t.Parallel()
	ids, ok := uniquePositiveIDs([]int64{3, 1, 2})
	if !ok || len(ids) != 3 {
		t.Fatalf("valid IDs rejected: %#v", ids)
	}
	for _, input := range [][]int64{{1, 1}, {0}, {-1, 2}} {
		if _, ok := uniquePositiveIDs(input); ok {
			t.Fatalf("invalid selection accepted: %#v", input)
		}
	}
}

func TestIsPublicShowtime(t *testing.T) {
	now := time.Now()
	showtime := &entity.Showtime{
		StartsAt: now.Add(time.Hour),
		Status:   enum.ShowtimeStatusPublished,
		Hall: &entity.Hall{
			Status: enum.HallStatusActive,
			Cinema: &entity.Cinema{Status: enum.CinemaStatusActive},
		},
		Movie: &entity.Movie{Status: enum.MovieNowShowing},
	}

	if !isPublicShowtime(showtime, now) {
		t.Fatal("published upcoming showtime in active resources must be public")
	}

	checks := []struct {
		name   string
		mutate func(*entity.Showtime)
	}{
		{name: "draft", mutate: func(s *entity.Showtime) { s.Status = enum.ShowtimeStatusDraft }},
		{name: "past", mutate: func(s *entity.Showtime) { s.StartsAt = now.Add(-time.Minute) }},
		{name: "inactive hall", mutate: func(s *entity.Showtime) { s.Hall.Status = enum.HallStatusInactive }},
		{name: "inactive cinema", mutate: func(s *entity.Showtime) { s.Hall.Cinema.Status = enum.CinemaStatusInactive }},
		{name: "ended movie", mutate: func(s *entity.Showtime) { s.Movie.Status = enum.MovieEnded }},
	}

	for _, check := range checks {
		t.Run(check.name, func(t *testing.T) {
			candidate := *showtime
			hall := *showtime.Hall
			cinema := *showtime.Hall.Cinema
			movie := *showtime.Movie
			hall.Cinema = &cinema
			candidate.Hall = &hall
			candidate.Movie = &movie
			check.mutate(&candidate)
			if isPublicShowtime(&candidate, now) {
				t.Fatal("non-public showtime was exposed")
			}
		})
	}
}
