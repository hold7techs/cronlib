package cronlib

import (
	"time"
)

//go:generate mockgen -source=./locker.go -package cron -destination locker_mock.go

// Locker 分布式锁接口
// 该接口主要是针对多个JOB执行时刻，进行并发控制，确保在不同机器上的Job仅执行一次行
type Locker interface {
	// Lock 分布式加锁
	Lock(key string, ttl time.Duration) error

	// UnLock 分布式解锁
	UnLock(key string) error
}
