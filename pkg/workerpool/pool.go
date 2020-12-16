package workerpool

import (
	"context"
	"sync"
)

type Job interface {
	// Process is any function doing some job
	Process(ctx context.Context)
}

type WorkerPool struct {
	poolSize int
	wg       *sync.WaitGroup
}

func NewWorkerPool(poolSize int) *WorkerPool {
	cp := WorkerPool{
		wg:       &sync.WaitGroup{},
		poolSize: poolSize,
	}

	return &cp
}

func (p *WorkerPool) Run(ctx context.Context, jobs []Job) []Job {
	jobsBucket := make(chan Job, len(jobs))

	for _, job := range jobs {
		jobsBucket <- job
	}

	close(jobsBucket)

	for i := 0; i < p.poolSize; i++ {
		p.wg.Add(1)

		go func(ctx context.Context, jobsBucket chan Job) {
			for job := range jobsBucket {
				select {
				case <-ctx.Done():
					return
				default:
					job.Process(ctx)
				}
			}

			p.wg.Done()
		}(ctx, jobsBucket)
	}

	p.wg.Wait()

	return jobs
}
