package cmd

import (
	"blue/client"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(clientCmd)
}

var clientCmd = &cobra.Command{
	Use: "client",
	Run: cli,
}

func cli(cmd *cobra.Command, args []string) {
	client.RunMain()
}
