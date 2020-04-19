package core

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"strings"
	"time"
)

type outputSingleRedisConfig struct {
	Network  string `json:"network"`
	Url      string `json:"url"`
	Password string `json:"password"`
	DBNumber int    `json:"dbNumber"`
}

func (o *outputSingleRedisConfig) newOutput() output {
	return NewOutputSingleRedis(o)
}

func NewOutputSingleRedis(conf *outputSingleRedisConfig) *OutputSingleRedis {

	if conf == nil || len(conf.Url) == 0 {
		return nil
	}

	if len(conf.Network) == 0 {
		conf.Network = "tcp"
	}
	ret := &OutputSingleRedis{}
	ret.client = redis.NewClient(&redis.Options{
		Network:      conf.Network,
		Addr:         conf.Url,
		Password:     conf.Password,
		DB:           conf.DBNumber,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	})
	_, err := ret.client.Ping().Result()
	if err != nil {
		panic(err)
	}
	return ret
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
