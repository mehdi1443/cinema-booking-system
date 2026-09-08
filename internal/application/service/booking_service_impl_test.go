package service

import (
	"testing"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
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
