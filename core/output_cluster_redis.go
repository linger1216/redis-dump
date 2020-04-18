package core

import (
	"github.com/go-redis/redis/v7"
)

type OutputClusterRedisConfig struct {
	Url      []string
	Password string
}

type OutputClusterRedis struct {
	client *redis.ClusterClient
}

func (x *OutputClusterRedis) save(commands [][]string) error {
	pipe := x.client.Pipeline()
	for i := range commands {
		// todo
		// may be pool
		args := make([]interface{}, 0, len(commands[i]))
		for j := range commands[i] {
			args = append(args, commands[i][j])
		}
		pipe.Do(args...)
	}
	_, err := pipe.Exec()
	return err
}

func (x *OutputClusterRedis) close() {
	_ = x.client.Close()
}

func NewOutputClusterRedis(c *OutputClusterRedisConfig) *OutputClusterRedis {
	ret := &OutputClusterRedis{}
	ret.client = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    c.Url,
		Password: c.Password,
	})
	_, err := ret.client.Ping().Result()
	if err != nil {
		panic(err)
	}
	return ret
}
