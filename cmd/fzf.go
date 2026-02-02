package cmd

import (
	"fmt"
	"log"

	"github.com/niedch/mux-session/internal/conf"
	"github.com/niedch/mux-session/internal/fzf"
	"github.com/spf13/cobra"
)

// fzfCmd represents the fzf command
var fzfCmd = &cobra.Command{
	Use:   "fzf",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := conf.Load(configFile)
		if err != nil {
			log.Fatal(err)
			return
		}

		provider := fzf.NewDirectoryProvider(config.SearchPaths)
		selected, err := fzf.StartApp(provider)
		if err != nil {
			log.Fatal(err)
		}
		if selected != "" {
			fmt.Println(selected)
		}
	},
}

func init() {
	rootCmd.AddCommand(fzfCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fzfCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fzfCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
