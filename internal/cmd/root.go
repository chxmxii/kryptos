package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "kryptos",
	Short:         "Kryptos is a simple key-value store that encrypts values before storing them in Redis.",
	SilenceUsage:  true,
	SilenceErrors: true,
	Version:       "0.0.1",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help() //nolint: errcheck
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
