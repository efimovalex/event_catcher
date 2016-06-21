package consumerapi

import (
	"github.com/efimovalex/EventKitAPI/adaptors/database"
	"log"
)

type Dispatcher struct {
	JobQueue chan Job
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Job
	maxWorkers int
	logger     *log.Logger
	dbAdaptor  *database.Adaptor
}

func NewDispatcher(maxWorkers int, maxJobs int, logger *log.Logger, dbAdaptor *database.Adaptor) *Dispatcher {
	jobQueue := make(chan Job, maxJobs)
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{
		maxWorkers: maxWorkers,
		JobQueue:   jobQueue,
		WorkerPool: pool,
		logger:     logger,
		dbAdaptor:  dbAdaptor,
	}
}

func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(i, d.WorkerPool, d.logger, d.dbAdaptor)
		worker.Start()
	}

	d.logger.Println("INFO: Workers started")

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.JobQueue:
			// a job request has been received
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

func (d *Dispatcher) AddJob(work Job) {
	d.JobQueue <- work
}
