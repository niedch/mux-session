package cmd

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/niedch/mux-session/internal/logger"
	"github.com/spf13/cobra"
)

//go:embed config.toml.template
var defaultConfigTemplate []byte

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates a default config file in ~/.config/mux-session/",
	Long: `This command creates a default configuration file for mux-session.
It will create a 'mux-session' directory inside your default config directory
(e.g., ~/.config on Linux) and place a 'config.toml' file there.`,
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			logger.SetEnabled(true)
		}
		logger.Printf("Initializing mux-session configuration\n")
		homeDir, err := os.UserHomeDir()
		if err != nil {
			logger.Fatalf("Failed to get user home directory: %v\n", err)
		}
		muxSessionConfigDir := filepath.Join(homeDir, ".config", "mux-session")
		logger.Printf("Config directory: %s\n", muxSessionConfigDir)

		if err := os.MkdirAll(muxSessionConfigDir, 0755); err != nil {
			logger.Fatalf("Failed to create config directory '%s': %v\n", muxSessionConfigDir, err)
		}
		logger.Printf("Config directory created/verified\n")

		configFilePath := filepath.Join(muxSessionConfigDir, "config.toml")

		if _, err := os.Stat(configFilePath); !os.IsNotExist(err) {
			fmt.Printf("Config file already exists at %s. Aborting.\n", configFilePath)
			return
		}
		logger.Printf("Config file does not exist, proceeding with creation\n")

		fmt.Print("Enter the path to the folder which holds your projects (comma-separated): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		searchPaths := strings.TrimSpace(input)
		logger.Printf("User provided search paths: %s\n", searchPaths)

		tmpl, err := template.New("config").Parse(string(defaultConfigTemplate))
		if err != nil {
			logger.Fatalf("Failed to parse config template: %v\n", err)
		}

		file, err := os.Create(configFilePath)
		if err != nil {
			logger.Fatalf("Failed to create config file: %v\n", err)
		}
		defer file.Close()
		logger.Printf("Config file created: %s\n", configFilePath)

		data := struct {
			SearchPaths string
		}{
			SearchPaths: fmt.Sprintf("\"%s\"", strings.Join(strings.Split(searchPaths, ","), "\", \"")),
		}

		if err := tmpl.Execute(file, data); err != nil {
			logger.Fatalf("Failed to write config file: %v\n", err)
		}
		logger.Printf("Config template executed successfully\n")

		fmt.Printf("Successfully created config file at %s 🚀\n", configFilePath)
		logger.Printf("Initialization complete\n")
		fmt.Println("You can now run mux-session ✨")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
