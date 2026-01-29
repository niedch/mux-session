package cmd

import (
	"fmt"

	"github.com/niedch/mux-session/internal/conf"
	"github.com/spf13/cobra"
)

// configValidateCmd represents the configValidate command
var configValidateCmd = &cobra.Command{
	Use:   "config-validate",
	Short: "Validate and pretty print the configuration",
	Long: `Loads the configuration from config.toml and displays it in a formatted
JSON structure. This command is useful for validating that your configuration
is properly loaded and structured.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := conf.Load(configFile)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}
		config.PrettyPrint()
	},
}

func init() {
	configValidateCmd.Flags().StringVarP(&configFile, "file", "f", "", "Path to config file (default is XDG_CONFIG/mux-session/config.toml)")
	rootCmd.AddCommand(configValidateCmd)
}
