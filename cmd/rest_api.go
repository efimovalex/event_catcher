package cmd

import (
	"log"

	"github.com/efimovalex/EventKitAPI/restapi"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
)

var restAPICmd = &cobra.Command{
	Use:   "restapi",
	Short: "Runs REST api http server",
}

func init() {
	restAPICmd.Run = restapiServer
}

func restapiServer(cmd *cobra.Command, args []string) {
	var config restapi.Config

	if err := envconfig.Process("REST", &config); err != nil {
		log.Fatal("Error occurred during startup:", err.Error())
	}
	log.Print(`.          +         .         .                 .  .
       .                 .                   .               .
               .    ,,o         .                  __.o+.
     .            od8^                  .      oo888888P^b           .
        .       ,".o'      .     .             'b^'""'b -'b   .
              ,'.'o'             .   .          t. = -'b -'t.    .
             ; d o' .        ___          _.--.. 8  -  'b  ='b
         .  dooo8<       .o:':__;o.     ,;;o88%%8bb - = 'b  ='b.    .
     .     |^88^88=. .,x88/::/ | \\';;;;;;d%%%%%88%88888/%x88888
           :-88=88%%L8'%'|::|_>-<_||%;;%;8%%=;:::=%8;;\%%%%\8888
       .   |=88 88%%|HHHH|::| >-< |||;%;;8%%=;:::=%8;;;%%%%+|]88        .
           | 88-88%%LL.%.%b::Y_|_Y/%|;;;;'%8%%oo88%:o%.;;;;+|]88  .
           Yx88o88^^'"'^^%8boooood..-\H_Hd%P%%88%P^%%^'\;;;/%%88
          . '"\^\          ~"""""'      d%P """^" ;   = '+' - P
    .        '.'.b   .                :<%%>  .   :  -   d' - P      . .
               .'.b     .        .    '788      ,'-  = d' =.'
        .       ''.b.                           :..-  :'  P
             .   'q.>b         .               '^^^:::::,'       .
     LS            ""^^               .                     .
   .                                           .               .       .
     .         .          .                 .        +         .
                                  .`)
	log.Printf("[+] Started REST API on: %s:%d", config.Interface, config.Port)

	RESTAPIService := restapi.NewService(&config)

	if err := RESTAPIService.Start(&config); err != nil {
		log.Fatal("Error occurred during startup:", err.Error())
	}
}
