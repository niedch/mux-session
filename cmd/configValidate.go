package cmd

import (
	"log"

	"github.com/niedch/mux-session/internal/conf"
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
		config, err := conf.Load(configFile)
		if err != nil {
			log.Fatalf("Error loading config: %v\n", err)
		}

		config.PrettyPrint()
	},
}

func init() {
	configValidateCmd.Flags().StringVarP(&configFile, "file", "f", "", "Path to config file (default is XDG_CONFIG/mux-session/config.toml)")
	rootCmd.AddCommand(configValidateCmd)
}
