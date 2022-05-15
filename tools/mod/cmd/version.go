package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	cliVersion = "v1.0.6"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print CLI version",
	Long:  `Print CLI version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Current version: ", cliVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
