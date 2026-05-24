package limiter

import (
	"sync"

	"github.com/vshishkar/rate-limiter-lab/internal"
)

type NaiveWindowLimiter struct {
	records map[string]*NaiveWindow
	mu      sync.Mutex

	MaxRequests int
	Clock       internal.Clock
	WindowSize  int64
}

func NewNaiveWindowLimiter(maxRequests int, windowSize int64, clock internal.Clock) *NaiveWindowLimiter {
	return &NaiveWindowLimiter{
		records:     make(map[string]*NaiveWindow),
		MaxRequests: maxRequests,
		Clock:       clock,
		WindowSize:  windowSize,
	}
}

func (w *NaiveWindowLimiter) Allow(key string) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	window, exists := w.records[key]
	if !exists {
		w.records[key] = &NaiveWindow{Count: 1, Timestamp: w.Clock.Now()}
		return true
	}

	// lazy eviction
	if w.Clock.Now()-window.Timestamp >= w.WindowSize {
		window.Count = 1
		window.Timestamp = w.Clock.Now()
		return true
	}

	window.Count++
	return window.Count <= w.MaxRequests
}

type NaiveWindow struct {
	Count     int
	Timestamp int64
}
