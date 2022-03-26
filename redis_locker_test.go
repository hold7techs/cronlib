package cronlib

import (
	"context"
	"testing"
	"time"
)

// TestRedisLocker_Lock 自测
func TestRedisLocker_Lock(t *testing.T) {
	// 本地测试依赖redis实例
	t.SkipNow()

	redisURL := "redis://127.0.0.1:6379/0?dial_timeout=3&read_timeout=6s&max_retries=2"
	locker, err := NewRedisLocker(redisURL)
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
