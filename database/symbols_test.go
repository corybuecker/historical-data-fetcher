package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	redis "gopkg.in/redis.v5"
)

var symbols *Symbols

func init() {
	redisConnection = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "master",
		SentinelAddrs: []string{":26379"},
		DB:            1,
	})

	redisConnection.FlushDb()

	now = time.Now().Truncate(time.Second)

	symbols = &Symbols{client: &RedisClient{
		Client: redisConnection,
	}}
}

func TestInitializeWithoutKeys(t *testing.T) {
	assert.Nil(t, symbols.Initialize())
	assert.Empty(t, symbols.Symbols)
}

func TestInitializeWithKeys(t *testing.T) {
	redisConnection.SAdd("NASDAQ", "AAPL")
	assert.Nil(t, symbols.Initialize())
	assert.NotEmpty(t, symbols.Symbols)
}

func TestFilterNil(t *testing.T) {
	symbols.Initialize()
	symbols.Filter(time.Now())
	assert.Empty(t, symbols.Symbols)
}

func TestFilter(t *testing.T) {
	redisConnection.HSet("NASDAQ:AAPL", "last_date_fetched", time.Now().Format(time.RFC3339))
	symbols.Initialize()
	symbols.Filter(time.Now().Add(time.Minute * -1))
	assert.Empty(t, symbols.Symbols)
}

func TestFilterMatch(t *testing.T) {
	redisConnection.HSet("NASDAQ:AAPL", "last_date_fetched", time.Now().Add(time.Minute*-1).Format(time.RFC3339))
	symbols.Initialize()
	symbols.Filter(time.Now())
	assert.NotEmpty(t, symbols.Symbols)
}
