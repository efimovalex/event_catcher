package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{}

func init() {
	RootCmd.AddCommand(consumerCmd)
	RootCmd.AddCommand(restAPICmd)
	RootCmd.AddCommand(cassandraCmd)
}
