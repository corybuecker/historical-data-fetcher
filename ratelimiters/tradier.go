package ratelimiters

import (
	"log"
	"strconv"
	"time"
)

type ClockInterface interface {
	sleep(time.Duration)
	now() time.Time
}

type Clock struct{}

func (clock *Clock) sleep(d time.Duration) {
	log.Printf("sleeping for %s to respect rate limit", d.String())
	time.Sleep(d)
}
func (clock *Clock) now() time.Time { return time.Now() }

type TradierRateLimiter struct {
	Clock ClockInterface
}

func (rateLimiter *TradierRateLimiter) ObeyRateLimit(headers map[string]string) error {
	if rateLimiter.Clock == nil {
		rateLimiter.Clock = &Clock{}
	}

	var rateLimitExpires time.Duration
	var err error
	var ratelimitAvailableTemp, ratelimitExpiresTemp int64

	if ratelimitAvailableTemp, err = strconv.ParseInt(headers["X-Ratelimit-Available"], 10, 8); err != nil {
		return err
	}

	if ratelimitExpiresTemp, err = strconv.ParseInt(headers["X-Ratelimit-Expiry"], 10, 64); err != nil {
		return err
	}

	rateLimitExpires = time.Unix(ratelimitExpiresTemp/1000, 0).Sub(rateLimiter.Clock.now())

	rateLimiter.Clock.sleep(calculateSleepTime(ratelimitAvailableTemp, int64(rateLimitExpires)))

	return nil
}

func calculateSleepTime(available int64, expires int64) time.Duration {
	if available > 10 {
		return time.Duration(0)
	}

	return time.Duration(expires) + time.Second*5
}
