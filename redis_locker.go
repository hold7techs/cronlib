package cronlib

import (
	"context"
	"sync"
	"time"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type rdsLockInfo struct {
	ctx    context.Context
	key    string
	locker *redislock.Lock
}

// RedisLocker 基于Redis实现的分布式锁，
type RedisLocker struct {
	mu         *sync.Mutex
	rdsClient  *redis.Client
	lockClient *redislock.Client
	locks      map[string]*rdsLockInfo
}

// NewRedisLocker 创建一个Redis分布式锁
// Examples:
//		redis://user:password@localhost:6789/3?dial_timeout=3&db=1&read_timeout=6s&max_retries=2
//		is equivalent to:
//		&Options{
//			Network:     "tcp",
//			Addr:        "localhost:6789",
//			DB:          1,               // path "/3" was overridden by "&db=1"
//			DialTimeout: 3 * time.Second, // no time unit = seconds
//			ReadTimeout: 6 * time.Second,
//			MaxRetries:  2,
//		}
func NewRedisLocker(redisURL string) (*RedisLocker, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	rdsClient := redis.NewClient(opts)

	return &RedisLocker{
		mu:         &sync.Mutex{},
		rdsClient:  rdsClient,
		lockClient: redislock.New(rdsClient),
		locks:      make(map[string]*rdsLockInfo),
	}, nil
}

// Lock 分布式加锁
func (r *RedisLocker) Lock(key string, ttl time.Duration) error {
	ctx := context.Background()

	// cron不做重试
	locker, err := r.lockClient.Obtain(ctx, key, ttl, nil)
	if err != nil {
		return errors.Wrap(err, "redis locker obtain lock fail")
	}
	r.mu.Lock()
	r.locks[key] = &rdsLockInfo{
		ctx:    ctx,
		key:    key,
		locker: locker,
	}
	r.mu.Unlock()
	return nil
}

// UnLock 分布式解锁
func (r *RedisLocker) UnLock(key string) error {
	if l, ok := r.locks[key]; ok {
		ctx, locker := l.ctx, l.locker
		if err := locker.Release(ctx); err != nil {
			return err
		}

		// 锁释放后，从map移除key对应的锁信息
		delete(r.locks, key)
		return nil
	}
	return nil
}
