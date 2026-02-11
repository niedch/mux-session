package cmd

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

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
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get user home directory: %v\n", err)
			os.Exit(1)
		}
		muxSessionConfigDir := filepath.Join(homeDir, ".config", "mux-session")

		if err := os.MkdirAll(muxSessionConfigDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create config directory '%s': %v\n", muxSessionConfigDir, err)
			os.Exit(1)
		}

		configFilePath := filepath.Join(muxSessionConfigDir, "config.toml")

		if _, err := os.Stat(configFilePath); !os.IsNotExist(err) {
			fmt.Printf("Config file already exists at %s. Aborting.\n", configFilePath)
			return
		}

		fmt.Print("Enter the path to the folder which holds your projects (comma-separated): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		searchPaths := strings.TrimSpace(input)

		tmpl, err := template.New("config").Parse(string(defaultConfigTemplate))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse config template: %v\n", err)
			os.Exit(1)
		}

		file, err := os.Create(configFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create config file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		data := struct {
			SearchPaths string
		}{
			SearchPaths: fmt.Sprintf("\"%s\"", strings.Join(strings.Split(searchPaths, ","), "\", \"")),
		}

		if err := tmpl.Execute(file, data); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write config file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully created config file at %s 🚀\n", configFilePath)
		fmt.Println("You can now run mux-session ✨")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
