package core

import "github.com/go-redis/redis/v7"

type sourceSingleRedis struct {
	c     *redis.Client
	TTL   bool
	Match string
	Count int64

	cursor uint64
	end    bool
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
