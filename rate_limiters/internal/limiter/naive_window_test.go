package limiter

import (
	"testing"
)

func TestNaiveWindowLimiter_AllowsWithinLimit(t *testing.T) {
	clock := &FakeClock{}

	limiter := NewNaiveWindowLimiter(
		3,
		10,
		clock,
	)

	if !limiter.Allow("alice") {
		t.Fatal("expected first request allowed")
	}

	if !limiter.Allow("alice") {
		t.Fatal("expected second request allowed")
	}

	if !limiter.Allow("alice") {
		t.Fatal("expected third request allowed")
	}
}

func TestNaiveWindowLimiter_WindowExpires(t *testing.T) {
	clock := &FakeClock{}

	limiter := NewNaiveWindowLimiter(
		2,
		10,
		clock,
	)

	if !limiter.Allow("alice") {
		t.Fatal("expected allowed")
	}

	if !limiter.Allow("alice") {
		t.Fatal("expected allowed")
	}

	if limiter.Allow("alice") {
		t.Fatal("expected rejected ")
	}

	clock.Advance(10)

	if !limiter.Allow("alice") {
		t.Fatal("expected allowed after expiration")
	}
}
