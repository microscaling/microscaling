// Package utils contains common shared code.
package utils

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// Backoff holds the number of attempts as well as the min and max backoff delays.
type Backoff struct {
	sync.RWMutex
	attempt, Factor      int
	maxAttempts, waiting bool
	Min, Max             time.Duration
	Timer                *time.Timer
}

// Reset clears the number of attempts once the API call has succeeded.
func (b *Backoff) Reset() {
	b.Lock()
	defer b.Unlock()

	if b.attempt > 0 {
		log.Debugf("Backoff succeeded after %d attempts", b.attempt)
	}

	if b.waiting {
		log.Errorf("Backoff waiting flag is unexpectedly set")
	}

	b.attempt = 0
}

// Waiting flag is true while waiting for the backoff duration. Prevents
// any scaling actions.
func (b *Backoff) Waiting() bool {
	b.RLock()
	defer b.RUnlock()

	return b.waiting
}

// MaxAttempts returns true when the wait duration has reached the maximum value.
func (b *Backoff) MaxAttempts() bool {
	b.RLock()
	defer b.RUnlock()

	return b.maxAttempts
}

// SetTimer calculates the duration and sets an appropriate timer. When it pops it will send on the channel. Returns true if maxAttempts has been reached or exceeded.
func (b *Backoff) Backoff(c chan struct{}) error {
	b.Lock()
	defer b.Unlock()

	multiplier := math.Pow(float64(b.Factor), float64(b.attempt))
	duration := time.Duration(float64(b.Min) * multiplier)
	log.Debugf("Backing off for %s", duration)

	// Check whether we've reached the max backoff duration
	if duration > b.Max {
		return fmt.Errorf("Exceeded max backoff attempts")
	}

	b.waiting = true
	b.attempt++
	b.Timer = time.NewTimer(duration)
	go func() {
		<-b.Timer.C
		log.Debug("Backff expired")
		b.waiting = false
		c <- struct{}{}
	}()

	return nil
}

func (b *Backoff) Stop() {
	b.Lock()
	defer b.Unlock()

	if b.waiting {
		b.Timer.Stop()
		b.waiting = false
	}
	return
}
