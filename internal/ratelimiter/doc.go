// Package ratelimiter implements a token-bucket rate limiter intended for
// use by the logslice output pipeline.
//
// When processing large log archives it can be useful to throttle how quickly
// matching lines are written so that downstream consumers (e.g. a pager or a
// streaming API) are not overwhelmed.
//
// Usage:
//
//	rl := ratelimiter.New(500) // 500 lines per second
//	defer rl.Stop()
//
//	for _, line := range lines {
//		if err := rl.Wait(ctx); err != nil {
//			break
//		}
//		fmt.Fprintln(w, line)
//	}
//
// A rate of zero or less disables limiting entirely.
package ratelimiter
