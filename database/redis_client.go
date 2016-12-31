package database

import redis "gopkg.in/redis.v5"

type DatabaseClient interface {
	HGetAll(string) (map[string]string, error)
	HSet(string, string, string) error
	HIncrBy(string, string, int64) error
	SMembers(string) ([]string, error)
}

type RedisClient struct {
	Client *redis.Client
}

func (redisClient *RedisClient) HGetAll(key string) (results map[string]string, err error) {
	command := redisClient.Client.HGetAll(key)
	results, err = command.Result()
	return results, err
}

func (redisClient *RedisClient) HSet(key, field, value string) (err error) {
	command := redisClient.Client.HSet(key, field, value)
	err = command.Err()
	return err
}

func (redisClient *RedisClient) HIncrBy(key, field string, value int64) (err error) {
	command := redisClient.Client.HIncrBy(key, field, value)
	err = command.Err()
	return err
}
func (redisClient *RedisClient) SMembers(key string) (results []string, err error) {
	command := redisClient.Client.SMembers(key)
	results, err = command.Result()
	return results, err
}
