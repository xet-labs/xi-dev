package db

import (
	"context"
	"strings"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
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
	lazyInit   func()

	mu   sync.RWMutex
	once sync.Once
}

// Global instance
var Rdb = &RedisLib{
	prefix:     "app",
	defaultCli: "redis",
	clients:    sharedClients,
	ctx:        context.Background(),
}

// RegisterLazyFn allows deferred initialization
func (r *RedisLib) RegisterLazyFn(fn func()) {
	r.lazyInit = fn
}

// Ensures lazyInit runs once
func (r *RedisLib) lazyFnOnce() {
	r.once.Do(func() {
		if r.lazyInit != nil {
			r.lazyInit()
		}
	})
}

// New returns a new RedisLib instance with optional prefix/context
func (r *RedisLib) New(defaultCli string, opts ...any) *RedisLib {
	Db.Init()
	r.lazyFnOnce()

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
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.clients[name]; exists {
		log.Warn().Msgf("Redis client '%s' already exists", name)
	}
	for n, c := range r.clients {
		if c == client {
			log.Warn().Msgf("Redis client already registered as '%s'", n)
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
	r.lazyFnOnce()
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := r.defaultCli
	if len(name) > 0 && strings.TrimSpace(name[0]) != "" {
		key = name[0]
	}
	return r.clients[key]
}

// SetDefault sets the default Redis client by name
func (r *RedisLib) SetDefault(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.clients) == 0 {
		r.defaultCli = name
		return
	}

	if cli, ok := r.clients[name]; ok {
		r.defaultCli = name
		r.client = cli
		log.Info().Msgf("Redis default: client set to '%s'", name)
	} else {
		log.Warn().Msgf("Redis default: client '%s' not found", name)
	}
}

// SetPrefix updates the Redis key prefix
func (r *RedisLib) SetPrefix(prefix string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.prefix = strings.TrimSpace(prefix)
}

// GetPrefix returns current Redis key prefix
func (r *RedisLib) GetPrefix() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.prefix
}

// SetCtx sets Redis context
func (r *RedisLib) SetCtx(ctx context.Context) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ctx = ctx
}

// GetCtx returns current context
func (r *RedisLib) GetCtx() context.Context {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.ctx
}

// GetDefault returns default client name
func (r *RedisLib) GetDefault() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.defaultCli
}

// With returns a new RedisLib bound to the given client name
func (r *RedisLib) With(cliName string) *RedisLib {
	return r.New(cliName, r.prefix)
}
