package cmd

import (
	"github.com/origine-run/figs/pkg/validator"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen", "g"},
	Short:   "Generate the required file type",
	Long:    "...",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		data, _ := cmd.Flags().GetString("data")
		template, _ := cmd.Flags().GetString("template")
		format, _ := cmd.Flags().GetString("format")
		storage, _ := cmd.Flags().GetString("storage")

		//validate
		val := validator.NewJson()
		status := val.Validate(source, scheme)

		//map (for pdf only)

		//convert

		//save file to disk
	},
}

func init() {
	RootCmd.AddCommand(generateCmd)
}
