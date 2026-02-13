package cmd

import (
	"fmt"

	"github.com/niedch/mux-session/internal/conf"
	"github.com/niedch/mux-session/internal/dataproviders"
	"github.com/niedch/mux-session/internal/logger"
	"github.com/niedch/mux-session/internal/tmux"
	"github.com/niedch/mux-session/internal/tree"
	"github.com/spf13/cobra"
)

var listSessionsCmd = &cobra.Command{
	Use:   "list-sessions",
	Short: "List all directories and active tmux sessions",
	Long:  `Displays all configured directories and currently active tmux sessions, similar to what appears in the interactive picker.`,
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
		tmuxWrapper, err := tmux.NewTmux(socket)
		if err != nil {
			logger.Fatalf("Failed to initialize tmux: %v\n", err)
		}

		directoryProvider := dataproviders.NewDirectoryProvider(config.SearchPaths)
		tmuxProvider := dataproviders.NewTmuxProvider(tmuxWrapper)
		composedProvider := dataproviders.NewDeduplicatorProvider(directoryProvider, tmuxProvider).WithMarkDuplicates(true)

		logger.Printf("Getting items from providers\n")
		items, err := composedProvider.GetItems()
		if err != nil {
			logger.Fatalf("Failed to get items: %v\n", err)
		}
		logger.Printf("Found %d items\n", len(items))

		if len(items) == 0 {
			fmt.Println("No items found")
			return
		}

		flattenedItems := tree.FlattenItems(items)
		logger.Printf("Displaying %d items\n", len(flattenedItems))
		for _, item := range flattenedItems {
			fmt.Println(item.Display)
		}
	},
}

func init() {
	rootCmd.AddCommand(listSessionsCmd)
}
