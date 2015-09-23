package ratelimiters

import (
	"net/http"
	"strconv"
	"testing"
	"time"
)

type TestClock struct {
	sleepDuration int64
}

func (clock *TestClock) Sleep(d time.Duration) {
	clock.sleepDuration = int64(d)
}
func TestRateLimitDivideZero(t *testing.T) {
	headers := http.Header{}
	headers.Add("X-Ratelimit-Available", "0")
	headers.Add("X-Ratelimit-Expiry", strconv.FormatInt(time.Now().Add(time.Second*10).Unix()*1000, 10))
	parser := TradierRateLimiter{Clock: &TestClock{}}

	if err := parser.ObeyRateLimit(headers); err != nil {
		t.Fatalf("should not have failed, %s", err.Error())
	}
}
