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

// Shared global Redis clients
var sharedClients = make(map[string]*redis.Client)

// RedisLib wraps Redis client management and access
type RedisLib struct {
	prefix     string
	defaultCli string
	clients    map[string]*redis.Client
	client     *redis.Client
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

// RegisterLazyInit allows deferred initialization
func (r *RedisLib) RegisterLazyInit(fn func()) {
	r.lazyInit = fn
}

// New returns a new RedisLib instance with optional prefix/context
func (r *RedisLib) New(defaultCli string, opts ...any) *RedisLib {
	r.lazyInitOnce()

	prefix, ctx := r.prefix, r.ctx
	for _, opt := range opts {
		switch v := opt.(type) {
		case string:
			if s := strings.TrimSpace(v); s != "" {
				prefix = s
			}
		case context.Context:
			ctx = v
		}
	}

	return &RedisLib{
		prefix:     prefix,
		defaultCli: defaultCli,
		client:     r.GetCli(defaultCli),
		clients:    sharedClients,
		ctx:        ctx,
	}
}

// SetCli registers a new Redis client
func (r *RedisLib) SetCli(name string, client *redis.Client) {
	r.rw.Lock()
	defer r.rw.Unlock()

	if _, exists := r.clients[name]; exists {
		log.Printf("⚠️  Redis client '%s' already exists", name)
	}
	for n, c := range r.clients {
		if c == client {
			log.Printf("⚠️  Redis client already registered as '%s'", n)
			break
		}
	}

	r.clients[name] = client
	if r.client == nil || r.defaultCli == name || strings.TrimSpace(r.defaultCli) == "" {
		r.client = client
		r.defaultCli = name
	}
}

// GetCli returns a Redis client by name or default
func (r *RedisLib) GetCli(name ...string) *redis.Client {
	r.lazyInitOnce()
	r.rw.RLock()
	defer r.rw.RUnlock()

	key := r.defaultCli
	if len(name) > 0 && strings.TrimSpace(name[0]) != "" {
		key = name[0]
	}
	return r.clients[key]
}

// SetDefault sets the default Redis client by name
func (r *RedisLib) SetDefault(name string) {
	r.rw.Lock()
	defer r.rw.Unlock()

	if len(r.clients) == 0 {
		r.defaultCli = name
		return
	}

	if cli, ok := r.clients[name]; ok {
		r.defaultCli = name
		r.client = cli
		log.Printf("✅ Redis default: client set to '%s'", name)
	} else {
		log.Printf("⚠️  Redis default: client '%s' not found", name)
	}
}

// SetPrefix updates the Redis key prefix
func (r *RedisLib) SetPrefix(prefix string) {
	r.rw.Lock()
	defer r.rw.Unlock()
	r.prefix = strings.TrimSpace(prefix)
}

// GetPrefix returns current Redis key prefix
func (r *RedisLib) GetPrefix() string {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return r.prefix
}

// SetCtx sets Redis context
func (r *RedisLib) SetCtx(ctx context.Context) {
	r.rw.Lock()
	defer r.rw.Unlock()
	r.ctx = ctx
}

// GetCtx returns current context
func (r *RedisLib) GetCtx() context.Context {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return r.ctx
}

// GetDefault returns default client name
func (r *RedisLib) GetDefault() string {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return r.defaultCli
}

// With returns a new RedisLib bound to the given client name
func (r *RedisLib) With(cliName string) *RedisLib {
	return r.New(cliName, r.prefix)
}

// Internal: Ensures lazyInit runs once
func (r *RedisLib) lazyInitOnce() {
	r.once.Do(func() {
		if r.lazyInit != nil {
			r.lazyInit()
		}
	})
}

// ---------------- Redis command wrappers ----------------

func (r *RedisLib) key(k string) string {
	return r.prefix + ":" + k
}

func (r *RedisLib) Set(k string, v any, ttl time.Duration) error {
	if r.client == nil {
		return redis.ErrClosed
	}
	return r.client.Set(r.ctx, r.key(k), v, ttl).Err()
}

func (r *RedisLib) Get(k string) (string, error) {
	if r.client == nil {
		return "", redis.ErrClosed
	}
	return r.client.Get(r.ctx, r.key(k)).Result()
}

func (r *RedisLib) GetBytes(k string) ([]byte, error) {
	if r.client == nil {
		return nil, redis.ErrClosed
	}
	return r.client.Get(r.ctx, r.key(k)).Bytes()
}

func (r *RedisLib) SetJson(k string, v any, ttl time.Duration) error {
	val, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}
	return r.Set(k, val, ttl)
}

func (r *RedisLib) GetJson(k string, out any) error {
	v, err := r.GetBytes(k)
	if err != nil {
		return err
	}
	return json.Unmarshal(v, out)
}

func (r *RedisLib) Del(keys ...string) error {
	if r.client == nil {
		return redis.ErrClosed
	}
	if len(keys) == 0 {
		return nil
	}

	var redisKeys []string
	for _, k := range keys {
		redisKeys = append(redisKeys, r.key(k))
	}

	return r.client.Del(r.ctx, redisKeys...).Err()
}

func (r *RedisLib) Exists(k string) (bool, error) {
	if r.client == nil {
		return false, redis.ErrClosed
	}
	n, err := r.client.Exists(r.ctx, r.key(k)).Result()
	return n > 0, err
}

func (r *RedisLib) Keys(pattern string) ([]string, error) {
	if r.client == nil {
		return nil, redis.ErrClosed
	}
	return r.client.Keys(r.ctx, r.key(pattern)).Result()
}

func (r *RedisLib) FlushAll() error {
	if r.client == nil {
		return redis.ErrClosed
	}
	return r.client.FlushAll(r.ctx).Err()
}
