package cronlib

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

// TestRedisLocker_Lock 自测
func TestRedisLocker_Lock(t *testing.T) {
	// 本地测试依赖redis实例
	t.SkipNow()

	locker, err := NewRedisLocker("redis://127.0.0.1:6379")
	if err != nil {
		t.Fatalf("cron new redis locker got err: %s", err)
	}

	type args struct {
		key string
		ttl time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"lock01", args{"job1", 30 * time.Second}, false},
		{"lock02", args{"job1", 30 * time.Second}, true}, // 再次申请锁会失败
		{"lock03", args{"job1", 30 * time.Second}, true}, // 再次申请锁会失败
	}

	r := locker
	r.rdsClient.Del(context.Background(), "job1")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.Lock(tt.args.key, tt.args.ttl); (err != nil) != tt.wantErr {
				t.Errorf("Lock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}

// BenchmarkRedisLockerLock 并发测试，依赖redis
func BenchmarkRedisLockerLock(b *testing.B) {
	b.SkipNow()

	rdsLocker, err := NewRedisLocker("redis://127.0.0.1:6379")
	if err != nil {
		b.Fatal("redis connect error")
	}

	for i := 0; i < b.N; i++ {
		rdsLocker.Lock(strconv.FormatUint(rand.Uint64(), 10), 10*time.Millisecond)
	}
}

// TestRedisLockerLock 数据竟态检测: go test -run=TestRedisLockerLock -race -v
func TestRedisLockerLock(t *testing.T) {
	t.SkipNow()

	rand.Seed(time.Now().UnixNano())
	rdsLocker, err := NewRedisLocker("redis://127.0.0.1:6379")
	if err != nil {
		t.Fatal("redis connect error")
	}

	wg := sync.WaitGroup{}
	sema := make(chan bool, 200)
	for i := 0; i < 1e4; i++ {
		wg.Add(1)
		sema <- true
		go func() {
			key := fmt.Sprintf("key:%d", rand.Uint64())
			rdsLocker.Lock(key, 10*time.Millisecond)
			wg.Done()
			<-sema
		}()
	}
}
