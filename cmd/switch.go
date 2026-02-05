package cmd

import (
	"fmt"
	"log"

	"github.com/niedch/mux-session/internal/conf"
	"github.com/niedch/mux-session/internal/dataproviders"
	"github.com/niedch/mux-session/internal/multiplexer"
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
		config, err := conf.Load(configFile)
		if err != nil {
			log.Fatal(err)
		}

		tmux, err := tmux.NewTmux(socket)
		if err != nil {
			log.Fatal(err)
		}

		directoryProvider := dataproviders.NewDirectoryProvider(config.SearchPaths)
		tmuxProvider := dataproviders.NewTmuxProvider(tmux)
		composedProvider := dataproviders.NewComposeProvider(directoryProvider, tmuxProvider).WithMarkDuplicates(true)

		items, err := composedProvider.GetItems()
		if err != nil {
			log.Fatal(err)
		}

		item, err := findItem(args[0], items)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Found item with id '%s', and display '%s'\n", item.Id, item.Display)
		
		multiService := multiplexer.NewMultiplexerService(tmux)
		projectConfig := config.GetProjectConfig(item.Id)

		ok, err := multiService.SwitchSession(item)
		if err != nil {
			log.Fatal(err)
		}

		if ok {
			return
		}

		if err := multiService.CreateSession(item, projectConfig); err != nil {
			log.Fatal(err)
		}
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
	switchCmd.Flags().StringVarP(&configFile, "file", "f", "", "Path to config file (default is XDG_CONFIG/mux-session/config.toml)")
	switchCmd.Flags().StringVarP(&socket, "socket", "L", "", "tmux socket name for targeting a specific server")
	rootCmd.AddCommand(switchCmd)
}
