// +build consumer_acceptance

/*
  Each set of acceptance tests should be in their own file
  This file is a good place to put common test files, such as
  functions and structs that are re-used.
  There should be no actual tests here.
*/

package consumer_acceptance

import (
	"github.com/efimovalex/EventKitAPI/consumerapi"
	"github.com/kelseyhightower/envconfig"
	"log"
)

var acceptanceConfig = consumerapi.Config{}

func init() {
	if err := envconfig.Process("CONSUMER", &acceptanceConfig); err != nil {
		log.Fatal("Error occurred during startup:", err.Error())
	}
}
