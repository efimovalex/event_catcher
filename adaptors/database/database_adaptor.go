package database

import (
	"github.com/gocql/gocql"
)

// Adaptor for DB
type Adaptor struct {
	Session gocql.Session
}

// New db adaptor
func New(ips []string) {
	cluster := gocql.NewCluster(ips...)
	cluster.Keyspace = "event"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()
}
