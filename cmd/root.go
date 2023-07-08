/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/detect"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "shush",
	Short: "A simple application that redacts secrets from strings",
	Long: `Shush is an application that takes in a block of text and redacts secrets from that text.

	Great when needing to share config files for debugging or similar reasons.`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetConfigType("toml")
		err := viper.ReadConfig(strings.NewReader(config.DefaultConfig))
		if err != nil {
			panic(err)
		}

		defaultViperConfig := config.ViperConfig{}
		if err := viper.Unmarshal(&defaultViperConfig); err != nil {
			panic(err)
		}

		cfg, err := defaultViperConfig.Translate()
		if err != nil {
			panic(err)
		}

		detector := detect.NewDetector(cfg)
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}

		s := string(b)
		findings := detector.DetectString(s)

		if err != nil {
			panic(err)
		}

		for _, finding := range findings {
			s = strings.ReplaceAll(s, finding.Secret, "[[REDACTED]]")
		}

		_, err = os.Stdout.WriteString(s)

		if err != nil {
			panic(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.shush.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
