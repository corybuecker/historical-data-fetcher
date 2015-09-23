package ratelimiters

import (
	"net/http"
	"strconv"
	"time"
)

type ClockInterface interface {
	Sleep(time.Duration)
}

type Clock struct{}

func (clock *Clock) Sleep(d time.Duration) { time.Sleep(d) }

type TradierRateLimiter struct {
	Clock ClockInterface
}

func (rateLimiter *TradierRateLimiter) ObeyRateLimit(headers http.Header) error {
	if rateLimiter.Clock == nil {
		rateLimiter.Clock = &Clock{}
	}
	var rateLimitAvailable int64
	var rateLimitExpires time.Duration
	var err error
	if headers.Get("X-Ratelimit-Available") != "" {
		var ratelimitAvailableTemp int64
		if ratelimitAvailableTemp, err = strconv.ParseInt(headers.Get("X-Ratelimit-Available"), 10, 8); err != nil {
			return err
		}
		rateLimitAvailable = ratelimitAvailableTemp
	}

	if ratelimitExpiresHeader, ok := headers["X-Ratelimit-Expiry"]; ok {
		var ratelimitExpiresTemp int64
		if ratelimitExpiresTemp, err = strconv.ParseInt(ratelimitExpiresHeader[0], 10, 64); err != nil {
			return err
		}
		rateLimitExpires = time.Unix(ratelimitExpiresTemp/1000, 0).Sub(time.Now())
	}

	if rateLimitAvailable > 0 {
		rateLimiter.Clock.Sleep(time.Duration(int64(rateLimitExpires)/rateLimitAvailable) + time.Millisecond*100)
	} else {
		rateLimiter.Clock.Sleep(time.Minute)
	}
	return nil
}
