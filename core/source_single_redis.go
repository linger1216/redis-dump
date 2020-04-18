package core

import "github.com/go-redis/redis/v7"

type sourceSingleRedisConfig struct {
	Network  string `json:"network"`
	Url      string `json:"url"`
	Password string `json:"password"`
	DBNumber int    `json:"dbNumber"`
	TTL      bool   `json:"ttl"`
	Match    string `json:"match"`
	Count    int64  `json:"count"`
}

func (s sourceSingleRedisConfig) newSource() source {
	return newSourceSingleRedis(s)
}

type sourceSingleRedis struct {
	c     *redis.Client
	TTL   bool
	Match string
	Count int64

	cursor uint64
	end    bool
}

func newSourceSingleRedis(conf sourceSingleRedisConfig) *sourceSingleRedis {
	ret := &sourceSingleRedis{}
	ret.TTL = conf.TTL
	ret.Match = conf.Match
	ret.Count = conf.Count

	ret.c = redis.NewClient(&redis.Options{
		Network:  conf.Network,
		Addr:     conf.Url,
		Password: conf.Password,
		DB:       conf.DBNumber,
	})
	return ret
}

func (s *sourceSingleRedis) has() bool {
	return !s.end
}

func (s *sourceSingleRedis) next() ([][]string, error) {
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

func (s *sourceSingleRedis) close() {
	_ = s.c.Close()
}
