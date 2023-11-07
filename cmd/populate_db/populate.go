package populate_db

import (
	"context"
	"fmt"
	"os"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/configuration"
	"github.com/AndreyAD1/helsinki-guide/internal/populator"
	"github.com/caarlos0/env/v9"
	"github.com/spf13/cobra"
)

var (
	dbURL string
	PopulateCmd   = &cobra.Command{
		Use:   "populate <path-to-a-source-xlsx-file>",
		Short: "Populate a database",
		Long: "This command transfers data from an xlsx file to a database",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(args[0])
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
}

func run(sourceFilename string) error {
	os.Setenv("PopulatorSource", sourceFilename)
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
	populator.Run(ctx)
	return nil
}