package cmd

import (
	"log"
	"os"
	"path/filepath"
	"slices"

	"github.com/niedch/mux-session/internal/conf"
	"github.com/niedch/mux-session/internal/fzf"
	"github.com/niedch/mux-session/internal/tmux"
	"github.com/spf13/cobra"
)

var configFile string

var rootCmd = &cobra.Command{
	Use:   "mux-session",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := conf.Load(configFile)
		if err != nil {
			log.Fatal(err)
			return
		}

		tmux, err := tmux.NewTmux()
		if err != nil {
			log.Fatal(err)
			return
		}

		provider := fzf.NewDirectoryProvider(config.SearchPaths)
		dir, err := fzf.StartFzf(provider)
		if err != nil {
			log.Fatal(err)
			return
		}

		dir_name := filepath.Base(*dir)
		projectConfig := config.GetProjectConfig(dir_name)

		sessions, err := tmux.ListSessions()
		if err != nil {
			log.Fatal(err)
			return
		}

		if slices.Contains(sessions, dir_name) {
			if err := tmux.SwitchSession(dir_name); err != nil {
				log.Fatal(err)
				return
			}

			return
		}

		if err := tmux.CreateSession(*dir, projectConfig); err != nil {
			log.Fatal(err)
			return
		}
	},
}

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

	rootCmd.Flags().StringVarP(&configFile, "file", "f", "", "Path to config file (default is XDG_CONFIG/mux-session/config.toml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
