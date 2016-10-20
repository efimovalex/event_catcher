package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/efimovalex/EventKitAPI/adaptors/database"
	"github.com/efimovalex/EventKitAPI/consumerapi"
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

	DBAdaptor := database.NewAdaptor(strings.Split(config.CassandraInterfaces, ","), config.CassandraUser, config.CassandraPassword)

	for _, indexField := range database.IndexFields {
		eventsTable := DBAdaptor.GetEventMultimapTable(indexField)

		if err := eventsTable.Create(); err != nil {
			log.Print(err.Error())
		}
	}

	eventTable := DBAdaptor.GetTimeSeriesEventTable()
	if err := eventTable.Create(); err != nil {
		log.Print(err.Error())
	}

	eventMapTable := DBAdaptor.GetEventMapTable()
	if err := eventMapTable.Create(); err != nil {
		log.Print(err.Error())
	}
}
