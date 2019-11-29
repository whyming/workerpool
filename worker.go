package workerpool

// Worker ...
type Worker struct {
	pool *pool
	jobq chan *Job
}

func newWorker(p *pool) *Worker {
	w := &Worker{
		pool: p,
		jobq: make(chan *Job),
	}
	go w.start()
	return w
}

func (w *Worker) start() {

	for {
		select {
		case job := <-w.jobq:
			job.Do(w.pool)
			w.pool.workers <- w
		}
	}
}
