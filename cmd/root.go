package cmd

import (
	"log"
	"os"

	"github.com/GianlucaP106/gotmux/gotmux"
	"github.com/niedch/mux-session/internal/conf"
	"github.com/niedch/mux-session/internal/fzf"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mux-session",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := conf.Load()
		if err != nil {
			log.Fatal(err)
			return
		}

		provider := fzf.NewDirectoryProvider(config.Search_paths)
		selectedIndices, err := fzf.StartFzf(provider)
		if err != nil {
			log.Fatal(err)
			return
		}

		tmux, _ := gotmux.DefaultTmux()
		sessions, _ := tmux.ListSessions()
		
		for _, session := range sessions {
			log.Println(session.Name)
		}

		items, err := provider.GetItems()
		if err != nil {
			log.Fatal(err)
			return
		}

		// Print selected directories
		for _, i := range selectedIndices {
			if i < len(items) {
				log.Printf("Selected directory: %s", items[i])
			}
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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mux-session.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
