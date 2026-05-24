package limiter

import (
	"container/list"
	"sync"

	"github.com/vshishkar/rate-limiter-lab/internal"
)

type SlidingWindowLimiter struct {
	records     map[string]*SlidingWindow
	MaxRequests int
	Clock       internal.Clock
	WindowSize  int64
}

func NewSlidingWindowLimiter(maxRequests int, windowSize int64, clock internal.Clock) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		records:     make(map[string]*SlidingWindow),
		MaxRequests: maxRequests,
		Clock:       clock,
		WindowSize:  windowSize,
	}
}

func (s *SlidingWindowLimiter) Allow(key string) bool {
	now := s.Clock.Now()

	window, exists := s.records[key]
	if !exists {
		s.records[key] = &SlidingWindow{Timestamps: list.New()}
		window = s.records[key]
	}

	window.mu.Lock()
	defer window.mu.Unlock()

	for window.Timestamps.Len() > 0 && window.Timestamps.Front().Value.(int64) <= now-s.WindowSize {
		window.Timestamps.Remove(window.Timestamps.Front())
	}

	if window.Timestamps.Len() < s.MaxRequests {
		window.Timestamps.PushBack(now)
		return true
	}

	return false
}

type SlidingWindow struct {
	Timestamps *list.List
	mu         sync.Mutex
}
