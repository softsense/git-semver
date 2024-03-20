package main

import (
	"fmt"
	"log"
	"os"

	"github.com/softsense/git-semver/pkg/git"
	"github.com/softsense/git-semver/pkg/semver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-semver",
	Short: "A tool for bumping semantic versions based on git tags.",
	Long:  `A tool for bumping semantic versions based on git tags.`,
	Run: func(cmd *cobra.Command, args []string) {
		var below *semver.Version
		if viper.GetString("below") != "" {
			v, err := semver.Parse(viper.GetString("below"))
			if err != nil {
				log.Fatal(err)
			}
			below = &v
		}
		g, err := git.Open(viper.GetString("repo"), git.Config{
			Prefix:    viper.GetString("prefix"),
			Below:     below,
			IncludeRC: viper.GetBool("rc"),
		})
		if err != nil {
			log.Fatal(err)
		}

		n, err := g.Increment(viper.GetBool("major"), viper.GetBool("minor"), viper.GetBool("patch"), viper.GetBool("snapshot"), viper.GetBool("rc"))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(n.String())
	},
}

func init() {
	rootCmd.PersistentFlags().String("repo", "./", "path to git repository")
	if err := viper.BindPFlag("repo", rootCmd.PersistentFlags().Lookup("repo")); err != nil {
		log.Fatal(err)
	}

	rootCmd.Flags().Bool("major", false, "bump major version")
	if err := viper.BindPFlag("major", rootCmd.Flags().Lookup("major")); err != nil {
		log.Fatal(err)
	}

	rootCmd.Flags().Bool("minor", false, "bump minor version")
	if err := viper.BindPFlag("minor", rootCmd.Flags().Lookup("minor")); err != nil {
		log.Fatal(err)
	}

	rootCmd.Flags().Bool("patch", true, "bump patch version")
	if err := viper.BindPFlag("patch", rootCmd.Flags().Lookup("patch")); err != nil {
		log.Fatal(err)
	}

	rootCmd.PersistentFlags().Bool("rc", false, "bump rc version. will bump other version if an rc does not already exist.")
	if err := viper.BindPFlag("rc", rootCmd.Flags().Lookup("rc")); err != nil {
		log.Fatal(err)
	}

	rootCmd.Flags().Bool("snapshot", false, "set snapshot version")
	if err := viper.BindPFlag("snapshot", rootCmd.Flags().Lookup("snapshot")); err != nil {
		log.Fatal(err)
	}

	rootCmd.PersistentFlags().String("prefix", "", "use a prefix")
	if err := viper.BindPFlag("prefix", rootCmd.PersistentFlags().Lookup("prefix")); err != nil {
		log.Fatal(err)
	}

	rootCmd.PersistentFlags().String("below", "", "only look at tags below version")
	if err := viper.BindPFlag("below", rootCmd.PersistentFlags().Lookup("below")); err != nil {
		log.Fatal(err)
	}
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
