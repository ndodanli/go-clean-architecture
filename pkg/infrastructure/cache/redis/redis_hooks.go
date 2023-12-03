package redissrv

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"net"
	"time"
)

type RedisHook struct{}

func (l *RedisHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

func (l *RedisHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		return next(ctx, cmd)
	}
}

func (l *RedisHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		return next(ctx, cmds)
	}
}

func (l *RedisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (l *RedisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	return nil
}

func (l *RedisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (l *RedisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}

func (l *RedisHook) OnBeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (l *RedisHook) OnAfterProcess(ctx context.Context, cmd redis.Cmder) error {
	return nil
}

func (l *RedisHook) OnBeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (l *RedisHook) OnAfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}

func (l *RedisHook) OnConnect(ctx context.Context, cn *redis.Conn) error {
	log.Println("Connected to Redis")
	return nil
}

func (l *RedisHook) OnDisconnect(ctx context.Context, err error) {
	log.Printf("Disconnected from Redis: %v\n", err)
}

func (l *RedisHook) OnReconnect(ctx context.Context, n int, delay time.Duration) {
	log.Printf("Reconnecting to Redis (attempt %d) after %s\n", n, delay)
}
