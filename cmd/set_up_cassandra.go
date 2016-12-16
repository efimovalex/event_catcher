package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/efimovalex/EventKitAPI/adaptors/database"
	"github.com/efimovalex/EventKitAPI/consumerapi"
	"github.com/gocql/gocql"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
)

var cassandraCmd = &cobra.Command{
	Use:   "set_up_cassandra",
	Short: "Runs cassandra migration",
}

func init() {
	cassandraCmd.Run = setUp
}

func setUp(cmd *cobra.Command, args []string) {
	fmt.Println("[+] Running Cassandra set up ")

	var config consumerapi.Config

	if err := envconfig.Process("CONSUMER", &config); err != nil {
		log.Fatal("Error occurred during startup:", err.Error())
	}

	cluster := gocql.NewCluster(strings.Split(config.CassandraInterfaces, ",")...)
	cluster.Timeout = 1 * time.Minute
	cluster.ProtoVersion = 3

	session, err := cluster.CreateSession()
	if err != nil {
		log.Print(err.Error())

		return
	}

	qd := session.Query("DROP KEYSPACE IF EXISTS events;", "")
	if err := qd.Exec(); err != nil {
		log.Print(err.Error())
	}

	q := session.Query("CREATE KEYSPACE IF NOT EXISTS events WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1};", "")
	if err := q.Exec(); err != nil {
		log.Print(err.Error())
	}

	DBAdaptor := database.NewAdaptor(strings.Split(config.CassandraInterfaces, ","), config.CassandraUser, config.CassandraPassword)

	for _, indexField := range database.IndexFields {
		eventsTable := DBAdaptor.GetEventMultimapTable(indexField)

		if err := eventsTable.CreateIfNotExist(); err != nil {
			log.Print(err.Error())
		}
	}

	eventTable := DBAdaptor.GetTimeSeriesEventTable()
	if err := eventTable.CreateIfNotExist(); err != nil {
		log.Print(err.Error())
	}

	eventMapTable := DBAdaptor.GetEventMapTable()
	if err := eventMapTable.CreateIfNotExist(); err != nil {
		log.Print(err.Error())
	}
}
