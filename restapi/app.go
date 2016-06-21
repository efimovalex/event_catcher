package restapi

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	mw "github.com/labstack/echo/middleware"

	"github.com/efimovalex/EventKitAPI/adaptors/cache"
	"github.com/efimovalex/EventKitAPI/adaptors/database"
)

var URL = ""

type Service struct {
	Logger    *log.Logger
	config    *Config
	DBAdaptor *database.Adaptor
	Router    *echo.Echo
}

type ServiceInterface interface {
	Start(*Config) error
}

// NewService loads configs and starts listeners
func NewService(config *Config) Service {
	// setup echo
	echo := echo.New()

	logger := log.New(os.Stderr, "EVENT CONSUMER:", log.Ldate|log.Ltime|log.Lshortfile)

	DBAdaptor := database.NewAdaptor(strings.Split(config.CassandraInterfaces, ","))

	CacheAdaptor := cache.NewAdaptor(config.CacheURL)

	service := Service{
		Logger: logger,
		config: config,
		Router: Routes(
			// add each dependent service as a dependency to the router
			dependencies{
				config:       config,
				DBAdaptor:    DBAdaptor,
				CacheAdaptor: CacheAdaptor,
				Router:       echo,
			},
		),
		DBAdaptor: DBAdaptor,
	}

	return service
}

// Start runs the entire service
func (s *Service) Start(config *Config) error {
	RESTError := make(chan error)
	go func() {
		RESTError <- s.StartHTTP()
	}()

	for {
		select {
		case RESTMSG := <-RESTError:
			close(RESTError)
			return RESTMSG
		}
	}
}

// StartHTTP listens on the configured ports for the REST application
func (s *Service) StartHTTP() error {
	address := fmt.Sprintf("%s:%d", s.config.Interface, s.config.Port)

	URL = address
	// Use middlewares
	s.Router.Use(mw.Gzip())
	s.Router.Use(mw.Logger())
	s.Router.Run(standard.New(address))

	return nil
}
