package cmd

import (
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

		multiService := multiplexer.NewMultiplexerService(tmux)

		item := &dataproviders.Item{
			Id:      args[0],
			Display: args[0],
		}

		projectConfig := config.GetProjectConfig(item.Id)

		ok, err := multiService.SwitchSession(item)
		if err != nil {
			log.Fatal(err)
		}

		if ok {
			return
		}

		if err := multiService.CreateSession(item.Id, projectConfig); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	switchCmd.Flags().StringVarP(&configFile, "file", "f", "", "Path to config file (default is XDG_CONFIG/mux-session/config.toml)")
	switchCmd.Flags().StringVarP(&socket, "socket", "L", "", "tmux socket name for targeting a specific server")
	rootCmd.AddCommand(switchCmd)
}
