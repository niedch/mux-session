package cmd

import (
	"github.com/niedch/mux-session/internal/conf"
	"github.com/niedch/mux-session/internal/logger"
	"github.com/spf13/cobra"
)

// configValidateCmd represents the configValidate command
var configValidateCmd = &cobra.Command{
	Use:   "config-validate",
	Short: "Validate and display configuration",
	Long: `Loads the mux-session configuration file and displays it in a formatted
JSON structure. This command validates that your configuration is properly
parsed and shows the current settings including search paths and project
configurations.`,
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			logger.SetEnabled(true)
		}
		logger.Printf("Loading configuration from: %s\n", configFile)
		config, err := conf.Load(configFile)
		if err != nil {
			logger.Fatalf("Error loading config: %v\n", err)
		}
		logger.Printf("Configuration loaded successfully\n")

		config.PrettyPrint()
	},
}

func init() {
	rootCmd.AddCommand(configValidateCmd)
}
