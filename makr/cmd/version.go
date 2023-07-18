package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var (
	version string
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Print the version number of Makr",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("Makr - v%s", viper.GetString("version")))
	},
}
