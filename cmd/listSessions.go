package cmd

import (
	"fmt"
	"log"

	"github.com/niedch/mux-session/internal/conf"
	"github.com/niedch/mux-session/internal/dataproviders"
	"github.com/niedch/mux-session/internal/tmux"
	"github.com/niedch/mux-session/internal/tree"
	"github.com/spf13/cobra"
)

var listSessionsCmd = &cobra.Command{
	Use:   "list-sessions",
	Short: "List all directories and active tmux sessions",
	Long:  `Displays all configured directories and currently active tmux sessions, similar to what appears in the interactive picker.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := conf.Load(configFile)
		if err != nil {
			log.Fatal(err)
			return
		}

		tmuxWrapper, err := tmux.NewTmux(socket)
		if err != nil {
			log.Fatal(err)
			return
		}

		directoryProvider := dataproviders.NewDirectoryProvider(config.SearchPaths)
		tmuxProvider := dataproviders.NewTmuxProvider(tmuxWrapper)
		composedProvider := dataproviders.NewDeduplicatorProvider(directoryProvider, tmuxProvider).WithMarkDuplicates(true)

		items, err := composedProvider.GetItems()
		if err != nil {
			log.Fatal(err)
			return
		}

		if len(items) == 0 {
			fmt.Println("No items found")
			return
		}

		flattenedItems := tree.FlattenItems(items)
		for _, item := range flattenedItems {
			fmt.Println(item.Display)
		}
	},
}

func init() {
	listSessionsCmd.Flags().StringVarP(&configFile, "file", "f", "", "Path to config file (default is XDG_CONFIG/mux-session/config.toml)")
	listSessionsCmd.Flags().StringVarP(&socket, "socket", "L", "", "tmux socket name for targeting a specific server")
	rootCmd.AddCommand(listSessionsCmd)
}
