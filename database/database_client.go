package database

type DatabaseClient interface {
	HGetAll(string) (map[string]string, error)
	HSet(string, string, interface{}) error
	HIncrBy(string, string, int64) error
	SMembers(string) ([]string, error)
}
