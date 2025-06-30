package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Shared global Redis clients and lock
var (
	sharedClients = make(map[string]*redis.Client)
	// sharedLock    = &sync.RWMutex{}
)

// Central utility
type RedisLib struct {
	prefix     string
	defaultCli string
	cli        *redis.Client
	clients    map[string]*redis.Client
	ctx        context.Context
	rw         sync.RWMutex
	once       sync.Once
	lazyInit   func()
}

// Global instance
var Redis = &RedisLib{
	prefix:     "app",
	defaultCli: "redis",
	clients:    sharedClients,
	ctx:        context.Background(),
}

// RegisterLazyInit sets a callback for deferred initialization.
func (r *RedisLib) RegisterLazyInit(fn func()) {
	r.lazyInit = fn
}

// New creates a new RedisLib instance
func (r *RedisLib) New(defaultCli string, opts ...interface{}) *RedisLib {
	r.once.Do(func() {
		if r.lazyInit != nil {
			r.lazyInit()
		}
	})
	cli := r.GetCli(defaultCli)

	prefix := r.prefix
	ctx := r.ctx

	for _, opt := range opts {
		switch v := opt.(type) {
		case string:
			if strings.TrimSpace(v) != "" {
				prefix = v
			}
		case context.Context:
			ctx = v
		}
	}

	return &RedisLib{
		prefix:     prefix,
		defaultCli: defaultCli,
		cli:        cli,
		clients:    sharedClients,
		ctx:        ctx,
		// rw:         sharedLock,
	}
}

// SetCli registers a new redis client
func (r *RedisLib) SetCli(name string, cli *redis.Client) {
	r.rw.Lock()
	defer r.rw.Unlock()

	// Check if client name already exists
	if _, exists := r.clients[name]; exists {
		log.Printf("⚠️  Redis client name '%s' already exists", name)
	}

	// Check if cli pointer is already stored under another name
	for n, c := range r.clients {
		if c == cli {
			log.Printf("⚠️  Redis client instance already registered as '%s'", n)
			break
		}
	}

	r.clients[name] = cli

	// Set as default if unset, or if default name matches, or default name is empty
	if Redis.cli == nil || Redis.defaultCli == name || strings.TrimSpace(Redis.defaultCli) == "" {
		Redis.cli = cli
		Redis.defaultCli = name
	}
}

// GetCli returns client by name or default
func (r *RedisLib) GetCli(name ...string) *redis.Client {
	r.once.Do(func() {
		if r.lazyInit != nil {
			r.lazyInit()
		}
	})
	r.rw.RLock()
	defer r.rw.RUnlock()

	key := r.defaultCli
	if len(name) > 0 && strings.TrimSpace(name[0]) != "" {
		key = name[0]
	}
	return r.clients[key]
}

// SetPrefix updates the Redis key prefix
func (r *RedisLib) SetPrefix(prefix string) {
	r.rw.Lock()
	defer r.rw.Unlock()
	r.prefix = strings.TrimSpace(prefix)
}

// GetPrefix returns the current Redis key prefix
func (r *RedisLib) GetPrefix() string {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return r.prefix
}

// SetDefault sets the default Redis client by name
func (r *RedisLib) SetDefault(name string) {
	r.rw.Lock()
	defer r.rw.Unlock()

	// If no clients have been added yet
	if len(r.clients) == 0 {
		r.defaultCli = name
		// log.Printf("⚠️  Redis client map empty, set default to '%s'", name)
		return
	}

	// If client with given name exists
	if cli, ok := r.clients[name]; ok {
		r.defaultCli = name
		r.cli = cli
		log.Printf("✅ Redis default: client set to '%s'", name)
	} else {
		log.Printf("⚠️  Redis default: client '%s' not found", name)
	}
}

// GetDefault returns the name of the default Redis client
func (r *RedisLib) GetDefault() string {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return r.defaultCli
}

// SetCtx updates the Redis context
func (r *RedisLib) SetCtx(ctx context.Context) {
	r.rw.Lock()
	defer r.rw.Unlock()
	r.ctx = ctx
}

// GetCtx returns the current Redis context
func (r *RedisLib) GetCtx() context.Context {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return r.ctx
}

// With returns new instance bound to given client name
func (r *RedisLib) With(cliName string) *RedisLib {
	return r.New(cliName, r.prefix)
}

// ---- Redis command wrappers ----

func (r *RedisLib) key(k string) string {
	return r.prefix + ":" + k
}

func (r *RedisLib) Get(k string) (string, error) {
	if r.cli == nil {
		return "", redis.ErrClosed
	}
	return r.cli.Get(r.ctx, r.key(k)).Result()
}

func (r *RedisLib) Set(k, val string, ttl time.Duration) error {
	if r.cli == nil {
		return redis.ErrClosed
	}
	return r.cli.Set(r.ctx, r.key(k), val, ttl).Err()
}

func (r *RedisLib) GetBytes(k string) ([]byte, error) {
	if r.cli == nil {
		return nil, redis.ErrClosed
	}
	return r.cli.Get(r.ctx, r.key(k)).Bytes()
}

func (r *RedisLib) SetBytes(k string, data []byte, ttl time.Duration) error {
	if r.cli == nil {
		return redis.ErrClosed
	}
	return r.cli.Set(r.ctx, r.key(k), data, ttl).Err()
}

func (r *RedisLib) GetJson(k string, out interface{}) error {
	data, err := r.GetBytes(k)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, out)
}

func (r *RedisLib) SetJson(k string, val interface{}, ttl time.Duration) error {
	data, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}
	return r.SetBytes(k, data, ttl)
}

func (r *RedisLib) Exists(k string) (bool, error) {
	if r.cli == nil {
		return false, redis.ErrClosed
	}
	n, err := r.cli.Exists(r.ctx, r.key(k)).Result()
	return n > 0, err
}

func (r *RedisLib) Del(k string) error {
	if r.cli == nil {
		return redis.ErrClosed
	}
	return r.cli.Del(r.ctx, r.key(k)).Err()
}

func (r *RedisLib) Keys(pattern string) ([]string, error) {
	if r.cli == nil {
		return nil, redis.ErrClosed
	}
	return r.cli.Keys(r.ctx, r.key(pattern)).Result()
}

func (r *RedisLib) FlushAll() error {
	if r.cli == nil {
		return redis.ErrClosed
	}
	return r.cli.FlushAll(r.ctx).Err()
}
