package main

import (
	"fmt"
	"log"

	"github.com/santokan/gator/internal/config"
)

func main() {
	// Read the config file.
	// Set the current user to "lane" (actually, you should use your name instead) and update the config file on disk.
	// Read the config file again and print the contents of the config struct to the terminal.

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	fmt.Printf("Read config: %+v\n", cfg)

	err = cfg.SetUser("santokan")

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	fmt.Printf("Read config: %+v\n", cfg)
}
