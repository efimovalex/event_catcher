package consumerapi

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/efimovalex/EventKitAPI/adaptors/database"

	_ "net/http/pprof"
	"runtime"
	"runtime/pprof"
)

type Service struct {
	Logger     *log.Logger
	config     *Config
	Dispatcher *Dispatcher
	DBAdaptor  *database.Adaptor
}

type ServiceInterface interface {
	Start(*Config) error
}

type myHandler struct {
}

var mux map[string]func(http.ResponseWriter, *http.Request)

func NewService(config *Config) Service {
	logger := log.New(os.Stderr, "EVENT CONSUMER:", log.Ldate|log.Ltime|log.Lshortfile)

	DBAdaptor := database.NewAdaptor(strings.Split(config.CassandraInterfaces, ","), config.CassandraUser, config.CassandraPassword)

	dispatcher := NewDispatcher(config.MaxWorker, config.MaxJobQueue, logger, DBAdaptor)
	return Service{
		Logger:     logger,
		config:     config,
		Dispatcher: dispatcher,
		DBAdaptor:  DBAdaptor,
	}
}

// Start starts listeners
func (s *Service) Start() error {
	s.Dispatcher.Run()

	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	f2, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f2); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
	f2.Close()

	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", s.config.Interface, s.config.Port),
		Handler:        &myHandler{},
		ReadTimeout:    2 * time.Second,
		WriteTimeout:   2 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/v1/events"] = s.eventConsumerHandler

	return server.ListenAndServe()
}

func (m *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// if r.Method != "POST" {
	// 	w.WriteHeader(http.StatusMethodNotAllowed)

	// 	return
	// }
	log.Println(r.URL.String())
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)

		return
	}
}
