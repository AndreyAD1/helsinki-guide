package translate

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/AndreyAD1/helsinki-guide/infrastructure"
	ts "github.com/AndreyAD1/helsinki-guide/translator"
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
	TranslateCmd.Flags().StringVarP(
		&apiKey, 
		"api-key", 
		"k",
		"",
		"A key for a translation API (required)",
	)
	TranslateCmd.MarkFlagRequired("api-key")
}

func run() {
	translator := ts.NewTranslator(infrastructure.NewGoogleClient(apiKey))
	if err := translator.Run(context.Background()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}