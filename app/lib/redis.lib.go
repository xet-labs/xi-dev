package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisPrefix = os.Getenv("APP_ABBR")
	RedisDefRdb = "redis"
	RedisDefCli *redis.Client

	RedisClis = make(map[string]*redis.Client)
	RedisCtx  = context.Background()
	redisLock sync.RWMutex
)

// RedisLib represents the global redis utility
type RedisLib struct {
	Prefix string
	DefRdb string
	DefCli *redis.Client
	Clis   map[string]*redis.Client
	Ctx    context.Context
	Lock   sync.RWMutex

	Cli func(cli ...string) *redis.Client

	GetBytes  func(key string) ([]byte, error)
	GetString func(key string) (string, error)
	GetJson   func(key string, target interface{}) error
	SetBytes  func(key string, value []byte, ttl time.Duration) error
	SetString func(key string, value string, ttl time.Duration) error
	SetJson   func(key string, value interface{}, ttl time.Duration) error
	Del       func(key string) error
	Exists    func(key string) (bool, error)
	Key       func(key string) string
	Keys      func(pattern string) ([]string, error)
	FlushAll  func() error
	With      func(cli string) RedisInstc
}

// --- Initialization of RedisLib interface
var Redis = RedisLib{
	Cli:       redisCli,
	GetBytes:  RedisGetBytes,
	GetString: RedisGetString,
	GetJson:   RedisGetJson,
	SetBytes:  RedisSetBytes,
	SetString: RedisSetString,
	SetJson:   RedisSetJson,
	Del:       RedisDel,
	Exists:    RedisExists,
	Key:       RedisKey,
	Keys:      RedisKeys,
	FlushAll:  RedisFlushAll,
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
func RedisKey(key string) string {
	return RedisPrefix + ":" + key
}

// Get []byte value
func RedisGetBytes(key string) ([]byte, error) {
	if RedisDefCli == nil {
		return nil, redis.ErrClosed
	}
	return RedisDefCli.Get(RedisCtx, RedisKey(key)).Bytes()
}

// Get string value
func RedisGetString(key string) (string, error) {
	if RedisDefCli == nil {
		return "", redis.ErrClosed
	}
	return RedisDefCli.Get(RedisCtx, RedisKey(key)).Result()
}

// Get and unmarshal JSON into target (pass pointer)
func RedisGetJson(key string, target interface{}) error {
	if RedisDefCli == nil {
		return redis.ErrClosed
	}

	data, err := RedisDefCli.Get(RedisCtx, RedisKey(key)).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, target)
}

func RedisSetBytes(key string, value []byte, ttl time.Duration) error {
	if RedisDefCli == nil {
		return redis.ErrClosed
	}
	return RedisDefCli.Set(RedisCtx, RedisKey(key), value, ttl).Err()
}

func RedisSetString(key string, value string, ttl time.Duration) error {
	if RedisDefCli == nil {
		return redis.ErrClosed
	}
	return RedisDefCli.Set(RedisCtx, RedisKey(key), value, ttl).Err()
}

func RedisSetJson(key string, value interface{}, ttl time.Duration) error {
	if RedisDefCli == nil {
		return redis.ErrClosed
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}

	return RedisDefCli.Set(RedisCtx, RedisKey(key), data, ttl).Err()
}

func RedisDel(key string) error {
	if RedisDefCli == nil {
		return redis.ErrClosed
	}
	return RedisDefCli.Del(RedisCtx, RedisKey(key)).Err()
}

func RedisExists(key string) (bool, error) {
	if RedisDefCli == nil {
		return false, redis.ErrClosed
	}
	n, err := RedisDefCli.Exists(RedisCtx, RedisKey(key)).Result()
	return n > 0, err
}

func RedisKeys(pattern string) ([]string, error) {
	if RedisDefCli == nil {
		return nil, redis.ErrClosed
	}
	return RedisDefCli.Keys(RedisCtx, RedisKey(pattern)).Result()
}

func RedisFlushAll() error {
	if RedisDefCli == nil {
		return redis.ErrClosed
	}
	return RedisDefCli.FlushAll(RedisCtx).Err()
}
