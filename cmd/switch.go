package cmd

import (
	"fmt"

	"github.com/niedch/mux-session/internal/conf"
	"github.com/niedch/mux-session/internal/dataproviders"
	"github.com/niedch/mux-session/internal/logger"
	"github.com/niedch/mux-session/internal/orchestrator"
	"github.com/niedch/mux-session/internal/tmux"
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch <id>",
	Short: "Switch to an existing or create a new tmux session",
	Long: `Switch to an existing tmux session or create a new one from the given ID.
The ID can be a session name or directory path.`,
	Args: cobra.ExactArgs(1),
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

		directoryProvider := dataproviders.NewDirectoryProvider(config.SearchPaths)
		tmuxProvider := dataproviders.NewTmuxProvider(tmux)
		composedProvider := dataproviders.NewDeduplicatorProvider(directoryProvider, tmuxProvider).WithMarkDuplicates(true)

		logger.Printf("Getting items from providers\n")
		items, err := composedProvider.GetItems()
		if err != nil {
			logger.Fatalf("Failed to get items: %v\n", err)
		}
		logger.Printf("Found %d items\n", len(items))

		id := args[0]
		logger.Printf("Looking for item with id: %s\n", id)
		item, err := findItem(id, items)
		if err != nil {
			logger.Fatalf("Failed to find item: %v\n", err)
		}
		logger.Printf("Found item: id=%s, display=%s\n", item.Id, item.Display)

		multiService := orchestrator.New(tmux)
		projectConfig := config.GetProjectConfig(item.Id)

		logger.Printf("Switching to session: %s\n", item.Id)
		ok, err := multiService.SwitchSession(item)
		if err != nil {
			logger.Fatalf("Failed to switch session: %v\n", err)
		}

		if ok {
			logger.Printf("Switched to existing session: %s\n", item.Id)
			return
		}

		logger.Printf("Creating new session: %s\n", item.Id)
		if err := multiService.CreateSession(item, projectConfig); err != nil {
			logger.Fatalf("Failed to create session: %v\n", err)
		}
		logger.Printf("Session created successfully: %s\n", item.Id)
	},
}

func findItem(id string, items []dataproviders.Item) (*dataproviders.Item, error) {
	for _, item := range items {
		if item.Id == id {
			return &item, nil
		}
	}

	return nil, fmt.Errorf("Cannot find Id '%s'", id)
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
