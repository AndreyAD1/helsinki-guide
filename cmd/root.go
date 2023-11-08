package cmd

import (
	"github.com/spf13/cobra"

	"github.com/AndreyAD1/helsinki-guide/cmd/bot"
	"github.com/AndreyAD1/helsinki-guide/cmd/populate_db"
	"github.com/AndreyAD1/helsinki-guide/cmd/translate"
)

var debug *bool
var rootCmd = &cobra.Command{
	Use:   "helsinki-guide",
	Short: "The 'HelsinkiGuide' telegram bot provides information about notable Helsinki buildings.",
	Long: `This bot is designed to provide information about notable buildings in Helsinki.
Data source: History of buildings in Helsinki. The maintainer of the dataset is 
Helsingin kulttuurin ja vapaa-ajan toimiala / Kaupunginmuseo and 
the original author is Tmi Hilla Tarjanne.
https://hri.fi/data/en_GB/dataset/helsinkilaisten-rakennusten-historiatietoja`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(translate.TranslateCmd)
	rootCmd.AddCommand(bot.BotCmd)
	rootCmd.AddCommand(populate_db.PopulateCmd)
	debug = rootCmd.Flags().BoolP("debug", "d", false, "Run in a debug mode")
}
