package redis

import (
	"app/pkg/config"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis interface {
	Ping(ctx context.Context) error
	Get(ctx context.Context, name string) (string, error)
	Set(ctx context.Context, name string, value string) error
	SetWithDuration(ctx context.Context, name string, value string, d time.Duration) error
	Delete(ctx context.Context, name string) error
	PrintKeys(ctx context.Context)
	Exists(ctx context.Context, name string) (bool, error)
	Incr(ctx context.Context, name string) (int64, error)
}

// Setup Redis
func NewRedis(rdc *config.RedisConfig, expiracy int) Redis {
	redis := redis.NewClient(&redis.Options{
		Addr:     rdc.Host + ":" + rdc.Port,
		Password: rdc.Password,
		DB:       0,
	})

	return &rds{
		rdb:      redis,
		expiracy: time.Duration(expiracy) * time.Second,
		prefix:   rdc.Prefix,
	}
}

type rds struct {
	rdb      *redis.Client
	expiracy time.Duration
	prefix   string
}

func (c *rds) PrintKeys(ctx context.Context) {
	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = c.rdb.Scan(ctx, cursor, "", 0).Result()
		if err != nil {
			panic(err)
		}

		for _, key := range keys {
			fmt.Println("key", key)
		}

		if cursor == 0 {
			break
		}
	}
}

func (c *rds) SetWithDuration(ctx context.Context, name string, value string, d time.Duration) error {
	return c.rdb.Set(ctx, c.prefix+"_"+name, value, d).Err()
}

func (c *rds) Set(ctx context.Context, name string, value string) error {
	return c.rdb.Set(ctx, c.prefix+"_"+name, value, c.expiracy).Err()
}

func (c *rds) Get(ctx context.Context, name string) (string, error) {
	return c.rdb.Get(ctx, c.prefix+"_"+name).Result()
}

func (c *rds) Delete(ctx context.Context, name string) error {
	return c.rdb.Del(ctx, c.prefix+"_"+name).Err()
}

func (c *rds) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

func (c *rds) Exists(ctx context.Context, name string) (bool, error) {
	exists, err := c.rdb.Exists(ctx, c.prefix+"_"+name).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (c *rds) Incr(ctx context.Context, name string) (int64, error) {
	return c.rdb.Incr(ctx, c.prefix+"_"+name).Result()
}
