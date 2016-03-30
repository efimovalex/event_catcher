package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/efimovalex/EventKitAPI/adaptors/database"
	"github.com/efimovalex/EventKitAPI/consumerapi"
	"github.com/hailocab/gocassa"
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

	keySpace, err := gocassa.ConnectToKeySpace("events", strings.Split(config.CassandraInterfaces, ","), "", "")
	if err != nil {
		panic(err)
	}

	bounceTable := keySpace.Table("bounce_events", &database.BounceEvent{}, gocassa.Keys{
		PartitionKeys: []string{"sg_event_id"},
	})

	if err := bounceTable.Recreate(); err != nil {
		log.Print(err.Error())
	}

	deferredTable := keySpace.Table("deferred_events", &database.DeferredEvent{}, gocassa.Keys{
		PartitionKeys: []string{"sg_event_id"},
	})

	if err := deferredTable.Recreate(); err != nil {
		log.Print(err.Error())
	}

	deliveredTable := keySpace.Table("delivered_events", &database.DeliveredEvent{}, gocassa.Keys{
		PartitionKeys: []string{"sg_event_id"},
	})

	if err := deliveredTable.Recreate(); err != nil {
		log.Print(err.Error())
	}

	clickTable := keySpace.Table("click_events", &database.ClickEvent{}, gocassa.Keys{
		PartitionKeys: []string{"sg_event_id"},
	})

	if err := clickTable.Recreate(); err != nil {
		log.Print(err.Error())
	}

	dropTable := keySpace.Table("dropped_events", &database.DroppedEvent{}, gocassa.Keys{
		PartitionKeys: []string{"sg_event_id"},
	})

	if err := dropTable.Recreate(); err != nil {
		log.Print(err.Error())
	}

	groupUnsubscribeTable := keySpace.Table("group_unsubscribe_events", &database.GroupUnsubscribeEvent{}, gocassa.Keys{
		PartitionKeys: []string{"sg_event_id"},
	})

	if err := groupUnsubscribeTable.Recreate(); err != nil {
		log.Print(err.Error())
	}

	groupResubscribeTable := keySpace.Table("group_resubscribe_events", &database.GroupResubscribeEvent{}, gocassa.Keys{
		PartitionKeys: []string{"sg_event_id"},
	})

	if err := groupResubscribeTable.Recreate(); err != nil {
		log.Print(err.Error())
	}

	openTable := keySpace.Table("open_events", &database.OpenEvent{}, gocassa.Keys{
		PartitionKeys: []string{"sg_event_id"},
	})

	if err := openTable.Recreate(); err != nil {
		log.Print(err.Error())
	}

	processedTable := keySpace.Table("processed_events", &database.ProcessedEvent{}, gocassa.Keys{
		PartitionKeys: []string{"sg_event_id"},
	})

	if err := processedTable.Recreate(); err != nil {
		log.Print(err.Error())
	}

	spamTable := keySpace.Table("spam_events", &database.SpamEvent{}, gocassa.Keys{
		PartitionKeys: []string{"sg_event_id"},
	})

	if err := spamTable.Recreate(); err != nil {
		log.Print(err.Error())
	}

	unsubscribeTable := keySpace.Table("unsubscribe_events", &database.UnsubscribeEvent{}, gocassa.Keys{
		PartitionKeys: []string{"sg_event_id"},
	})

	if err := unsubscribeTable.Recreate(); err != nil {
		log.Print(err.Error())
	}
}
