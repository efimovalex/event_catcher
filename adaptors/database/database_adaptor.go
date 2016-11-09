package database

import (
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
	"github.com/hailocab/gocassa"
)

// Adaptor for DB
type Adaptor struct {
	Session            gocassa.KeySpace
	urls               []string
	eventTimeDateRange time.Duration
	Username           string
	Password           string
}

type Interface interface {
}

var IndexFields = []string{"event", "email", "sg_message_id"}

// New db adaptor
func NewAdaptor(urls []string, username, password string) *Adaptor {
	cluster := gocql.NewCluster(urls...)
	cluster.Keyspace = "events"
	cluster.ProtoVersion = 3

	// cluster.Authenticator = gocql.PasswordAuthenticator{
	// 	Username: username,
	// 	Password: password,
	// }

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Cassandra conn error: %s", err.Error())
	}
	qe := gocassa.GoCQLSessionToQueryExecutor(session)
	conn := gocassa.NewConnection(qe)
	gocassaSession := conn.KeySpace("events")

	return &Adaptor{
		Session:            gocassaSession,
		urls:               urls,
		eventTimeDateRange: 7 * 24 * time.Hour,
		Username:           username,
		Password:           password,
	}
}

func (a *Adaptor) ReestablishConnection() {
	log.Println("Reestablishing cql connection")
	cluster := gocql.NewCluster(a.urls[0])
	cluster.ProtoVersion = 3
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: a.Username,
		Password: a.Password,
	}
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Cassandra conn error: %s", err.Error())
	}
	qe := gocassa.GoCQLSessionToQueryExecutor(session)
	a.Session = gocassa.NewConnection(qe).KeySpace("events")
}

func (a *Adaptor) AddEvent(eventMap map[string]interface{}) error {
	e := Event{}
	e.MapEvent(eventMap)

	return a.Save(e)
}

func (a *Adaptor) Save(e Event) error {
	for _, indexField := range IndexFields {
		table := a.GetEventMultimapTable(indexField)

		if err := table.Set(e).Run(); err != nil {
			log.Println(err.Error())

			return fmt.Errorf("Save error: %s", err.Error())
		}
	}

	table := a.GetEventMapTable()
	if err := table.Set(e).Run(); err != nil {
		log.Println(err.Error())

		return fmt.Errorf("Save error: %s", err.Error())
	}

	timeSeriesTable := a.GetTimeSeriesEventTable()
	if err := timeSeriesTable.Set(e).Run(); err != nil {
		log.Println(err.Error())

		return fmt.Errorf("Save error: %s", err.Error())
	}

	return nil
}

func (a *Adaptor) DeleteEvent(SGEventID string) error {
	table := a.GetEventMapTable()
	var e Event

	err := table.Read(SGEventID, &e).Run()
	if err != nil {
		a.ReestablishConnection()

		return fmt.Errorf("Delete error, event not found: %s", err.Error())
	}

	if err := table.Delete(e.SGEventID).Run(); err != nil {
		a.ReestablishConnection()

		log.Printf("Delete error: %s, couldn't delete event from map table\n", err.Error())
	}

	for _, indexField := range IndexFields {
		table := a.GetEventMultimapTable(indexField)

		if err := table.Delete(e.SGEventID, e.SGEventID).Run(); err != nil {
			a.ReestablishConnection()

			log.Printf("Delete error: %s, couldn't delete event with index %s \n", err.Error(), indexField)
		}
	}

	timeSeriesTable := a.GetTimeSeriesEventTable()
	if err := timeSeriesTable.Delete(e.Timestamp, e.SGEventID).Run(); err != nil {
		a.ReestablishConnection()

		log.Printf("Delete error: %s, couldn't delete time series event\n", err.Error())
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

	mapTable := a.GetEventMapTable()

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
