package core

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"strings"
)

type OutputSingleRedisConfig struct {
	Network  string
	Url      string
	Password string
	DBNumber int
}

type OutputSingleRedis struct {
	client *redis.Client
}

func (x *OutputSingleRedis) save(commands [][]string) error {
	i := 0
	size := len(commands)
	if size > 0 {
		// resp first command is select
		// so we handle or ignore it
		cmd := commands[0]
		if len(cmd) == 2 && strings.ToLower(cmd[0]) == "select" {
			i = 1
		}
	}

	pipe := x.client.Pipeline()
	for ; i < size; i++ {
		// todo
		// may be pool
		args := make([]interface{}, 0, len(commands[i]))
		for j := range commands[i] {
			if commands[i][j] == "" {
				fmt.Printf("ca")
			}
			args = append(args, commands[i][j])
		}
		pipe.Do(args...)
	}
	_, err := pipe.Exec()
	return err
}

func (x *OutputSingleRedis) close() {
	_ = x.client.Close()
}

func NewOutputSingleRedis(c *OutputSingleRedisConfig) *OutputSingleRedis {
	if len(c.Network) == 0 {
		c.Network = "tcp"
	}
	ret := &OutputSingleRedis{}
	ret.client = redis.NewClient(&redis.Options{
		Network:  c.Network,
		Addr:     c.Url,
		Password: c.Password,
		DB:       c.DBNumber,
	})
	_, err := ret.client.Ping().Result()
	if err != nil {
		panic(err)
	}
	return ret
}
