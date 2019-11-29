package workerpool

import (
	"context"
	"time"
)

// JobFunc ...
type JobFunc func() error

// Job ...
type Job struct {
	pool   Pool
	ctx    context.Context
	f      JobFunc
	conf   *Conf
	rtimes int8          // retry times
	done   chan struct{} //
}

func jobWithCtx(ctx context.Context, j JobFunc, conf *Conf) *Job {
	return &Job{
		ctx:  ctx,
		done: make(chan struct{}),
		f:    j,
		conf: conf,
	}
}

// Do ...
func (j *Job) Do(p *pool) error {
	defer func() {
		if err := recover(); err != nil && j.conf.PanicHD != nil {
			j.conf.PanicHD(err)
		}
	}()
	select {
	case <-j.ctx.Done():
		return j.ctx.Err()
	case <-j.done:
		return nil
	default:
		if err := j.f(); err != nil {
			if (j.conf.Retry == nil || j.conf.Retry(err)) && j.rtimes < j.conf.RTimes {
				j.rtimes++
				time.Sleep(j.conf.Interval)
				p.addRetryJob(j.ctx, j)
			} else {
				close(j.done)
				return err
			}
		} else {
			close(j.done)
		}
	}
	return nil
}
