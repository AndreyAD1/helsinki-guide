package translate

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/AndreyAD1/helsinki-guide/translator"
)


var apiKey string
var TranslateCmd = &cobra.Command{
  Use:   "translate",
  Short: "Translate a building dataset",
  Run: func(cmd *cobra.Command, args []string) {
    run()
  },
}

func init() {
	TranslateCmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "A key for a translation API (required)")
	TranslateCmd.MarkFlagRequired("api-key")
}

func run() {
	translator.Run(context.Background(), apiKey)
}