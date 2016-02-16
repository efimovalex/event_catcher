package cmd

import (
	"fmt"
	"log"

	"github.com/efimovalex/EventKitAPI/consumerapi"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Runs consumer http server",
}

func init() {
	serverCmd.Run = server
}

func server(cmd *cobra.Command, args []string) {
	fmt.Println("[+] Running Consumer API main on: ")

	var config consumerapi.Config

	if err := envconfig.Process("CONSUMER", &config); err != nil {
		log.Fatal("Error occurred during startup:", err.Error())
	}
	log.Printf("[+]Started API on: %s:%d, with MaxWokers: %d and MaxJobs: %d\n", config.Interface, config.Port, config.MaxWorker, config.MaxJobQueue)

	consumerAPIService := consumerapi.NewService(&config)

	if err := consumerAPIService.Start(); err != nil {
		log.Fatal("Error occurred during startup:", err.Error())
	}
}
