package restapi

import (
	"github.com/labstack/echo"

	"github.com/efimovalex/EventKitAPI/adaptors/cache"
	"github.com/efimovalex/EventKitAPI/adaptors/database"
)

var Resources = map[string]string{
	"eventsEndpoint": "/v1/events",
	"eventEndpoint":  "/v1/event/:sg_event_id",
}

type dependencies struct {
	config       *Config
	DBAdaptor    *database.Adaptor
	CacheAdaptor *cache.Adaptor
	Router       *echo.Echo
}

// Routes returns an http.Handler with the available RESTful routes for the service
func Routes(deps dependencies) *echo.Echo {
	var resource string

	// place all routes here to make it easier to find
	eventsEndpoints := &EventsEndpoints{DBAdaptor: deps.DBAdaptor, CacheAdaptor: deps.CacheAdaptor}
	eventEndpoints := &EventEndpoints{DBAdaptor: deps.DBAdaptor, CacheAdaptor: deps.CacheAdaptor}

	resource = Resources["eventEndpoint"]
	deps.Router.Delete(resource, dispatch(resource, eventEndpoints.Delete))
	deps.Router.Get(resource, dispatch(resource, eventEndpoints.Get))

	resource = Resources["eventsEndpoint"]
	deps.Router.Get(resource, dispatch(resource, eventsEndpoints.Get))

	return deps.Router
}

// dispatch add route name in the context
func dispatch(route string, handler echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("route", route)

		return handler(c)
	}
}
