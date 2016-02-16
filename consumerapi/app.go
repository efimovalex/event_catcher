package consumerapi

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Service struct {
	Logger     *log.Logger
	config     *Config
	Dispatcher *Dispatcher
}

type ServiceInterface interface {
	Start(*Config) error
}

type myHandler struct {
}

var mux map[string]func(http.ResponseWriter, *http.Request)

func NewService(config *Config) Service {
	logger := log.New(os.Stderr, "EVENT CONSUMER:", log.Ldate|log.Ltime|log.Lshortfile)
	dispatcher := NewDispatcher(config.MaxWorker, config.MaxJobQueue, logger)

	return Service{
		Logger:     logger,
		config:     config,
		Dispatcher: dispatcher,
	}
}

// Start starts listeners
func (s *Service) Start() error {
	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", s.config.Interface, s.config.Port),
		Handler:        &myHandler{},
		ReadTimeout:    2 * time.Second,
		WriteTimeout:   2 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/v1/event"] = s.eventConsumerHandler

	s.Dispatcher.Run()

	return server.ListenAndServe()
}

func (m *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)

		return
	}
}
