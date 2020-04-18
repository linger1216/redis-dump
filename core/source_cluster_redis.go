package core

import "github.com/go-redis/redis/v7"

type sourceClusterRedisConfig struct {
	Url      []string `json:"url"`
	Password string   `json:"password"`
	TTL      bool     `json:"ttl"`
	Match    string   `json:"match"`
	Count    int64    `json:"count"`
}

func (s sourceClusterRedisConfig) newSource() source {
	return newSourceClusterRedis(s)
}

func newSourceClusterRedis(conf sourceClusterRedisConfig) *sourceClusterRedis {
	ret := &sourceClusterRedis{}
	ret.TTL = conf.TTL
	ret.Match = conf.Match
	ret.Count = conf.Count
	ret.c = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    conf.Url,
		Password: conf.Password,
	})
	return ret
}

type sourceClusterRedis struct {
	c     *redis.ClusterClient
	TTL   bool
	Match string
	Count int64

	cursor uint64
	end    bool
}

func (s *sourceClusterRedis) has() bool {
	return !s.end
}

func (s *sourceClusterRedis) next() ([][]string, error) {
	keys, cursor := s.c.Scan(s.cursor, s.Match, s.Count).Val()
	pipe := s.c.Pipeline()
	commands, err := dumpKeys(pipe, keys, s.TTL)
	if err != nil {
		return nil, err
	}
	s.cursor = cursor
	if cursor == 0 {
		s.end = true
	}
	return commands, err
}

func (s *sourceClusterRedis) close() {
	_ = s.c.Close()
}
