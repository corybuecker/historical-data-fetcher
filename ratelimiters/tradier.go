package ratelimiters

import (
	"log"
	"strconv"
	"time"
)

type ClockInterface interface {
	Sleep(time.Duration)
	Now() time.Time
}

type Clock struct{}

func (clock *Clock) Sleep(d time.Duration) { time.Sleep(d) }
func (clock *Clock) Now() time.Time        { return time.Now() }

type TradierRateLimiter struct {
	Clock ClockInterface
}

func (rateLimiter *TradierRateLimiter) ObeyRateLimit(headers map[string]string) error {
	if rateLimiter.Clock == nil {
		rateLimiter.Clock = &Clock{}
	}

	var rateLimitAvailable int64
	var rateLimitExpires time.Duration
	var err error
	var ratelimitAvailableTemp int64

	if ratelimitAvailableTemp, err = strconv.ParseInt(headers["X-Ratelimit-Available"], 10, 8); err != nil {
		return err
	}
	rateLimitAvailable = ratelimitAvailableTemp

	var ratelimitExpiresTemp int64
	if ratelimitExpiresTemp, err = strconv.ParseInt(headers["X-Ratelimit-Expiry"], 10, 64); err != nil {
		return err
	}
	rateLimitExpires = time.Unix(ratelimitExpiresTemp/1000, 0).Sub(rateLimiter.Clock.Now())

	if rateLimitAvailable > 0 {
		sleepTime := time.Duration(int64(rateLimitExpires)/rateLimitAvailable) + time.Millisecond*100
		log.Printf("sleeping for %s to respect rate limit", sleepTime.String())
		rateLimiter.Clock.Sleep(sleepTime)
	} else {
		rateLimiter.Clock.Sleep(time.Minute)
	}
	return nil
}
