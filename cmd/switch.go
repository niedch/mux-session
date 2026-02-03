package cmd

import (
	"log"

	"github.com/niedch/mux-session/internal/conf"
	"github.com/niedch/mux-session/internal/dataproviders"
	"github.com/niedch/mux-session/internal/fzf"
	"github.com/niedch/mux-session/internal/multiplexer"
	"github.com/niedch/mux-session/internal/tmux"
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch to an existing or create a new tmux session",
	Long: `Opens an interactive fzf picker to select a directory or existing tmux session.
If the selected session exists, switches to it. If not, creates a new session
with the configured windows.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := conf.Load(configFile)
		if err != nil {
			log.Fatal(err)
			return
		}

		tmux, err := tmux.NewTmux(socket)
		if err != nil {
			log.Fatal(err)
			return
		}

		multiService := multiplexer.NewMultiplexerService(tmux)

		provider := dataproviders.NewDirectoryProvider(config.SearchPaths)
		tmuxProvider := dataproviders.NewTmuxProvider(tmux)

		selected, err := fzf.StartApp(dataproviders.NewComposeProvider(provider, tmuxProvider))
		if err != nil {
			log.Fatal(err)
		}

		if selected == nil {
			log.Println("No selection done")
			return
		}

		projectConfig := config.GetProjectConfig(selected.Id)
		ok, err := multiService.SwitchSession(selected)
		if err != nil {
			log.Fatal(err)
		}

		if ok {
			return
		}

		if err := multiService.CreateSession(selected.Display, projectConfig); err != nil {
			log.Fatal(err)
			return
		}
	},
}

func init() {
	switchCmd.Flags().StringVarP(&configFile, "file", "f", "", "Path to config file (default is XDG_CONFIG/mux-session/config.toml)")
	switchCmd.Flags().StringVarP(&socket, "socket", "L", "", "tmux socket name for targeting a specific server")
	rootCmd.AddCommand(switchCmd)
}
