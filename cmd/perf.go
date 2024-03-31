package cmd

import (
	"blue/perf"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(perfCmd)
	perfCmd.PersistentFlags().Int32Var(&perf.TcpConnNum, "tcp_conn_num", 90000, "tcp 连接的数量，默认30000")
}

var perfCmd = &cobra.Command{
	Use: "perf",
	Run: perfHandle,
}

func perfHandle(cmd *cobra.Command, args []string) {
	perf.RunMain()
}
