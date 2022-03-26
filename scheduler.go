package cronlib

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/pkg/errors"
)

// cron格式`*/1 * * * *`，分(0-59) 时(0-23) 日(1-31) 月(1-12) 天(0-6)
const (
	ExpressEveryMin   = "* * * * *"
	ExpressEvery5Min  = "*/5 * * * *"
	ExpressEvery10Min = "*/10 * * * *"
	ExpressEvery30Min = "*/30 * * * *"
	ExpressEveryHour  = "0 * * * *" // 每小时整点执行
	ExpressEveryDay   = "0 1 * * *" // 凌晨1:00执行
	ExpressEveryMonth = "0 1 1 * *" // 每个月
)

// Scheduler Cron调度程序
type Scheduler struct {
	cfg             *Config           // 调度器配置
	gocronScheduler *gocron.Scheduler // gocron 调度器
	locker          Locker            // 在分布式模式cron job并发执行时，仅在抢到锁的执行真正的任务
}

// Config Cron调度器配置
type Config struct {
	TimeZone         string // 时区配置，默认为Asia/Shanghai
	Async            bool   // 启动job的方式是阻塞还是非阻塞
	SingletonModeAll bool   // 启动job是否采用单例模式，单例模式下若job如果之前有运行且未完成，则调度器不会重复调度同类型新的job任务(若无特殊要求，推荐开启)
}

// NewScheduler 初始化一个cron 调度器
func NewScheduler(cfg *Config, locker Locker) *Scheduler {
	scheduler := gocron.NewScheduler(parseTimeZone(cfg.TimeZone))

	// singleton model
	if cfg.SingletonModeAll {
		scheduler.SingletonModeAll()
	}

	return &Scheduler{
		cfg:             cfg,
		gocronScheduler: scheduler,
		locker:          locker,
	}
}

// AddJobs 新增一个cron job任务
func (s *Scheduler) AddJobs(jobs ...*Job) error {
	for _, job := range jobs {
		cronJob, err := s.gocronScheduler.
			Cron(job.express).
			Do(wrapLockerForJob(s.locker, job))
		if err != nil {
			return errors.Wrapf(err, "[error] gocronScheduler add job got err")
		}

		log.Printf("add job[%s] success...", job.name)
		cronJob.Tag(job.name)
	}
	return nil
}

// Start 启动Cron定时服务
func (s *Scheduler) Start() {
	// 是否启动非阻塞任务
	log.Println("start cron jobs ...")
	if s.cfg.Async {
		s.gocronScheduler.StartAsync()
	} else {
		s.gocronScheduler.StartBlocking()
	}
}

// parseTimeZone 解析调度的时区，默认时区为 Asia/Shanghai
func parseTimeZone(name string) *time.Location {
	if name == "" {
		name = "Asia/Shanghai"
	}
	cronTimeLoc, err := time.LoadLocation(name)
	if err != nil {
		log.Fatalf("parse cron timezone name got err: %s", err)
	}
	return cronTimeLoc
}

// wrapLockerForJob 对Job任务包裹分布式锁
func wrapLockerForJob(locker Locker, job *Job) func() {
	return func() {
		if locker != nil {
			// 抢分布式锁，延时30s自动解锁，限定cron调度的并发
			if err := locker.Lock(job.name, 30*time.Second); err != nil {
				log.Printf("locker required, get cron locker fail, %s", err)
				return
			}

			// 获得锁的执行任务
			log.Printf("locker required, get cron locker success, exec cron job [#%s]...", job.name)
			job.exec()
			return
		}

		// 无锁直接执行
		log.Printf("no locker required, exec cron job [#%s]...", job.name)
		job.exec()
	}
}
