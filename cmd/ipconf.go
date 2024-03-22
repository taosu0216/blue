package cmd

import (
	"blue/ipconf"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(ipConfCmd)
}

var ipConfCmd = &cobra.Command{
	Use: "ipconf",
	Run: ipconfHandle,
}

func ipconfHandle(cmd *cobra.Command, args []string) {
	ipconf.RunMain(ConfigPath)
}
