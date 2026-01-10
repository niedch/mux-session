package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/koki-develop/go-fzf"
	"github.com/niedch/mux-session/internal/conf"
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
		// Load config
		config, err := conf.Load()
		if err != nil {
			log.Fatal(err)
			return
		}

		var dirs []string
		for _, searchPath := range config.Search_paths {
			entries, err := os.ReadDir(searchPath)
			if err != nil {
				log.Printf("Error reading %s: %v", searchPath, err)
				continue
			}

			for _, entry := range entries {
				if entry.IsDir() {
					// Skip hidden directories
					if !strings.HasPrefix(entry.Name(), ".") {
						fullPath := filepath.Join(searchPath, entry.Name())
						dirs = append(dirs, fullPath)
					}
				}
			}
		}

		// Create fzf instance
		f, err := fzf.New(
			fzf.WithInputPosition(fzf.InputPositionBottom),
		)
		if err != nil {
			log.Fatal(err)
			return
		}

		// Run fzf to select directory
		idxs, err := f.Find(dirs, func(i int) string {
			return dirs[i]
		}, fzf.WithPreviewWindow(func(i, width, height int) string {
			if i == -1 {
				return ""
			}

			// Get directory path
			dirPath := dirs[i]

			// Read directory contents
			entries, err := os.ReadDir(dirPath)
			if err != nil {
				return fmt.Sprintf("Error reading %s: %v", dirPath, err)
			}

			// Build preview content
			var preview strings.Builder
			preview.WriteString(fmt.Sprintf("Directory: %s\n\n", dirPath))

			for _, entry := range entries {
				if entry.IsDir() {
					preview.WriteString(fmt.Sprintf("📁 %s/\n", entry.Name()))
				} else {
					preview.WriteString(fmt.Sprintf("📄 %s\n", entry.Name()))
				}
			}

			return preview.String()
		}))
		if err != nil {
			log.Fatal(err)
			return
		}

		// Print selected directory
		for _, i := range idxs {
			selectedDir := dirs[i]
			log.Printf("Selected directory: %s", selectedDir)
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
