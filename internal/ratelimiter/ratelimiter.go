// Package ratelimiter provides a token-bucket style rate limiter for
// controlling how many log lines are emitted per second during output.
package ratelimiter

import (
	"context"
	"time"
)

// RateLimiter controls the throughput of log line emission.
type RateLimiter struct {
	ticker  *time.Ticker
	tokens  chan struct{}
	cancel  context.CancelFunc
	ctx     context.Context
	linesPS int
}

// New creates a RateLimiter that allows at most linesPerSecond tokens per
// second. If linesPerSecond is zero or negative, no limiting is applied and
// Wait always returns immediately.
func New(linesPerSecond int) *RateLimiter {
	ctx, cancel := context.WithCancel(context.Background())
	rl := &RateLimiter{
		cancel:  cancel,
		ctx:     ctx,
		linesPS: linesPerSecond,
	}
	if linesPerSecond > 0 {
		interval := time.Second / time.Duration(linesPerSecond)
		rl.ticker = time.NewTicker(interval)
		rl.tokens = make(chan struct{}, linesPerSecond)
		go rl.produce()
	}
	return rl
}

// produce fills the token channel on each tick.
func (rl *RateLimiter) produce() {
	for {
		select {
		case <-rl.ctx.Done():
			return
		case <-rl.ticker.C:
			select {
			case rl.tokens <- struct{}{}:
			default:
			}
		}
	}
}

// Wait blocks until a token is available or ctx is cancelled.
// If no rate limit is configured, Wait returns nil immediately.
func (rl *RateLimiter) Wait(ctx context.Context) error {
	if rl.linesPS <= 0 {
		return nil
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-rl.ctx.Done():
		return rl.ctx.Err()
	case <-rl.tokens:
		return nil
	}
}

// LinesPerSecond returns the configured rate limit.
func (rl *RateLimiter) LinesPerSecond() int { return rl.linesPS }

// Stop releases resources held by the RateLimiter.
func (rl *RateLimiter) Stop() {
	rl.cancel()
	if rl.ticker != nil {
		rl.ticker.Stop()
	}
}
