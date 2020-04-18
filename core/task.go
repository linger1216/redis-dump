package core

import (
	"github.com/go-redis/redis/v7"
)

type Task interface {
	exec() error
	close()
}

type StandaloneTask struct {
	cli      *redis.Client
	dbNumber int
	ttl      bool
	match    string
	batch    int
	out      output
}

func (c *StandaloneTask) exec() error {
	panic("implement me")
}

func (c *StandaloneTask) close() {
	panic("implement me")
}

type ClusterTask struct {
	cli   *redis.ClusterClient
	ttl   bool
	match string
	batch int
	out   []output
}

func (c *ClusterTask) exec() error {
	panic("implement me")
}

func (c *ClusterTask) close() {
	panic("implement me")
}

func newTasks(c *DumpConfig) ([]Task, error) {
	//if c == nil {
	//	return nil, nil
	//}
	//
	//task := make([]Task, 0)
	//if len(c.Src.Url) > 1 {
	//	client := redis.NewClusterClient(&redis.ClusterOptions{
	//		Addrs:        c.Src.Url,
	//		Password:     c.Src.Password,
	//		ReadTimeout:  time.Second * 5,
	//		WriteTimeout: time.Second * 5,
	//	})
	//	if _, err := client.Ping().Result(); err != nil {
	//		return nil, err
	//	}
	//
	//	outs := make([]output, 0)
	//	if len(c.Dest.File.FileName) > 0 {
	//		outs = append(outs, NewOutputCSV(&outputFileConfig{
	//			Flag:      "trunk",
	//			Filename:  c.Dest.File.FileName,
	//			WriteSize: 4096,
	//		}))
	//	}
	//
	//	if len(c.Dest.Redis.Url) > 0 {
	//
	//	}
	//
	//	task = append(task, &ClusterTask{
	//		cli:   client,
	//		TTL:   c.Src.TTL,
	//		Match: c.Src.Match,
	//		Batch: c.Common.Batch,
	//		out:   nil,
	//	})
	//
	//} else {
	//
	//}
	//
	//for _, Url := range c.Src.Url {
	//
	//}

	return nil, nil
}
