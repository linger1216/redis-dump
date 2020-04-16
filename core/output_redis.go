package core

import (
	"bufio"
	"bytes"
	"github.com/go-redis/redis/v7"
	"os"
	"strings"
)

type OutputRedis struct {
	client        *redis.Client
	clusterClient *redis.ClusterClient
}

func (x *OutputRedis) save(commands []string) error {
	var content bytes.Buffer
	for i := range commands {
		_, err := content.WriteString(commands[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (x *OutputRedis) close() {
	_ = x.w.Flush()
	_ = x.f.Close()
}

func NewOutputRedis(cfg *OutputFileConfig) *OutputRedis {
	ret := &OutputRedis{}
	flag := 0
	if strings.ToLower(cfg.Flag) == "append" {
		flag = os.O_RDWR | os.O_CREATE | os.O_APPEND
	} else if strings.ToLower(cfg.Flag) == "trunc" {
		flag = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	}

	writeSize := cfg.WriteSize
	if writeSize == 0 {
		writeSize = 4096
	}

	obj, err := os.OpenFile(cfg.Filename, flag, 0644)
	if err != nil {
		panic(err)
	}
	ret.f = obj
	ret.w = bufio.NewWriterSize(ret.f, writeSize)
	return ret
}
