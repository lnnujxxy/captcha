//
// @Author: zhouweiwei
// @Date: 2020/7/11 11:47 上午
//

package captcha

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

type RedisOption struct {
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var defaultRedisOption = RedisOption{
	MaxIdle:     10,
	MaxActive:   10,
	IdleTimeout: 30 * time.Second,
}

type FnOption func(option *RedisOption)

func WithMaxIdle(maxIdle int) FnOption {
	return func(option *RedisOption) {
		option.MaxIdle = maxIdle
	}
}

func WithMaxActive(maxActive int) FnOption {
	return func(option *RedisOption) {
		option.MaxActive = maxActive
	}
}

func WithIdleTimeout(idleTimeout time.Duration) FnOption {
	return func(option *RedisOption) {
		option.IdleTimeout = idleTimeout
	}
}

type redisStore struct {
}

var RedisConn *redis.Pool

func NewRedis(host, password string, opts ...FnOption) Store {
	option := &defaultRedisOption
	// 调用设置属性
	for _, opt := range opts {
		opt(option)
	}

	rs := &redisStore{}
	RedisConn = &redis.Pool{
		MaxIdle:     option.MaxIdle,
		MaxActive:   option.MaxActive,
		IdleTimeout: option.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return rs
}

func (r *redisStore) Set(id string, digits []byte) {
	conn := RedisConn.Get()
	defer conn.Close()

	conn.Do("SET", id, digits)
	conn.Do("EXPIRE", id, 3600)
}

func (r *redisStore) Get(id string, clear bool) (digits []byte) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", id))
	if err != nil {
		return nil
	}
	if clear {
		conn.Do("DEL", id)
	}

	return reply
}
