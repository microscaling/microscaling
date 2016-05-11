// Package utils contains common shared code.
package utils

import (
	"math"
	"time"
)

// Backoff holds the number of attempts as well as the min and max backoff delays.
type Backoff struct {
	attempt, Factor      int
	maxAttempts, waiting bool
	Min, Max             time.Duration
}

// Duration sets the waiting flag, calculates the backoff delay and increments
// the attempts count.
func (b *Backoff) Duration(attempt int) time.Duration {
	b.waiting = true

	d := b.CalcDuration(b.attempt)
	b.attempt++

	return d
}

// CalcDuration calculates the backoff delay and caps it at the maximum delay.
func (b *Backoff) CalcDuration(attempt int) time.Duration {
	// Default to sensible values when not configured.
	if b.Min == 0 {
		b.Min = 100 * time.Millisecond
	}
	if b.Max == 0 {
		b.Max = 10 * time.Second
	}
	if b.Factor == 0 {
		b.Factor = 2
	}

	// Calculate the wait duration.
	duration := float64(b.Min) * math.Pow(float64(b.Factor), float64(attempt))

	// Cap it at the maximum value.
	if duration > float64(b.Max) {
		b.maxAttempts = true
		return b.Max
	}

	return time.Duration(duration)
}

// Reset clears the number of attempts once the API call has succeeded.
func (b *Backoff) Reset() {
	b.attempt = 0
}

// Waiting flag is true while waiting for the backoff duration. Prevents
// any scaling actions.
func (b *Backoff) Waiting() bool {
	return b.waiting
}

// Clear the waiting flag after the backoff duration.
func (b *Backoff) Clear() {
	b.waiting = false
}

// Attempt returns the number of times the API call has failed.
func (b *Backoff) Attempt() int {
	return b.attempt
}

// MaxAttempts returns true when the wait duration has reached the maximum value.
func (b *Backoff) MaxAttempts() bool {
	return b.maxAttempts
}
