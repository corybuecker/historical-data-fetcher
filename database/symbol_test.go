package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	redis "gopkg.in/redis.v5"
)

var redisConnection *redis.Client
var symbol *Symbol
var err error

func init() {
	redisConnection = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "master",
		SentinelAddrs: []string{":26379"},
		DB:            1,
	})

	redisConnection.FlushDb()

	symbol = &Symbol{
		Symbol:   "TESTSYMBOL",
		Exchange: "TESTEXCHANGE",
		client: &RedisClient{
			Client: redisConnection,
		},
	}

	now = time.Now().Truncate(time.Second)
}

func TestUpdateFetched(t *testing.T) {
	var value string
	symbol.UpdateFetched(now)
	value, err = redisConnection.HGet("TESTEXCHANGE:TESTSYMBOL", "last_date_fetched").Result()

	assert.Nil(t, err)
	assert.Equal(t, now.Format(time.RFC3339), value)
}

func TestSetLastUpdated(t *testing.T) {
	var value string
	var valueTime time.Time

	symbol.SetLastUpdated()

	value, err = redisConnection.HGet("TESTEXCHANGE:TESTSYMBOL", "last_date_fetched").Result()
	assert.Nil(t, err)

	valueTime, err = time.Parse(time.RFC3339, value)
	assert.Nil(t, err)

	assert.WithinDuration(t, now, valueTime, time.Duration(time.Second))
}

func TestMarkPresentInWiki(t *testing.T) {
	var value string

	symbol.MarkPresentInWiki()

	value, err = redisConnection.HGet("TESTEXCHANGE:TESTSYMBOL", "present_in_wiki").Result()
	assert.Nil(t, err)

	assert.Equal(t, "true", value)
}

func TestGetLastDate(t *testing.T) {
	assert.True(t, symbol.getLastDate())
	assert.Equal(t, now, symbol.LastDateFetched)
}

func TestGetLastDateWithoutKey(t *testing.T) {
	redisConnection.Del("TESTEXCHANGE:TESTSYMBOL")
	assert.False(t, symbol.getLastDate())
	assert.True(t, symbol.LastDateFetched.IsZero())
}
