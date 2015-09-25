package ratelimiters

import (
	"net/http"
	"strconv"
	"testing"
	"time"
)

type TestClock struct {
	sleepDuration time.Duration
}

func (clock *TestClock) Sleep(d time.Duration) {
	clock.sleepDuration = d
}
func (clock *TestClock) Now() time.Time {
	return time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC)
}

var headers http.Header

func init() {
	headers = http.Header{}
	headers.Set("X-Ratelimit-Available", "60")
	headers.Set("X-Ratelimit-Expiry", strconv.FormatInt(time.Now().Add(time.Second).Unix()*1000, 10))
}

func TestCreatedClock(t *testing.T) {
	rateLimiter := new(TradierRateLimiter)
	rateLimiter.ObeyRateLimit(headers)
	if rateLimiter.Clock == nil {
		t.Fatalf("should have created a clock")
	}
}

func TestCorrectSleep(t *testing.T) {
	headers.Set("X-Ratelimit-Available", "60")

	rateLimiter := TradierRateLimiter{Clock: &TestClock{}}
	headers.Set("X-Ratelimit-Expiry", strconv.FormatInt(rateLimiter.Clock.Now().Add(time.Second*60).Unix()*1000, 10))

	rateLimiter.ObeyRateLimit(headers)
	if rateLimiter.Clock.(*TestClock).sleepDuration != time.Second+time.Millisecond*100 {
		t.Fatalf("expected around one second, got %s", rateLimiter.Clock.(*TestClock).sleepDuration)
	}
}

func TestExpiryHeaderParseIntError(t *testing.T) {
	headers.Set("X-Ratelimit-Expiry", "integer")
	headers.Set("X-Ratelimit-Available", "60")

	rateLimiter := TradierRateLimiter{Clock: &TestClock{}}
	if err := rateLimiter.ObeyRateLimit(headers); err == nil {
		t.Fatalf("should have received error")
	}
}

func TestAvailableHeaderParseIntError(t *testing.T) {
	headers.Set("X-Ratelimit-Expiry", "60")
	headers.Set("X-Ratelimit-Available", "integer")
	rateLimiter := TradierRateLimiter{Clock: &TestClock{}}
	if err := rateLimiter.ObeyRateLimit(headers); err == nil {
		t.Fatalf("should have received error")
	}
}

func TestRateLimitDivideZero(t *testing.T) {
	headers.Set("X-Ratelimit-Available", "0")
	headers.Set("X-Ratelimit-Expiry", strconv.FormatInt(time.Now().Add(time.Second*10).Unix()*1000, 10))
	rateLimiter := TradierRateLimiter{Clock: &TestClock{}}
	rateLimiter.ObeyRateLimit(headers)
	if rateLimiter.Clock.(*TestClock).sleepDuration != time.Minute {
		t.Fatalf("should have slept for 60 seconds, got %d", rateLimiter.Clock.(*TestClock).sleepDuration)
	}
}
