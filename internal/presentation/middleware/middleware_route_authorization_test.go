package middleware

import (
	"testing"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

func TestCriticalEndpointPolicies(t *testing.T) {
	t.Parallel()
	for _, route := range []string{"POST /v1/auth/otp/send", "POST /v1/auth/otp/verify", "GET /v1/movies", "GET /v1/cinemas/:id"} {
		if _, ok := publicEndpointPolicies[route]; !ok {
			t.Fatalf("expected %s to be public", route)
		}
	}
	for _, route := range []string{"GET /v1/cinemas/:id/halls", "GET /v1/halls/:id/seats", "GET /v1/showtimes", "GET /v1/showtimes/:id/seats"} {
		if _, ok := publicEndpointPolicies[route]; !ok {
			t.Fatalf("expected booking route %s to be public", route)
		}
	}

	testCases := map[string]enum.PermissionType{
		"POST /v1/movies":                   enum.MoviesCreate,
		"POST /v1/search/sync":              enum.SearchSync,
		"POST /v1/users/:id/wallet/balance": enum.WalletChangeBalance,
	}
	for route, expected := range testCases {
		policy, ok := protectedEndpointPolicies[route]
		if !ok || policy.permission != expected.String() {
			t.Fatalf("unexpected policy for %s", route)
		}
	}

	if !protectedEndpointPolicies["PUT /v1/cinemas/:id"].ownedCinema {
		t.Fatal("cinema updates must be ownership scoped")
	}
	if !protectedEndpointPolicies["POST /v1/halls/:id/showtimes"].ownedHall {
		t.Fatal("showtime creation must be hall ownership scoped")
	}
	if !protectedEndpointPolicies["PUT /v1/seats/:id"].ownedSeat {
		t.Fatal("seat updates must be ownership scoped")
	}
	if !protectedEndpointPolicies["PUT /v1/showtimes/:id/prices"].ownedShowtime {
		t.Fatal("showtime pricing must be ownership scoped")
	}
	if protectedEndpointPolicies["POST /v1/reservations"].permission != enum.ReservationsCreate.String() {
		t.Fatal("reservation creation permission mismatch")
	}
	if protectedEndpointPolicies["GET /v1/users/:id/profile"].selfPermission != enum.UsersReadProfile.String() {
		t.Fatal("profile reads must be self scoped")
	}
}
