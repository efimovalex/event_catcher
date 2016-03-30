package database

import (
	"github.com/hailocab/gocassa"
)

// Adaptor for DB
type Adaptor struct {
	Session gocassa.KeySpace
}

type Interface interface {
}

// New db adaptor
func NewAdaptor(ips []string) *Adaptor {
	session, err := gocassa.ConnectToKeySpace("events", ips, "", "")
	if err != nil {
		panic(err)
	}

	return &Adaptor{
		Session: session,
	}
}
