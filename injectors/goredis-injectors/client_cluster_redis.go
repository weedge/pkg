package goredisinjectors

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/weedge/pkg/option"
)

type RedisClusterClientOptions struct {
	Addrs []string `mapstructure:"addrs"`
	// To route commands by latency or randomly, enable one of the following.
	Route string `mapstructure:"route"`
	RedisClientCommonOptions

	// optional, more details see redis.ClusterOptions
	// if use redis.ClusterOptions, config options can't use
	redisClusterOpts *redis.ClusterOptions
}

func (m *RedisClusterClientOptions) String() string {
	res := fmt.Sprintf("route: %s", m.Route)
	if m.Addrs != nil {
		res += fmt.Sprintf("addrs: %v", m.Addrs)
	}
	res += fmt.Sprintf("RedisClientCommonOptions: %+v", m.RedisClientCommonOptions)
	if m.redisClusterOpts != nil {
		res += fmt.Sprintf("redisClusterOpts: %+v", m.redisClusterOpts)
	}
	return res
}

func WithRedisClusterAddrs(addrs []string) option.Option {
	return option.NewOpt(func(op option.OptPrinter) {
		o, ok := op.(*RedisClusterClientOptions)
		if !ok {
			return
		}
		o.Addrs = addrs
	})
}

func WithRedisClusterRoute(route string) option.Option {
	return option.NewOpt(func(op option.OptPrinter) {
		o, ok := op.(*RedisClusterClientOptions)
		if !ok {
			return
		}

		switch route {
		case "randomly", "latency":
			o.Route = route
		default:
			o.Route = "randomly"
		}
	})
}

func WithGoRedisClusterOpts(opts *redis.ClusterOptions) option.Option {
	return option.NewOpt(func(op option.OptPrinter) {
		o, ok := op.(*RedisClusterClientOptions)
		if !ok {
			return
		}
		o.redisClusterOpts = opts
	})
}

type RedisClientCommonOptions struct {
	Password string `mapstructure:"password"`
	Username string `mapstructure:"username"`

	MaxRetries      int           `mapstructure:"maxRetries"`
	MinRetryBackoff time.Duration `mapstructure:"minRetryBackoff"`
	MaxRetryBackoff time.Duration `mapstructure:"maxRetryBackoff"`
	DialTimeout     time.Duration `mapstructure:"dialTimeout"`
	ReadTimeout     time.Duration `mapstructure:"readTimeout"`
	WriteTimeout    time.Duration `mapstructure:"writeTimeout"`

	// connect pool
	PoolSize     int           `mapstructure:"poolSize"`
	MinIdleConns int           `mapstructure:"minIdleConns"`
	MaxIdleConns int           `mapstructure:"maxIdleConns"`
	MaxConnAge   time.Duration `mapstructure:"maxConnAge"`
	PoolTimeout  time.Duration `mapstructure:"poolTimeout"`
	IdleTimeout  time.Duration `mapstructure:"idleTimeout"`
}

func DefaultRedisClusterClientOptions() *RedisClusterClientOptions {
	return &RedisClusterClientOptions{
		//Addrs:    []string{":26379"},
		Addrs: []string{":26379", ":26380", ":26381", ":26382", ":26383", ":26384"},
		RedisClientCommonOptions: RedisClientCommonOptions{
			Password: "",
			Username: "",

			MaxRetries:      3,
			MinRetryBackoff: 3 * time.Second,
			MaxRetryBackoff: 5 * time.Second,
			DialTimeout:     5 * time.Second,
			ReadTimeout:     3 * time.Second,
			WriteTimeout:    3 * time.Second,

			// connect pool
			PoolSize:     100,
			MinIdleConns: 10,
			MaxConnAge:   60 * time.Second,
			PoolTimeout:  5 * time.Second,
			IdleTimeout:  30 * time.Second,
		},
		Route: "randomly",
	}
}

// InitRedisClusterClient init redis cluster instance
func InitRedisClusterClient(options ...option.Option) redis.UniversalClient {
	opts := getClusterClientOptions(options...)
	clusterOpts := &redis.ClusterOptions{
		Addrs:    opts.Addrs,
		Password: opts.Password,
		Username: opts.Username,

		MaxRetries:      opts.MaxRetries,
		MinRetryBackoff: opts.MinRetryBackoff,
		MaxRetryBackoff: opts.MaxRetryBackoff,
		DialTimeout:     opts.DialTimeout,
		ReadTimeout:     opts.ReadTimeout,
		WriteTimeout:    opts.WriteTimeout,

		// connect pool
		PoolSize:        opts.PoolSize,
		MinIdleConns:    opts.MinIdleConns,
		MaxIdleConns:    opts.MaxIdleConns,
		PoolTimeout:     opts.PoolTimeout,
		ConnMaxLifetime: opts.MaxConnAge,
		ConnMaxIdleTime: opts.IdleTimeout,

		// To route commands by latency or randomly, enable one of the following.
		//RouteByLatency: true,
		//RouteRandomly: true,
	}
	switch opts.Route {
	case "randomly":
		clusterOpts.RouteRandomly = true
	case "latency":
		clusterOpts.RouteByLatency = true
	default:
		clusterOpts.RouteRandomly = true
	}

	if opts.redisClusterOpts != nil {
		clusterOpts = opts.redisClusterOpts
	}

	return redis.NewClusterClient(clusterOpts)
}

func getClusterClientOptions(opts ...option.Option) *RedisClusterClientOptions {
	options := DefaultRedisClusterClientOptions()
	for _, o := range opts {
		o.Apply(options)
	}

	return options
}
