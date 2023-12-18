package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	databaseURL  string
	converterURL string
	RootCmd      = &cobra.Command{
		Use:   "convert",
		Short: "Convert ETRS89/GK25FIN coordinates into WGS84 coordinates.",
		Long: `This script gets projected coordinates from columns 'latitude_etrsgk25' and
'longitude_etrsgk25' of the table 'buildings', converts them into
geographic coordinates, and sets the columns 'latitude_wgs84' and 
'longitude_wgs84' respectively.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func init() {
	RootCmd.Flags().StringVarP(
		&databaseURL,
		"dburl",
		"u",
		"",
		"A database URL.",
	)
	RootCmd.Flags().StringVarP(
		&converterURL,
		"converter",
		"c",
		"https://epsg.io",
		"A base URL for a converter service. Default: https://epsg.io",
	)
	RootCmd.MarkFlagRequired("dburl")
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}