package consumerapi

import (
	"encoding/json"
	"log"
	"net/http"
)

func (s *Service) eventConsumerHandler(w http.ResponseWriter, r *http.Request) {
	// Read the body into a string for json decoding
	var content = []map[string]interface{}{}
	err := json.NewDecoder(r.Body).Decode(&content)
	if err != nil {
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
	for index, _ := range content {
		// let's create a job with the payload
		work := Job{Payload: &content[index]}
		// Push the work onto the queue.
		s.Dispatcher.JobQueue <- work
	}

	w.WriteHeader(http.StatusAccepted)
}
