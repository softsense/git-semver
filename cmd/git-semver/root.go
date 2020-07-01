package main

import (
	"fmt"
	"log"
	"os"

	"github.com/softsense/git-semver/pkg/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-semver",
	Short: "A tool for bumping semantic versions based on git tags.",
	Long:  `A tool for bumping semantic versions based on git tags.`,
	Run: func(cmd *cobra.Command, args []string) {
		g, err := git.Open(viper.GetString("repo"), git.Config{
			Prefix: viper.GetString("prefix"),
		})
		if err != nil {
			log.Fatal(err)
		}

		n, err := g.Increment(viper.GetBool("major"), viper.GetBool("minor"), viper.GetBool("patch"), viper.GetBool("snapshot"))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(n.String())
	},
}

func init() {
	rootCmd.Flags().String("repo", "./", "path to git repository")
	if err := viper.BindPFlag("repo", rootCmd.Flags().Lookup("repo")); err != nil {
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

	rootCmd.Flags().Bool("snapshot", false, "set snapshot version")
	if err := viper.BindPFlag("snapshot", rootCmd.Flags().Lookup("snapshot")); err != nil {
		log.Fatal(err)
	}

	rootCmd.Flags().String("prefix", "", "use a prefix")
	if err := viper.BindPFlag("prefix", rootCmd.Flags().Lookup("prefix")); err != nil {
		log.Fatal(err)
	}
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
