package workerpool

import (
	"context"
	"sync"
)

// Group is job group
type Group interface {
	AddJob(ctx context.Context, j JobFunc) error
	SetConf(*Conf) Group
	Done()
}

type group struct {
	pool Pool
	wg   sync.WaitGroup
	conf *Conf
}

// NewGroup new a group
func NewGroup(p Pool) Group {
	return &group{
		pool: p,
		conf: &Conf{},
	}
}
func (g *group) AddJob(ctx context.Context, job JobFunc) error {
	g.wg.Add(1)
	return g.pool.addJob(ctx, g.wrappJob(ctx, job))
}

func (g *group) Done() {
	g.wg.Wait()
}

func (g *group) SetConf(conf *Conf) Group {
	g.conf = conf
	return g
}

func (g *group) wrappJob(ctx context.Context, job JobFunc) *Job {
	j := jobWithCtx(ctx, job, g.conf)
	go func() {
		<-j.done
		g.wg.Done()
	}()
	return j
}
