package cmd

import (
	"blue/gateway"
	"github.com/spf13/cobra"
)

func init() {
	// Add the gateway command
	rootCmd.AddCommand(gatewayCmd)
}

var gatewayCmd = &cobra.Command{
	Use: "gateway",
	Run: GatewayHandler,
}

func GatewayHandler(cmd *cobra.Command, args []string) {
	// Handle the gateway command
	gateway.RunMain(ConfigPath)
}
