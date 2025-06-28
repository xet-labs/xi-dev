package lib

import (
	// "encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
)
// --- RedisInstc methods (custom client access)
func (r RedisInstc) Get(key string) (string, error) {
	if r.Cli == nil {
		return "", redis.ErrClosed
	}
	return r.Cli.Get(RedisCtx, RedisKey(key)).Result()
}

func (r RedisInstc) Set(key, value string, ttl ...time.Duration) error {
	if r.Cli == nil {
		return redis.ErrClosed
	}
	exp := time.Duration(0)
	if len(ttl) > 0 {
		exp = ttl[0]
	}
	return r.Cli.Set(RedisCtx, RedisKey(key), value, exp).Err()
}

func (r RedisInstc) Del(key string) error {
	if r.Cli == nil {
		return redis.ErrClosed
	}
	return r.Cli.Del(RedisCtx, RedisKey(key)).Err()
}

func (r RedisInstc) Exists(key string) (bool, error) {
	if r.Cli == nil {
		return false, redis.ErrClosed
	}
	n, err := r.Cli.Exists(RedisCtx, RedisKey(key)).Result()
	return n > 0, err
}

func (r RedisInstc) Keys(pattern string) ([]string, error) {
	if r.Cli == nil {
		return nil, redis.ErrClosed
	}
	return r.Cli.Keys(RedisCtx, RedisKey(pattern)).Result()
}
