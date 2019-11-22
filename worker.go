package workerpool

import "context"

// JobFunc ...
type JobFunc func() error

// PanicHandler ...
type PanicHandler func(interface{})

// Job ...
type Job struct {
	done <-chan struct{}
	f    JobFunc
	hd   PanicHandler
}

// Worker ...
type Worker struct {
	pool *Pool
}

func newWorker(p *Pool) *Worker {
	return &Worker{
		pool: p,
	}
}

// Do worker do a job
func (w *Worker) Do(j *Job) {
	defer func() {
		w.pool.workers <- w
		if err := recover(); err != nil && j.hd != nil {
			j.hd(err)
		}
	}()
	j.Do()
}

func jobWithCtx(ctx context.Context, j JobFunc, hd PanicHandler) *Job {
	return &Job{
		done: ctx.Done(),
		f:    j,
		hd:   hd,
	}
}

// Do ...
func (j *Job) Do() error {
	select {
	case <-j.done:
		return ErrJobIsCanceled
	default:
		return j.f()
	}
}
