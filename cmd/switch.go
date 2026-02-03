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
		log.Printf("[DEBUG] switch: Starting with ID: %s", args[0])

		config, err := conf.Load(configFile)
		if err != nil {
			log.Printf("[DEBUG] switch: Failed to load config: %v", err)
			log.Fatal(err)
			return
		}
		log.Printf("[DEBUG] switch: Config loaded successfully")

		tmux, err := tmux.NewTmux(socket)
		if err != nil {
			log.Printf("[DEBUG] switch: Failed to create tmux wrapper: %v", err)
			log.Fatal(err)
			return
		}
		log.Printf("[DEBUG] switch: Tmux wrapper created")

		multiService := multiplexer.NewMultiplexerService(tmux)

		item := &dataproviders.Item{
			Id:      args[0],
			Display: args[0],
		}
		log.Printf("[DEBUG] switch: Created item with ID: %s", item.Id)

		projectConfig := config.GetProjectConfig(item.Id)
		log.Printf("[DEBUG] switch: Got project config, windows: %d", len(projectConfig.WindowConfig))

		ok, err := multiService.SwitchSession(item)
		if err != nil {
			log.Printf("[DEBUG] switch: SwitchSession failed: %v", err)
			log.Fatal(err)
		}
		log.Printf("[DEBUG] switch: SwitchSession returned ok=%v", ok)

		if ok {
			log.Printf("[DEBUG] switch: Successfully switched to existing session")
			return
		}

		log.Printf("[DEBUG] switch: Session not found, attempting to create...")
		if err := multiService.CreateSession(item.Id, projectConfig); err != nil {
			log.Printf("[DEBUG] switch: CreateSession failed: %v", err)
			log.Fatal(err)
			return
		}
		log.Printf("[DEBUG] switch: Session created successfully")
	},
}

func init() {
	switchCmd.Flags().StringVarP(&configFile, "file", "f", "", "Path to config file (default is XDG_CONFIG/mux-session/config.toml)")
	switchCmd.Flags().StringVarP(&socket, "socket", "L", "", "tmux socket name for targeting a specific server")
	rootCmd.AddCommand(switchCmd)
}
