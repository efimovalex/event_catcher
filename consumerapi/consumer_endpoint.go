package consumerapi

import (
	"encoding/json"

	"log"
	"net/http"
)

func (s *Service) eventConsumerHandler(w http.ResponseWriter, r *http.Request) {
	// defer close of body
	if r.Body != nil {
		defer r.Body.Close()
	}
	// Read the body into a string for json decoding
	var content []map[string]interface{}
	d := json.NewDecoder(r.Body)
	d.UseNumber()
	if err := d.Decode(&content); err != nil {
		log.Printf("Error decoding request body: %s", err.Error())
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if len(s.Dispatcher.JobQueue) == cap(s.Dispatcher.JobQueue) {
		log.Println("ERROR: Job Queue is full")
		w.WriteHeader(http.StatusServiceUnavailable)

		return
	}

	// Go through each payload and queue items individually to be posted to Cassandra
	for _, job := range content {
		// let's create a job with the payload
		work := Job{Payload: job}
		// Push the work onto the queue.
		s.Dispatcher.AddJob(work)
	}

	w.WriteHeader(http.StatusAccepted)
}
