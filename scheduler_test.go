package cronlib

import (
	"testing"
	"time"
)

// TestCronJob 自测
func TestCronJob(t *testing.T) {
	t.SkipNow()

	var job1Fn = func(jobID int) {
		t.Logf("exec job[#%d] fn", jobID)
	}

	// 准备job任务
	job1 := NewJob("job1", func() { job1Fn(1) }, "* * * * *")
	job2 := NewJob("job2", func() { job1Fn(2) }, "*/2 * * * *")

	// 准备scheduler调度器
	// locker, err := NewRedisLocker("redis://127.0.0.1:6379/0?dial_timeout=3&db=1&read_timeout=6s&max_retries=2")
	// if err != nil {
	// 	t.Fatalf("cron new redis locker got err: %s", err)
	// }
	cron := NewScheduler(&Config{
		TimeZone:         "Asia/Shanghai",
		Async:            false,
		SingletonModeAll: true,
	}, nil)

	// 添加job任务
	if err := cron.AddJobs(job1, job2); err != nil {
		t.Fatalf("cron add job got err: %s", err)
	}

	// job执行
	cron.Start()
	// Output:
	//  add job[job1] success...
	//  add job[job2] success...
	//  start cron jobs ...
}

// TestAddJob 自测
func TestAddJob(t *testing.T) {
	t.SkipNow()

	var job1Fn = func(jobID int) {
		t.Logf("exec job[#%d] fn", jobID)
	}

	// 准备job任务
	job1 := NewJob("job1", func() { job1Fn(1) }, "* * * * *")
	job2 := NewJob("job2", func() { job1Fn(2) }, "* * * * *")

	// 准备scheduler调度器
	cron := NewScheduler(&Config{
		TimeZone:         "Asia/Shanghai",
		Async:            true,
		SingletonModeAll: true,
	}, nil)

	// 启动job执行
	cron.Start()

	// 添加job1任务
	if err := cron.AddJobs(job1); err != nil {
		t.Fatalf("cron add job got err: %s", err)
	}

	// 添加job2任务
	if err := cron.AddJobs(job2); err != nil {
		t.Fatalf("cron add job got err: %s", err)
	}

	time.Sleep(2 * time.Minute)
}
