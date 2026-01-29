package conf

import (
	"encoding/json"
	"fmt"
	"log"
)

func (c *Config) PrettyPrint() {
	jsonBytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Printf("Error marshaling config: %v\n", err)
		return
	}

	fmt.Println(string(jsonBytes))
}

func (c *ProjectConfig) PrettyPrint() {
	jsonBytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Printf("Error marshaling config: %v\n", err)
		return
	}

	fmt.Println(string(jsonBytes))
}
