package populate_db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/configuration"
	"github.com/AndreyAD1/helsinki-guide/internal/populator"
	"github.com/caarlos0/env/v9"
	"github.com/spf13/cobra"
)

var (
	dbURL       string
	sheetName   string
	PopulateCmd = &cobra.Command{
		Use:   "populate <finnish-file> <english-file> <russian-file>",
		Short: "Populate a database",
		Long:  "This command transfers data from xlsx files to a database",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(args[0], args[1], args[2])
		},
	}
)

func init() {
	PopulateCmd.Flags().StringVarP(
		&dbURL,
		"dburl",
		"d",
		"",
		"A database URL. You can also use an environment variable 'DatabaseURL'.",
	)
	PopulateCmd.Flags().StringVarP(
		&sheetName,
		"sheet",
		"s",
		"",
		"An xlsx sheet name to translate (required)",
	)
	PopulateCmd.MarkFlagRequired("sheet")
}

func run(finFile, enFilename, ruFilename string) error {
	if dbURL != "" {
		os.Setenv("DatabaseURL", dbURL)
	}
	config := configuration.PopulatorConfig{}
	err := env.Parse(&config)
	if err != nil {
		return fmt.Errorf("a configuration error: %w", err)
	}
	ctx := context.Background()
	populator, err := populator.NewPopulator(ctx, config)
	if err != nil {
		return err
	}
	err = populator.Run(ctx, sheetName, finFile, enFilename, ruFilename)
	if err != nil {
		log.Printf("unexpected error: %v", err)
		os.Exit(1)
	}
	return nil
}
