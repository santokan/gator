package main

import (
	"log"
	"os"

	"github.com/santokan/gator/internal/config"
)

type state struct {
	config *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
	}
	s := &state{config: &cfg}
	cmds := &commands{registeredCommands: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args..]")
		return
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(s, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
