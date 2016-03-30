package consumerapi

import (
	"github.com/efimovalex/EventKitAPI/adaptors/database"
	"log"
)

// Job represents the job to be run
type Job struct {
	Payload map[string]interface{}
}

// A buffered channel that we can send work requests on.
var JobQueue chan Job

// Worker represents the worker that executes the job
type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
	logger     *log.Logger
	dbAdaptor  *database.Adaptor
}

func NewWorker(workerPool chan chan Job, logger *log.Logger, dbAdaptor *database.Adaptor) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
		logger:     logger,
		dbAdaptor:  dbAdaptor,
	}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w *Worker) Start() {
	go func() {
		w.logger.Println("INFO: worker started")
		// register the current worker into the worker queue.
		w.WorkerPool <- w.JobChannel
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel
			select {
			case job := <-w.JobChannel:
				// we have received a work request.
				err := w.dbAdaptor.AddEvent(job.Payload)
				if err != nil {
					w.logger.Println(err.Error())
				}
				// register the current worker into the worker queue.
				w.WorkerPool <- w.JobChannel
			case <-w.quit:
				w.logger.Println("worker quit")
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w *Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
