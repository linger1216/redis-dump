package core

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"strconv"
	"strings"
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

// RedisCmdSerializer will serialize cmd to a string with redis commands
func RedisCmdSerializer(cmd []string) string {
	return strings.Join(cmd, " ")
}

func keysTypes(client *redis.Client, keys []string) ([]string, error) {
	pipe := client.Pipeline()
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

func DumpKeys(client *redis.Client, keys []string, ttl bool, s Serializer) ([]string, error) {
	if len(keys) == 0 {
		return nil, nil
	}

	if client == nil || s == nil {
		return nil, fmt.Errorf("invalid para")
	}

	types, err := keysTypes(client, keys)
	if err != nil {
		return nil, err
	}

	pipe := client.Pipeline()
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

	commands := make([]string, 0, len(values))
	for i, v := range values {
		switch types[i] {
		case "string":
			commands = append(commands, s(stringToRedisCmd(keys[i], v.(*redis.StringCmd).Val())))
		case "list":
			commands = append(commands, s(listToRedisCmd(keys[i], v.(*redis.StringSliceCmd).Val())))
		case "set":
			commands = append(commands, s(setToRedisCmd(keys[i], v.(*redis.StringSliceCmd).Val())))
		case "hash":
			commands = append(commands, s(hashToRedisCmd(keys[i], v.(*redis.StringStringMapCmd).Val())))
		case "zset":
			commands = append(commands, s(zsetToRedisCmd(keys[i], v.(*redis.StringSliceCmd).Val())))
		default:
			return nil, fmt.Errorf("Key %s is of unreconized type %s", keys[i], types[i])
		}

		if ttls != nil && ttls[i] != nil {
			ttl := ttls[i].(*redis.DurationCmd).Val().Seconds()
			if ttl > 0 {
				commands = append(commands, s(ttlToRedisCmd(keys[i], int64(ttl))))
			}
		}
	}
	return commands, nil
}
