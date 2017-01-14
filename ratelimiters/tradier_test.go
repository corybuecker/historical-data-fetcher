package ratelimiters

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestClock struct {
	sleepDuration time.Duration
}

func (clock *TestClock) sleep(d time.Duration) {
	clock.sleepDuration = d
}

func (clock *TestClock) now() time.Time {
	return time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC)
}

var headers = map[string]string{}

func init() {
	headers["X-Ratelimit-Available"] = "60"
	headers["X-Ratelimit-Expiry"] = strconv.FormatInt(time.Now().Add(time.Second).Unix()*1000, 10)
}

func TestCreatedClock(t *testing.T) {
	rateLimiter := new(TradierRateLimiter)
	rateLimiter.ObeyRateLimit(headers)

	assert.NotNil(t, rateLimiter.Clock)
}

func TestCorrectSleepWithLimitLeft(t *testing.T) {
	headers["X-Ratelimit-Available"] = "60"
	rateLimiter := TradierRateLimiter{Clock: &TestClock{}}
	headers["X-Ratelimit-Expiry"] = strconv.FormatInt(rateLimiter.Clock.now().Add(time.Second*60).Unix()*1000, 10)
	rateLimiter.ObeyRateLimit(headers)

	assert.Equal(t, time.Duration(0), rateLimiter.Clock.(*TestClock).sleepDuration)
}

func TestCorrectSleep(t *testing.T) {
	headers["X-Ratelimit-Available"] = "1"

	rateLimiter := TradierRateLimiter{Clock: &TestClock{}}
	headers["X-Ratelimit-Expiry"] = strconv.FormatInt(rateLimiter.Clock.now().Add(time.Second*60).Unix()*1000, 10)
	rateLimiter.ObeyRateLimit(headers)

	assert.Equal(t, time.Duration(65*time.Second), rateLimiter.Clock.(*TestClock).sleepDuration)
}

func TestExpiryHeaderParseIntError(t *testing.T) {
	headers["X-Ratelimit-Expiry"] = "integer"
	headers["X-Ratelimit-Available"] = "60"
	rateLimiter := TradierRateLimiter{Clock: &TestClock{}}

	assert.EqualError(t, rateLimiter.ObeyRateLimit(headers), "strconv.ParseInt: parsing \"integer\": invalid syntax")
}

func TestAvailableHeaderParseIntError(t *testing.T) {
	headers["X-Ratelimit-Expiry"] = "60"
	headers["X-Ratelimit-Available"] = "integer"
	rateLimiter := TradierRateLimiter{Clock: &TestClock{}}

	assert.EqualError(t, rateLimiter.ObeyRateLimit(headers), "strconv.ParseInt: parsing \"integer\": invalid syntax")
}
