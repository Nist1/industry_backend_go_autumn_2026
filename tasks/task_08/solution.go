package main

import (
	"sync"
	"time"
)

type Clock interface{ Now() time.Time }

type Limiter struct {
	mu     sync.Mutex
	clock  Clock
	rate   float64
	burst  int
	tokens float64
	last   time.Time
}

func NewLimiter(clock Clock, ratePerSec float64, burst int) *Limiter {
	l := &Limiter{
		clock:  clock,
		rate:   ratePerSec,
		burst:  burst,
		tokens: float64(max(burst, 0)), // при создании корзина полная
	}

	if clock != nil {
		l.last = clock.Now()
	}
	return l
}

func (l *Limiter) AllowN(n int) bool {
	if l.clock == nil || l.burst <= 0 || n <= 0 || n > l.burst {
		return false
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock.Now()
	if now.After(l.last) {
		if l.rate > 0 {
			elapsed := now.Sub(l.last).Seconds()
			// ограничиваем burst: простой не копит токены
			l.tokens = min(float64(l.burst), l.tokens+elapsed*l.rate)
		}
		l.last = now
	}

	need := float64(n)
	if l.tokens < need {
		// Токенов не хватает -> возвращаем false и ничего не списываем
		return false
	}
	l.tokens -= need
	return true
}
