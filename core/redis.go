package core

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"strconv"
	"time"
)

func ttlToRedisCmd(k string, val int64) []string {
	return []string{"EXPIREAT", k, fmt.Sprint(time.Now().Unix() + val)}
}

func stringToRedisCmd(k, val string) []string {
	return []string{"SET", k, val}
}

func hashToRedisCmd(k string, val map[string]string) []string {
	cmd := []string{"HSET", k}
	for k, v := range val {
		cmd = append(cmd, k, v)
	}
	return cmd
}

func setToRedisCmd(k string, val []string) []string {
	cmd := []string{"SADD", k}
	return append(cmd, val...)
}

func listToRedisCmd(k string, val []string) []string {
	cmd := []string{"RPUSH", k}
	return append(cmd, val...)
}

func zsetToRedisCmd(k string, val []string) []string {
	cmd := []string{"ZADD", k}
	var key string

	for i, v := range val {
		if i%2 == 0 {
			key = v
			continue
		}

		cmd = append(cmd, v, key)
	}
	return cmd
}

type Serializer func(cmd []string) string

// RESPSerializer will serialize cmd to RESP
func RESPSerializer(cmd []string) string {
	s := ""
	s += "*" + strconv.Itoa(len(cmd)) + "\r\n"
	for _, arg := range cmd {
		s += "$" + strconv.Itoa(len(arg)) + "\r\n"
		s += arg + "\r\n"
	}
	return s
}

func keysTypes(pipe redis.Pipeliner, keys []string) ([]string, error) {
	for _, key := range keys {
		pipe.Type(key)
	}
	res, err := pipe.Exec()
	if err != nil {
		return nil, err
	}
	ensure(len(keys) == len(res))
	types := make([]string, len(res))
	for i := range res {
		types[i] = res[i].(*redis.StatusCmd).Val()
	}
	return types, nil
}

func dumpKeys(pipe redis.Pipeliner, keys []string, ttl bool) ([][]string, error) {
	if len(keys) == 0 || pipe == nil {
		return nil, nil
	}

	types, err := keysTypes(pipe, keys)
	if err != nil {
		return nil, err
	}

	for i, t := range types {
		switch t {
		case "string":
			pipe.Get(keys[i])
		case "list":
			pipe.LRange(keys[i], 0, -1)
		case "set":
			pipe.SMembers(keys[i])
		case "hash":
			pipe.HGetAll(keys[i])
		case "zset":
			pipe.ZRangeByScore(keys[i], &redis.ZRangeBy{
				Min:    "-inf",
				Max:    "+inf",
				Offset: 0,
				Count:  0,
			})
		default:
			return nil, fmt.Errorf("Key %s is of unreconized type %s", keys[i], t)
		}
	}
	values, _ := pipe.Exec()

	ensure(values != nil)
	ensure(len(types) == len(keys))
	ensure(len(types) == len(values))

	var ttls []redis.Cmder
	if ttl {
		for _, key := range keys {
			pipe.TTL(key)
		}
		ttls, err = pipe.Exec()
		if err != nil {
			return nil, err
		}
	}

	ensure(ttls != nil)
	ensure(len(ttls) == len(keys))
	ensure(len(ttls) == len(values))

	commands := make([][]string, 0, len(values))
	for i, v := range values {
		switch types[i] {
		case "string":
			commands = append(commands, stringToRedisCmd(keys[i], v.(*redis.StringCmd).Val()))
		case "list":
			commands = append(commands, listToRedisCmd(keys[i], v.(*redis.StringSliceCmd).Val()))
		case "set":
			commands = append(commands, setToRedisCmd(keys[i], v.(*redis.StringSliceCmd).Val()))
		case "hash":
			commands = append(commands, hashToRedisCmd(keys[i], v.(*redis.StringStringMapCmd).Val()))
		case "zset":
			commands = append(commands, zsetToRedisCmd(keys[i], v.(*redis.StringSliceCmd).Val()))
		default:
			return nil, fmt.Errorf("Key %s is of unreconized type %s", keys[i], types[i])
		}
		if ttls != nil && ttls[i] != nil {
			ttl := ttls[i].(*redis.DurationCmd).Val().Seconds()
			if ttl > 0 {
				commands = append(commands, ttlToRedisCmd(keys[i], int64(ttl)))
			}
		}
	}
	return commands, nil
}
