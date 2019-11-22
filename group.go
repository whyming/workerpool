package workerpool

import (
	"context"
	"sync"
	"time"
)

// Group is job group
type Group interface {
	AddJob(ctx context.Context, j JobFunc) error
	Done()
	SetTryTimes(int8) Group
	SetPanicHd(PanicHandler) Group
	SetInterval(time.Duration) Group
}

type group struct {
	pool     *Pool
	wg       sync.WaitGroup
	times    int8
	interval time.Duration
	hd       PanicHandler
}

// NewGroup new a group
func NewGroup(p *Pool) Group {
	return &group{
		pool: p,
	}
}
func (g *group) AddJob(ctx context.Context, job JobFunc) error {
	g.wg.Add(1)
	return g.addTryJob(ctx, newTryJob(job))
}

func (g *group) Done() {
	g.wg.Wait()
}

func (g *group) SetTryTimes(n int8) Group {
	g.times = n
	return g
}

func (g *group) SetPanicHd(hd PanicHandler) Group {
	g.hd = hd
	return g
}

func (g *group) SetInterval(d time.Duration) Group {
	g.interval = d
	return g
}

func (g *group) try(ctx context.Context, j tryJob) {
	if j.times > g.times {
		g.wg.Done()
		return
	}
	time.Sleep(g.interval)
	g.addTryJob(ctx, j)
}

func (g *group) addTryJob(ctx context.Context, job tryJob) error {
	j := func() error {
		var err error
		defer func() {
			if err == nil {
				g.wg.Done()
			}
		}()
		err = job.job()
		if err != nil && err != ErrJobIsCanceled {
			job.times++
			g.try(ctx, job)
		}
		return err
	}
	return g.pool.AddJob(ctx, j, g.hd)
}

type tryJob struct {
	job   JobFunc
	times int8
}

func newTryJob(j JobFunc) tryJob {
	return tryJob{
		job: j,
	}
}
