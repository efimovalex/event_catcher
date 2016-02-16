package consumerapi

import (
	"log"
)

type Dispatcher struct {
	JobQueue chan Job
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Job
	maxWorkers int
	logger     *log.Logger
}

func NewDispatcher(maxWorkers int, maxJobs int, logger *log.Logger) *Dispatcher {
	jobQueue := make(chan Job, maxJobs)
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{
		JobQueue:   jobQueue,
		WorkerPool: pool,
		logger:     logger,
	}
}

func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.WorkerPool, d.logger)
		worker.Start()
	}

	d.logger.Println("Workers started")

	d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.JobQueue:
			// a job request has been received
			d.logger.Println("job received")
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}
