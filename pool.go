package workerpool

import (
	"context"
	"errors"
	"sync"
)

// All Error
var (
	ErrWorkerPoolFull = errors.New("worker pool is full")
	ErrWorkerPoolStop = errors.New("worker pool is stoped")
)

// Pool ...
type Pool interface {
	SetConf(*Conf) Pool
	AddJob(context.Context, JobFunc) error
	addJob(context.Context, *Job) error
	addRetryJob(context.Context, *Job) error
	NewGroup() Group
	Stop()
}

// Pool worker pool
type pool struct {
	workers chan *Worker
	queue   chan *Job // job queue
	req     chan *Job // retry job queue
	closed  chan struct{}
	mu      sync.RWMutex
	conf    *Conf
}

// NewPool new a pool
func NewPool(ws int64, maxq int64) Pool {
	p := &pool{
		workers: make(chan *Worker, ws),
		queue:   make(chan *Job, maxq),
		req:     make(chan *Job),
		closed:  make(chan struct{}),
		conf:    &Conf{},
	}
	for i := int64(0); i < ws; i++ {
		p.workers <- newWorker(p)
	}
	go p.start()
	return p
}

func (p *pool) SetConf(conf *Conf) Pool {
	p.conf = conf
	return p
}

func (p *pool) start() {
	p.mu.Lock()
	if p.closed == nil {
		p.closed = make(chan struct{})
	}
	p.mu.Unlock()
	for {
		select {
		case <-p.closed:
			return
		case job := <-p.queue:
			worker := <-p.workers
			worker.jobq <- job
		case job := <-p.req:
			worker := <-p.workers
			worker.jobq <- job
		}
	}
}

// Stop stop worker pool
func (p *pool) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed != nil {
		close(p.closed)
		p.closed = nil
	}
}

// AddJob add job to queue
func (p *pool) AddJob(ctx context.Context, job JobFunc) error {
	return p.addJob(ctx, jobWithCtx(ctx, job, p.conf))
}

func (p *pool) addJob(ctx context.Context, job *Job) error {
	select {
	case p.queue <- job:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-p.closed:
		return ErrWorkerPoolStop
	default:
		return ErrWorkerPoolFull
	}
}

func (p *pool) addRetryJob(ctx context.Context, job *Job) error {
	select {
	case p.req <- job:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-p.closed:
		return ErrWorkerPoolStop
	}
}

// NewGroup ...
func (p *pool) NewGroup() Group {
	return NewGroup(p)
}
