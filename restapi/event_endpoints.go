package restapi

import (
	"log"
	"net/http"

	"github.com/labstack/echo"

	"github.com/efimovalex/EventKitAPI/adaptors/cache"
	"github.com/efimovalex/EventKitAPI/adaptors/database"
)

type EventEndpoints struct {
	DBAdaptor    *database.Adaptor
	CacheAdaptor *cache.Adaptor
	Config       *Config
}

func (ee *EventEndpoints) Delete(c echo.Context) error {
	eventID := c.Param("sg_event_id")

	err := ee.DBAdaptor.DeleteEvent(eventID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ControllerError{
			Message: err.Error(),
		})
	}

	return c.NoContent(http.StatusOK)
}

func (ee *EventEndpoints) Get(c echo.Context) error {
	eventID := c.Param("sg_event_id")

	if ee.Config.EnableCaching {
		cachedResponse := ee.CacheAdaptor.GetEvent(eventID)
		if cachedResponse != "" {
			log.Println("returning cached response")

			return c.JSONBlob(http.StatusOK, []byte(cachedResponse))
		}
	}

	event, err := ee.DBAdaptor.GetEvent(eventID)
	if err != nil {
		log.Println("error returning data from cassandra: " + err.Error())

		return c.JSON(http.StatusInternalServerError, ControllerError{
			Message: "Something went wrong retrieving data from storage",
		})
	}

	return c.JSON(http.StatusOK, event)
}
