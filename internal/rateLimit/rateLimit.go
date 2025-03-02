package rateLimit

import (
	"errors"
	"iter"
	"time"
)

type RateLimit interface {
	// UpdateCount sets a new rate count per duration
	UpdateCount(count int)
	// UpdateDuration sets a new duration for rate
	UpdateDuration(duration time.Duration)
	// Iterate iterates synchronously
	Iterate() iter.Seq[bool]
	// IterateAsync iterates asynchronously, it can be called again after async iteration is stopped
	IterateAsync() (<-chan bool, error)
	// StopIterateAsync stops an async iteration to avoid a leaking goroutine
	StopIterateAsync()
}

type rateLimitImpl struct {
	count                int
	countDiff            int
	duration             time.Duration
	asyncChan            chan bool
	stopChan             chan struct{}
	isStarted            bool
	isReceivedStopSignal bool
}

func New(count int, duration time.Duration) RateLimit {
	return &rateLimitImpl{
		count:    count,
		duration: duration,
	}
}

func (rl *rateLimitImpl) UpdateCount(count int) {
	rl.countDiff = count - rl.count
	rl.count = count
}

func (rl *rateLimitImpl) UpdateDuration(duration time.Duration) {
	rl.duration = duration
}

func (rl *rateLimitImpl) Iterate() iter.Seq[bool] {
	return func(yield func(bool) bool) {
		start := time.Now()
		num := rl.count
		value := true
		for {
			if rl.countDiff > 0 {
				num += rl.countDiff
				rl.countDiff = 0
			}
			if time.Since(start) >= rl.duration {
				start = time.Now()
				num = rl.count
				value = true
			}
			if num > 0 {
				value = true
				num--
			} else {
				value = false
			}
			if !yield(value) {
				return
			}
		}
	}
}

func (rl *rateLimitImpl) IterateAsync() (<-chan bool, error) {
	if rl.isStarted {
		return nil, errors.New("iteration already running")
	}
	rl.asyncChan = make(chan bool, 1)
	rl.stopChan = make(chan struct{}, 1)
	go func() {
		defer func() { rl.isStarted = false }()
		defer close(rl.asyncChan)
		defer close(rl.stopChan)
		rl.isStarted = true
		rl.isReceivedStopSignal = false
		for value := range rl.Iterate() {
			select {
			case rl.asyncChan <- value:
			case <-rl.stopChan:
				break
			}
		}
	}()
	return rl.asyncChan, nil
}

func (rl *rateLimitImpl) StopIterateAsync() {
	if rl.isReceivedStopSignal {
		return
	}
	rl.isReceivedStopSignal = true
	defer close(rl.stopChan)
	select {
	case rl.stopChan <- struct{}{}:
	default:
	}
}
