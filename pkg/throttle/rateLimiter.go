package throttle

import (
	"sync"

	"golang.org/x/time/rate"
)

// Implements a threadsafe per host rate limiter
type HostRateLimiter struct {
	limiters          map[string]*rate.Limiter // host to rate.Limiter mapping
	lock              *sync.RWMutex
	defaultTps        float64
	defaultBucketSize int
}

func NewHostRateLimiter(defaultTps float64) *HostRateLimiter {
	return &HostRateLimiter{
		limiters:          make(map[string]*rate.Limiter),
		lock:              &sync.RWMutex{},
		defaultTps:        defaultTps,
		defaultBucketSize: 1,
	}
}

// Gets the rate limiter for a host. If none exists for the host, a new one is created.
// This call requires mutex lock acquisition
func (rl *HostRateLimiter) GetLimiter(host string) *rate.Limiter {

	rl.lock.RLock()
	limiter := rl.limiters[host]
	rl.lock.RUnlock()
	if limiter != nil {
		return limiter
	}

	// Create new limiter for unseen host
	limiter = rate.NewLimiter(rate.Limit(rl.defaultTps), rl.defaultBucketSize)
	rl.lock.Lock()
	defer rl.lock.Unlock()

	rl.limiters[host] = limiter
	return limiter
}
