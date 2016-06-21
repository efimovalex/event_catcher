package cmd

import (
	"log"

	"github.com/efimovalex/EventKitAPI/consumerapi"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
)

var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Runs consumer http server",
}

func init() {
	consumerCmd.Run = consumerServer
}

func consumerServer(cmd *cobra.Command, args []string) {
	var config consumerapi.Config

	if err := envconfig.Process("CONSUMER", &config); err != nil {
		log.Fatal("Error occurred during startup:", err.Error())
	}
	log.Print(`.    .                     .
                  _        .                          .            (
                 (_)        .       .                                     .
  .        ____.--^.
   .      /:  /    |                               +           .         .
         /:  '--=--'   .                                                .
        /: __[\=='-.___          *           .
       /__|\ _~~~~~~   ~~--..__            .             .
       \   \|::::|-----.....___|~--.                                 .
        \ _\_~~~~~-----:|:::______//---...___
    .   [\  \  __  --     \       ~  \_      ~~~===------==-...____
        [============================================================-
        /         __/__   --  /__    --       /____....----''''~~~~      .
  *    /  /   ==           ____....=---='''~~~~ .
 .    /____....--=-''':~~~~                      .                .
      .       ~--~              
                     .                                               .
                          .                      .                      .`)

	log.Printf("[+] Started Consumer API on: %s:%d, with MaxWokers: %d and MaxJobs: %d\n", config.Interface, config.Port, config.MaxWorker, config.MaxJobQueue)

	consumerAPIService := consumerapi.NewService(&config)

	if err := consumerAPIService.Start(); err != nil {
		log.Fatal("Error occurred during startup:", err.Error())
	}
}
