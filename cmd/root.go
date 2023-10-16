package cmd

import (
	"github.com/spf13/cobra"

	"github.com/AndreyAD1/helsinki-guide/cmd/translate"
)

var debug *bool
var rootCmd = &cobra.Command{
	Use:   "helsinki-guide",
	Short: "A telegram bot providing information about notable Helsinki buildings.",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(translate.TranslateCmd)
	debug = rootCmd.Flags().BoolP("debug", "d", false, "Run in a debug mode")
}