package cmd

import (
	"io"
	"os"
	"regexp"
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

		if showPasswords, _ := cmd.Flags().GetBool("showPasswords"); !showPasswords {
			cfg.Rules["shush:password"] = config.Rule{
				Description: "Password",
				RuleID:      "Generic password",
				Regex:       regexp.MustCompile("['\"]?(password|pass|pwd|passwd|passphrase)['\"]?[[:space:]]?[:-]?[[:space:]]?['\"]?([^'\"\r\n]*)['\"]?"),
				SecretGroup: 2,
				Keywords:    []string{"password", "pass", "pwd", "passwd", "passphrase", "passphrase"},
			}
		}

		if showSecrets, _ := cmd.Flags().GetBool("showSecrets"); !showSecrets {
			cfg.Rules["shush:secret"] = config.Rule{
				Description: "Secret",
				RuleID:      "Generic secrets",
				Regex:       regexp.MustCompile("['\"]?(secret)['\"]?[[:space:]]?[:-]?[[:space:]]?['\"]?([^'\"\r\n]*)['\"]?"),
				SecretGroup: 2,
				Keywords:    []string{"secret"},
			}
		}

		detector := detect.NewDetector(cfg)

		var b []byte
		if len(args) == 0 {
			b, err = io.ReadAll(os.Stdin)
		} else {
			b, err = os.ReadFile(args[0])
		}

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
	rootCmd.Flags().BoolP("showPasswords", "p", false, "Dont redact generic passwords")
	rootCmd.Flags().BoolP("showSecrets", "s", false, "Dont redact generic secrets")
}
