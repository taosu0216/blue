package cmd

import (
	"github.com/spf13/cobra"
)

var (
	ConfigPath string
)

func init() {

	rootCmd.PersistentFlags().StringVarP(&ConfigPath, "config", "c", "./blue.yaml", "config file (default is $HOME/.blue.yaml)")
}

var rootCmd = &cobra.Command{
	Use:   "blue",
	Short: "blue is a CLI IM tool",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
