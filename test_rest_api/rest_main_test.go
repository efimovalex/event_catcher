// +build rest_acceptance

/*
  Each set of acceptance tests should be in their own file
  This file is a good place to put common test files, such as
  functions and structs that are re-used.
  There should be no actual tests here.
*/

package rest_acceptance

import (
	"github.com/efimovalex/EventKitAPI/restapi"
	"github.com/kelseyhightower/envconfig"
	"log"
)

var acceptanceConfig = restapi.Config{}

func init() {
	if err := envconfig.Process("REST", &acceptanceConfig); err != nil {
		log.Fatal("Error occurred during startup:", err.Error())
	}
}
