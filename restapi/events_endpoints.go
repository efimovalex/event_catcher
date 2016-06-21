package restapi

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"

	"github.com/efimovalex/EventKitAPI/adaptors/cache"
	"github.com/efimovalex/EventKitAPI/adaptors/database"
	"github.com/efimovalex/EventKitAPI/common"
)

type EventsEndpoints struct {
	DBAdaptor    *database.Adaptor
	CacheAdaptor *cache.Adaptor
}

const Timestamp = "day"

const date = "2006-01-02"

// ControllerError structure returned to the client
type ControllerError struct {
	Message string      `json:"message"`
	Field   interface{} `json:"field,omitempty"`
}

func (ee *EventsEndpoints) Get(c echo.Context) error {
	headers := c.Response().Header()
	headers.Set("Cache-Control", "private, max-age=600")

	var nextPage string

	queryString := []byte(c.Request().URL().QueryString())

	cachedResponse := ee.CacheAdaptor.GetEventRequest(queryString)
	if cachedResponse != "" {
		log.Println("returning cached response")

		return c.JSONBlob(http.StatusOK, []byte(cachedResponse))
	}

	offsetID := c.QueryParam("offset_id")
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit == 0 {
		limit = 50
	}

	fieldName := c.QueryParam("field_name")

	var events []database.Event
	var getErr error

	if fieldName == Timestamp {
		startValue := c.QueryParam("start_date")
		if startValue == "" {
			return c.JSON(http.StatusBadRequest, ControllerError{
				Message: "Start date is required for fieldValue = day",
			})
		}
		startDate, timeErr := time.Parse(date, startValue)
		if timeErr != nil {
			log.Println("Could not convert start date: " + timeErr.Error())

			return c.JSON(http.StatusBadRequest, ControllerError{
				Message: "Could not convert start date ",
			})
		}

		endValue := c.QueryParam("end_date")
		if endValue != "" {
			endDate, timeErr := time.Parse(date, endValue)
			if timeErr != nil {
				log.Println("Could not convert end date: " + timeErr.Error())

				return c.JSON(http.StatusBadRequest, ControllerError{
					Message: "Could not convert end date ",
				})
			}
			events, getErr = ee.DBAdaptor.GetEventsInInterval(startDate, endDate, limit, offsetID)

		} else {
			events, getErr = ee.DBAdaptor.GetEventsInInterval(startDate, time.Now(), limit, offsetID)
		}
	} else {
		fieldValue := c.QueryParam("field_value")
		events, getErr = ee.DBAdaptor.GetEvents(fieldName, fieldValue, limit, offsetID)

		if len(events) == limit-1 {
			nextPage = fmt.Sprintf("%s%s?field_name=%s&field_value=%s&limit=%d&offset_id=%s", URL, c.Get("route"), fieldName, fieldValue, limit, events[limit-1].SGEventID)
		}
	}

	if getErr != nil {
		log.Println("error returning data from cassandra: " + getErr.Error())

		return c.JSON(http.StatusInternalServerError, ControllerError{
			Message: "Something went wrong retrieving data from storage",
		})
	}

	var response common.ListResponse

	if len(events) < limit {
		response = common.ListResponse{
			Events: events}
	} else {
		response = common.ListResponse{
			Events:   events,
			NextPage: nextPage,
		}
	}

	ee.CacheAdaptor.SaveEventRequest(queryString, response)

	return c.JSON(http.StatusOK, response)
}
