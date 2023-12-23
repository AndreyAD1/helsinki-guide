package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	databaseURL  string
	converterURL string
	limit        int
	offset       int
	RootCmd      = &cobra.Command{
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
		"A database URL",
	)
	RootCmd.Flags().StringVarP(
		&converterURL,
		"converter",
		"c",
		"https://epsg.io",
		"A base URL for a converter service",
	)
	RootCmd.Flags().IntVarP(
		&limit,
		"number",
		"n",
		100,
		"A number of buildings to convert",
	)
	RootCmd.Flags().IntVarP(
		&offset,
		"offset",
		"o",
		0,
		"A building offset for the table 'buildings' (default 0)",
	)
	RootCmd.MarkFlagRequired("dburl")
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
