package conf

import (
	"encoding/json"
	"fmt"

	"github.com/niedch/mux-session/internal/logger"
)

func (c *Config) PrettyPrint() {
	jsonBytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		logger.Printf("Error marshaling config: %v\n", err)
		return
	}

	fmt.Println(string(jsonBytes))
}

func (c *ProjectConfig) PrettyPrint() {
	jsonBytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		logger.Printf("Error marshaling config: %v\n", err)
		return
	}

	fmt.Println(string(jsonBytes))
}
