# cronlib

`cronlib`基于`github.com/go-co-op/gocron`的基础上，添加了redis分布式锁，并简化了定时任务调度器的使用；

基于分布式锁的操作，用以解决在集群模式下，多cron可能在同一时刻发起调度的问题

## cronlib 快速使用流程
1. 初始化Job业务逻辑：定义Job定期执行具体业务逻辑
2. 初始化Job任务: `NewJob()`方法创建调度job
3. 初始化Scheduler调度器
4. 添加Job任务到调度器
5. 调度器启动

### 简单的定时任务示例

```
func ExampleStartTest() {
	var job1Fn = func(jobID int) {
		log.Printf("exec job[#%d] fn", jobID)
	}

	// 准备job任务
	job1 := NewJob("job1", func() { job1Fn(1) }, "* * * * *")
	job2 := NewJob("job2", func() { job1Fn(2) }, "*/2 * * * *")

	// 准备scheduler调度器
	locker, err := NewRedisLocker("redis://127.0.0.1:6379")
	if err != nil {
		log.Fatalf("parse redis url got err: %s", err)
	}

	crond := NewScheduler(&Config{
		Async:            false,           // 不阻塞主协程
		SingletonModeAll: true,            // 调度器不会重复调度同类型新的job任务
	}, locker)

	// 添加job任务
	err = crond.AddJobs(job1, job2)
	if err != nil {
		log.Fatalf("crond add job got err: %s", err)
	}

	// job执行
	crond.Start()
}
```

### redis 依赖

通过docker快速启动一个本地redis实例：

```shell
$ docker-compose up -d
```