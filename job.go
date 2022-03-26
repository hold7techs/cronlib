package cronlib

// Job 待执行的JOB任务
type Job struct {
	name    string // job描述名称
	exec    func() // job调度后执行的内容
	express string // cron格式`*/1 * * * *`，分(0-59) 时(0-23) 日(1-31) 月(1-12) 天(0-6)
}

// NewJob 创建一个任务
func NewJob(name string, exec func(), express string) *Job {
	return &Job{
		name:    name,
		exec:    exec,
		express: express,
	}
}
