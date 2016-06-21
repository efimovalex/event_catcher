package database

import (
	"time"

	"github.com/hailocab/gocassa"
)

func (a *Adaptor) GetEventMultimapTable(index string) gocassa.MultimapTable {
	return a.Session.MultimapTable("events", index, "sg_event_id", &Event{})
}

func (a *Adaptor) GetEventMapTable() gocassa.MapTable {
	return a.Session.MapTable("events", "sg_event_id", &Event{})
}

func (a *Adaptor) GetTimeSeriesEventTable() gocassa.TimeSeriesTable {
	return a.Session.TimeSeriesTable("events", "timestamp", "sg_event_id", 24*time.Hour, &Event{})
}
