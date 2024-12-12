package limiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	limits   map[string]*userLimit
	maxReqs  int           // иаксимальное количество запросов
	interval time.Duration // интервал времени
	mu       sync.Mutex
}

type userLimit struct {
	requests int       // количество запросов
	reset    time.Time // время сброса
}

// создаёт новый RateLimiter
func NewRateLimiter(maxReqs int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		limits:   make(map[string]*userLimit),
		maxReqs:  maxReqs,
		interval: interval,
	}
}

// проверяет, можно ли выполнить запрос от пользователя
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// получаем данные для пользователя
	limit, exists := rl.limits[ip]
	if !exists || time.Now().After(limit.reset) {
		// если лимит для пользователя не существует или его время сброса прошло, сбрасываем лимит
		rl.limits[ip] = &userLimit{
			requests: 1,
			reset:    time.Now().Add(rl.interval),
		}
		return true
	}

	// проверяем, превышен ли лимит
	if limit.requests >= rl.maxReqs {
		return false
	}

	// увеличиваем количество запросов
	limit.requests++
	return true
}
