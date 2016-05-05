// Package utils contains common shared code.
package utils

import (
	"math"
	"time"
)

// Backoff holds the number of attempts as well as the min and max backoff delays.
type Backoff struct {
	attempt, Factor float64
	Min, Max        time.Duration
}

// Duration calculates the backoff delay and increments the attempts count.
func (b *Backoff) Duration(attempt float64) time.Duration {
	d := b.SetDuration(b.attempt)
	b.attempt++
	return d
}

// SetDuration calculates the backoff delay and caps it at the maximum delay.
func (b *Backoff) SetDuration(attempt float64) time.Duration {
	if b.Min == 0 {
		b.Min = 100 * time.Millisecond
	}

	if b.Max == 0 {
		b.Max = 10 * time.Second
	}

	// Calculate the wait duration.
	duration := float64(b.Min) * math.Pow(b.Factor, attempt)

	// Cap it at the maximum value.
	if duration > float64(b.Max) {
		return b.Max
	}

	return time.Duration(duration)
}

// Reset clears the number of attempts once the API call has succeeded.
func (b *Backoff) Reset() {
	b.attempt = 0
}

// Attempt returns the number of times the API call has failed.
func (b *Backoff) Attempt() float64 {
	return b.attempt
}
