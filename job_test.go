package cronlib

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
)

// Test_wrapLocker 自测
func Test_wrapLocker(t *testing.T) {
	var task = func() { t.Logf("exec task") }
	type args struct {
		locker Locker
		job    *Job
	}

	// mock
	ctrl := gomock.NewController(t)
	locker := NewMockLocker(ctrl)
	locker.EXPECT().Lock(gomock.Any(), gomock.Any()).MaxTimes(1).Return(nil)
	locker.EXPECT().Lock(gomock.Any(), gomock.Any()).AnyTimes().Return(errors.New("get lock fail"))
	locker.EXPECT().UnLock(gomock.Any()).AnyTimes()

	tests := []struct {
		name string
		args args
	}{
		{"t1", args{locker, NewJob("job1", task, "")}},
		{"t2", args{locker, NewJob("job1", task, "")}},
		{"t3", args{locker, NewJob("job1", task, "")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warpFn := wrapLockerForJob(tt.args.locker, tt.args.job)
			warpFn()
		})
	}
}
