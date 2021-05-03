package main

import (
	"fmt"
	"log"

	"github.com/softsense/git-semver/pkg/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(historyCmd)
	historyCmd.PersistentFlags().String("msg-prefix", "", "use a prefix for the messages")
	if err := viper.BindPFlag("msg-prefix", historyCmd.PersistentFlags().Lookup("msg-prefix")); err != nil {
		log.Fatal(err)
	}
}

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Print history since last tag.",
	Run: func(cmd *cobra.Command, args []string) {
		g, err := git.Open(viper.GetString("repo"), git.Config{
			Prefix: viper.GetString("prefix"),
		})
		if err != nil {
			log.Fatal(err)
		}

		history, err := g.History(viper.GetString("msg-prefix"))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(history)
	},
}
