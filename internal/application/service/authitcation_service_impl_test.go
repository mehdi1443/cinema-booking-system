package service

import "testing"

func TestNormalizeIranianPhone(t *testing.T) {
	t.Parallel()
	testCases := map[string]string{
		"09123456789":   "09123456789",
		"9123456789":    "09123456789",
		"+989123456789": "09123456789",
	}
	for input, expected := range testCases {
		actual, err := normalizeIranianPhone(input)
		if err != nil {
			t.Fatalf("normalize %q: %v", input, err)
		}
		if actual != expected {
			t.Fatalf("normalize %q: got %q, want %q", input, actual, expected)
		}
	}

	for _, input := range []string{"", "123", "08123456789", "09123abc789"} {
		if _, err := normalizeIranianPhone(input); err == nil {
			t.Fatalf("expected %q to be rejected", input)
		}
	}
}
