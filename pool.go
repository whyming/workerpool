package workerpool

import (
	"context"
	"errors"
	"sync"
)

// All Error
var (
	ErrWorkerPoolFull = errors.New("worker pool is full")
	ErrJobIsCanceled  = errors.New("job is canceled")
)

// Pool worker pool
type Pool struct {
	workers chan *Worker
	queue   chan *Job
	closed  chan struct{}
	mu      sync.RWMutex
}

// NewPool new a pool
func NewPool(ws int64, maxq int64) *Pool {
	p := &Pool{
		workers: make(chan *Worker, ws),
		queue:   make(chan *Job, maxq),
	}
	for i := int64(0); i < ws; i++ {
		p.workers <- newWorker(p)
	}
	go p.start()
	return p
}

func (p *Pool) start() {
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
			worker.Do(job)
		}
	}
}

// Stop stop worker pool
func (p *Pool) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed != nil {
		close(p.closed)
		p.closed = nil
	}
}

// AddJob add job to queue
func (p *Pool) AddJob(ctx context.Context, j JobFunc, hd PanicHandler) error {
	job := jobWithCtx(ctx, j, hd)
	select {
	case <-ctx.Done():
		return ErrJobIsCanceled
	default:
		select {
		case p.queue <- job:
			return nil
		default:
			return ErrWorkerPoolFull
		}
	}
}

// NewGroup ...
func (p *Pool) NewGroup() Group {
	return NewGroup(p)
}
