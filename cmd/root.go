package cmd

import (
	"github.com/spf13/cobra"

	"github.com/AndreyAD1/helsinki-guide/cmd/bot"
	"github.com/AndreyAD1/helsinki-guide/cmd/global_flags"
	"github.com/AndreyAD1/helsinki-guide/cmd/populate_db"
	"github.com/AndreyAD1/helsinki-guide/cmd/translate"
)

var RootCmd = &cobra.Command{
	Use:   "helsinki-guide",
	Short: "The 'HelsinkiGuide' telegram bot provides information about notable Helsinki buildings.",
	Long: `This bot is designed to provide information about notable buildings in Helsinki.
Data source: History of buildings in Helsinki. The maintainer of the dataset is 
Helsingin kulttuurin ja vapaa-ajan toimiala / Kaupunginmuseo and 
the original author is Tmi Hilla Tarjanne.
https://hri.fi/data/en_GB/dataset/helsinkilaisten-rakennusten-historiatietoja`,
}

func Execute() error {
	return RootCmd.Execute()
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&global_flags.Debug, "debug", "d", false, "Run in a debug mode")
	RootCmd.AddCommand(translate.TranslateCmd)
	RootCmd.AddCommand(bot.BotCmd)
	RootCmd.AddCommand(populate_db.PopulateCmd)
}
