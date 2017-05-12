package TokenBucket

import (
	"errors"
	"log"
	"time"
)

var (
	ErrNoTokenAvailable = errors.New("Running out of tokens")
	ErrNotInit          = errors.New("Uninitialized tokenbucket")
)

type RateLimiter interface {
	GetToken(timeout time.Duration) error
	Rate() uint32
	Stop() error
}

type TokenBucket struct {
	limit  uint32
	t      *time.Ticker
	bucket chan struct{}
	stop   chan bool
}

// NewTokenBucket creates a token bucket.
//     rate: tokens per second
func NewTokenBucket(rate uint32) RateLimiter {
	if rate == 0 {
		rate = 1
	}
	d := time.Second / time.Duration(rate)
	log.Printf("Token bucket ticker duration %v\n", d)
	tb := &TokenBucket{}
	tb.t = time.NewTicker(d)
	tb.limit = rate
	tb.bucket = make(chan struct{}, rate)
	tb.stop = make(chan bool, 1)
	go tb.run()
	return tb
}

func (tb *TokenBucket) run() {
	defer tb.t.Stop()
	for {
		select {
		case <-tb.t.C:
			select {
			case tb.bucket <- struct{}{}:
				// bucket not full
			default:
				// bucket is full
			}
		case <-tb.stop:
			close(tb.bucket)
			return
		}
	}
}

func (tb *TokenBucket) Stop() error {
	if tb == nil {
		return ErrNotInit
	}
	if tb.t == nil {
		return ErrNotInit
	}
	tb.stop <- true
	return nil
}

func (tb *TokenBucket) Rate() uint32 {
	if tb == nil {
		return 0
	}
	return tb.limit
}

// GetToken gets a token from the bucket.
// It returns nil if a token is obtained, non-nil if timed out.
// If timeout=0, it blocks until a token is available.
func (tb *TokenBucket) GetToken(timeout time.Duration) error {
	if tb == nil {
		return ErrNotInit
	}
	if timeout == 0 {
		<-tb.bucket
		return nil
	}
	select {
	case <-tb.bucket:
		return nil
	case <-time.After(timeout):
		return ErrNoTokenAvailable
	}
}
