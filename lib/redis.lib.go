package lib

import (
	"context"
	"github.com/redis/go-redis/v9"
	"fmt"
	"log"
"encoding/json"
	"os"
	"strings"
	"sync"
	"time"
)



// RedisLib represents the global redis utility
type RedisLib struct {
	Prefix   string
	DefRdb   string
	DefCli   *redis.Client
	Clis     map[string]*redis.Client
	Ctx      context.Context
	redisLock sync.RWMutex

	Cli      func(cli ...string) *redis.Client
	Get      func(key string) (string, error)
	Set      func(key string, value interface{}, ttl ...time.Duration) error
	Del      func(key string) error
	Exists   func(key string) (bool, error)
	Key      func(key string) string
	Keys     func(pattern string) ([]string, error)
	FlushAll func() error
	With     func(cli string) RedisInstc
}

var (
	RedisPrefix = os.Getenv("APP_ABBR")
	RedisDefRdb = "redis"
	RedisDefCli *redis.Client

	RedisClis = make(map[string]*redis.Client)
	RedisCtx  = context.Background()
	redisLock sync.RWMutex
)

// --- Initialization of RedisLib interface
var Redis = RedisLib{
	Cli:      redisCli,
	Get:      redisGet,
	Set:      redisSet,
	Del:      redisDel,
	Exists:   redisExists,
	Key:      redisKey,
	Keys:     redisKeys,
	FlushAll: redisFlushAll,
	With: func(cli string) RedisInstc {
		return RedisInstc{Cli: redisCli(cli)}
	},
}

type RedisInstc struct {
	Cli *redis.Client
}

// --- Internal redisCli function with optional argument
func redisCli(cli ...string) *redis.Client {
	redisLock.RLock()
	defer redisLock.RUnlock()

	name := RedisDefRdb
	if len(cli) > 0 && strings.TrimSpace(cli[0]) != "" {
		name = cli[0]
	}

	if RedisDefCli, ok := RedisClis[name]; ok {
		return RedisDefCli
	}

	available := make([]string, 0, len(RedisClis))
	for k := range RedisClis {
		available = append(available, k)
	}

	log.Printf("❌ Redis client '%s' not found", name)
	log.Printf("🧩 Available Redis clients: %s", strings.Join(available, ", "))
	return nil
}

// --- Key prefixing
func redisKey(key string) string {
	return RedisPrefix + ":" + key
}

// --- Default client CRUD operations
func redisGet(key string) (string, error) {
	if RedisDefCli == nil {
		return "", redis.ErrClosed
	}
	return RedisDefCli.Get(RedisCtx, redisKey(key)).Result()
}


// redisSet stores any Go value in Redis after JSON serialization
func redisSet(key string, value interface{}, ttl ...time.Duration) error {
	if RedisDefCli == nil {
		return redis.ErrClosed
	}

	// Marshal the value to JSON
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// Set TTL if provided
	exp := time.Duration(0)
	if len(ttl) > 0 {
		exp = ttl[0]
	}

	// Store the data in Redis
	return RedisDefCli.Set(RedisCtx, redisKey(key), data, exp).Err()
}

func redisDel(key string) error {
	if RedisDefCli == nil {
		return redis.ErrClosed
	}
	return RedisDefCli.Del(RedisCtx, redisKey(key)).Err()
}

func redisExists(key string) (bool, error) {
	if RedisDefCli == nil {
		return false, redis.ErrClosed
	}
	n, err := RedisDefCli.Exists(RedisCtx, redisKey(key)).Result()
	return n > 0, err
}

func redisKeys(pattern string) ([]string, error) {
	if RedisDefCli == nil {
		return nil, redis.ErrClosed
	}
	return RedisDefCli.Keys(RedisCtx, redisKey(pattern)).Result()
}

func redisFlushAll() error {
	if RedisDefCli == nil {
		return redis.ErrClosed
	}
	return RedisDefCli.FlushAll(RedisCtx).Err()
}

// --- RedisInstc methods (custom client access)
func (r RedisInstc) Get(key string) (string, error) {
	if r.Cli == nil {
		return "", redis.ErrClosed
	}
	return r.Cli.Get(RedisCtx, redisKey(key)).Result()
}

func (r RedisInstc) Set(key, value string, ttl ...time.Duration) error {
	if r.Cli == nil {
		return redis.ErrClosed
	}
	exp := time.Duration(0)
	if len(ttl) > 0 {
		exp = ttl[0]
	}
	return r.Cli.Set(RedisCtx, redisKey(key), value, exp).Err()
}

func (r RedisInstc) Del(key string) error {
	if r.Cli == nil {
		return redis.ErrClosed
	}
	return r.Cli.Del(RedisCtx, redisKey(key)).Err()
}

func (r RedisInstc) Exists(key string) (bool, error) {
	if r.Cli == nil {
		return false, redis.ErrClosed
	}
	n, err := r.Cli.Exists(RedisCtx, redisKey(key)).Result()
	return n > 0, err
}

func (r RedisInstc) Keys(pattern string) ([]string, error) {
	if r.Cli == nil {
		return nil, redis.ErrClosed
	}
	return r.Cli.Keys(RedisCtx, redisKey(pattern)).Result()
}
