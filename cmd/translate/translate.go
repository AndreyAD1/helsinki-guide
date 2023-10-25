package translate

import (
	"context"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/AndreyAD1/helsinki-guide/internal/infrastructure/clients"
	ts "github.com/AndreyAD1/helsinki-guide/internal/translator"
)

var (
	apiKey       string
	sheetName    string
	TranslateCmd = &cobra.Command{
		Use:   "translate <source> <target>",
		Short: "Translate a building dataset",
		Long: `This command translates a dataset from Finnish into English.
The dataset is located in an xlsx file that can be downloaded from https://hri.fi/data/en_GB/dataset/helsinkilaisten-rakennusten-historiatietoja`,
		Args: cobra.ExactArgs(2),
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
	translator := ts.NewTranslator(clients.NewGoogleClient(apiKey))
	source, target := args[0], args[1]
	if err := translator.Run(context.Background(), source, sheetName, target); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
