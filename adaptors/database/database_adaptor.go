package database

import (
	"fmt"
	"log"
	"time"

	"github.com/hailocab/gocassa"
)

// Adaptor for DB
type Adaptor struct {
	Session            gocassa.KeySpace
	urls               []string
	eventTimeDateRange time.Duration
}

type Interface interface {
}

var IndexFields = []string{"event", "email", "sg_message_id"}

// New db adaptor
func NewAdaptor(urls []string) *Adaptor {
	session, _ := gocassa.ConnectToKeySpace("events", urls, "", "")

	return &Adaptor{
		Session:            session,
		urls:               urls,
		eventTimeDateRange: 7 * 24 * time.Hour,
	}
}

func (a *Adaptor) ReestablishConnection() {
	log.Println("Reestablishing cql connection")
	a.Session, _ = gocassa.ConnectToKeySpace("events", a.urls, "", "")
}

func (a *Adaptor) AddEvent(eventMap map[string]interface{}) error {
	e := Event{}
	e.MapEvent(eventMap)

	return a.Save(e)
}

func (a *Adaptor) DeleteEvent(eventID string) error {
	e := Event{SGEventID: eventID}
	saveErr := a.Delete(e)

	return saveErr
}

func (a *Adaptor) Save(e Event) error {
	for _, indexField := range IndexFields {
		table := a.GetEventMultimapTable(indexField)

		if err := table.Set(e).Run(); err != nil {
			a.ReestablishConnection()

			return fmt.Errorf("Save error: %s", err.Error())
		}
	}

	table := a.GetTimeSeriesEventTable()
	if err := table.Set(e).Run(); err != nil {
		a.ReestablishConnection()

		return fmt.Errorf("Save error: %s", err.Error())
	}

	return nil
}

func (a *Adaptor) Delete(e Event) error {
	for _, indexField := range IndexFields {
		table := a.GetEventMultimapTable(indexField)

		if err := table.Delete(e, e.SGEventID).Run(); err != nil {
			a.ReestablishConnection()

			return fmt.Errorf("Delete error: %s", err.Error())
		}
	}

	table := a.GetTimeSeriesEventTable()
	if err := table.Delete(e.Timestamp, e.SGEventID).Run(); err != nil {
		a.ReestablishConnection()

		return fmt.Errorf("Delete error: %s", err.Error())
	}

	return nil
}

func (a *Adaptor) Update(e Event) error {
	for _, indexField := range IndexFields {
		table := a.GetEventMultimapTable(indexField)
		if err := table.Set(e).Run(); err != nil {
			a.ReestablishConnection()

			return fmt.Errorf("Update error: %s", err.Error())
		}
	}

	table := a.GetTimeSeriesEventTable()
	if err := table.Set(e).Run(); err != nil {
		a.ReestablishConnection()

		return fmt.Errorf("Update error: %s", err.Error())
	}

	return nil
}

func (a *Adaptor) GetEventsInInterval(startTime, endTime time.Time, limit int, offsetID string) ([]Event, error) {
	table := a.GetTimeSeriesEventTable()

	var result []Event

	if err := table.List(startTime, endTime, &result).Run(); err != nil {
		return []Event{}, fmt.Errorf("Get interval error: %s", err.Error())
	}

	return result, nil
}

func (a *Adaptor) GetEvents(field string, fieldValue interface{}, limit int, offsetID string) ([]Event, error) {
	if !stringInSlice(field, IndexFields) {
		return []Event{}, fmt.Errorf("Index field not found: %s", field)
	}
	var events []Event

	mapTable := a.GetEventMultimapTable(field)

	list := mapTable.List(fieldValue, offsetID, limit, &events)

	if offsetID == "" {
		list = mapTable.List(fieldValue, nil, limit, &events)
	}

	err := list.Run()
	if err != nil {
		a.ReestablishConnection()
		return []Event{}, err
	}

	return events, nil
}

func (a *Adaptor) GetEvent(SGEventID string) (Event, error) {
	var event Event

	mapTable := a.Session.MapTable("events", "sg_event_id", &Event{})

	list := mapTable.Read(SGEventID, &event)

	err := list.Run()
	if err != nil {
		a.ReestablishConnection()
		return Event{}, err
	}

	return event, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
