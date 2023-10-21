package translate

import (
	"context"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/AndreyAD1/helsinki-guide/infrastructure"
	ts "github.com/AndreyAD1/helsinki-guide/translator"
)

var (
	apiKey       string
	sheetName    string
	TranslateCmd = &cobra.Command{
		Use:   "translate <source> <target>",
		Short: "Translate a building dataset",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			run(args)
		},
	}
)

func init() {
	TranslateCmd.Flags().StringVarP(
		&apiKey,
		"api-key",
		"k",
		"",
		"A key for a translation API (required)",
	)
	TranslateCmd.Flags().StringVarP(
		&sheetName,
		"sheet",
		"s",
		"",
		"A name of Excel sheet to translate (required)",
	)
	TranslateCmd.MarkFlagRequired("api-key")
	TranslateCmd.MarkFlagRequired("sheet")
}

func run(args []string) {
	translator := ts.NewTranslator(infrastructure.NewGoogleClient(apiKey))
	source, target := args[0], args[1]
	if err := translator.Run(context.Background(), source, sheetName, target); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
