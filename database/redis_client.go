package database

import redis "gopkg.in/redis.v5"

type RedisClient struct {
	Client *redis.Client
}

func (redisClient *RedisClient) HGetAll(key string) (results map[string]string, err error) {
	command := redisClient.Client.HGetAll(key)
	results, err = command.Result()
	return
}

func (redisClient *RedisClient) HSet(key, field string, value interface{}) (err error) {
	command := redisClient.Client.HSet(key, field, value)
	err = command.Err()
	return
}

func (redisClient *RedisClient) HIncrBy(key, field string, value int64) (err error) {
	command := redisClient.Client.HIncrBy(key, field, value)
	err = command.Err()
	return
}
func (redisClient *RedisClient) SMembers(key string) (results []string, err error) {
	command := redisClient.Client.SMembers(key)
	results, err = command.Result()
	return
}
