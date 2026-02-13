package cmd

import (
	"os"

	"github.com/niedch/mux-session/internal/conf"
	"github.com/niedch/mux-session/internal/dataproviders"
	"github.com/niedch/mux-session/internal/fzf"
	"github.com/niedch/mux-session/internal/logger"
	"github.com/niedch/mux-session/internal/orchestrator"
	"github.com/niedch/mux-session/internal/tmux"
	"github.com/spf13/cobra"
)

var (
	configFile string
	socket     string
	verbose    bool
)

var rootCmd = &cobra.Command{
	Use:   "mux-session",
	Short: "Interactive tmux session manager",
	Long: `mux-session is an interactive tmux session manager that allows you to quickly
navigate to project directories and create or switch to tmux sessions.

It provides an fzf-powered interface for selecting directories from configured
search paths and automatically creates or attaches to tmux sessions with the
directory name as the session name.`,
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			logger.SetEnabled(true)
		}
		logger.Printf("Loading configuration from: %s\n", configFile)
		config, err := conf.Load(configFile)
		if err != nil {
			logger.Fatalf("Failed to load config: %v\n", err)
		}
		logger.Printf("Configuration loaded successfully\n")

		logger.Printf("Initializing tmux wrapper (socket: %s)\n", socket)
		tmux, err := tmux.NewTmux(socket)
		if err != nil {
			logger.Fatalf("Failed to initialize tmux: %v\n", err)
		}

		multiService := orchestrator.New(tmux)

		directoryProvider := dataproviders.NewDirectoryProvider(config.SearchPaths)
		tmuxProvider := dataproviders.NewTmuxProvider(tmux)

		logger.Printf("Starting interactive session selector\n")
		selected, err := fzf.Run(dataproviders.NewDeduplicatorProvider(directoryProvider, tmuxProvider).WithMarkDuplicates(true), config)

		if err != nil {
			logger.Fatalf("Session selector failed: %v\n", err)
		}

		if selected == nil {
			logger.Printf("No selection made, exiting\n")
			return
		}
		logger.Printf("Selected session: id=%s, display=%s\n", selected.Id, selected.Display)

		projectConfig := config.GetProjectConfig(selected.Id)
		logger.Printf("Switching to session: %s\n", selected.Id)
		ok, err := multiService.SwitchSession(selected)
		if err != nil {
			logger.Fatalf("Failed to switch session: %v\n", err)
		}

		if ok {
			logger.Printf("Switched to existing session: %s\n", selected.Id)
			return
		}

		logger.Printf("Creating new session: %s\n", selected.Id)
		if err := multiService.CreateSession(selected, projectConfig); err != nil {
			logger.Fatalf("Failed to create session: %v\n", err)
		}
		logger.Printf("Session created successfully: %s\n", selected.Id)
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

	rootCmd.PersistentFlags().StringVarP(&configFile, "file", "f", "", "Path to config file (default is XDG_CONFIG/mux-session/config.toml)")
	rootCmd.PersistentFlags().StringVarP(&socket, "socket", "L", "", "tmux socket name for targeting a specific server")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Toggle flag for testing")
}
