package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:                "version",
	Example:            "kryptos version",
	Short:              "Print the version number of Kryptos",
	Long:               `All software has versions. This is Kryptos's`,
	Run:                Version,
	Args:               cobra.NoArgs,
	SilenceUsage:       true,
	SilenceErrors:      true,
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(versionCmd)

}

func Version(cmd *cobra.Command, args []string) {
	fmt.Printf("Kryptos version: %s\n", rootCmd.Version)
}
