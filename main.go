package main

import (
	"log"
	"os"

	"github.com/bdibon/gator/internal/config"
)

type state struct {
	config *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v\n", err)
	}

	programState := state{&cfg}
	cmds := commands{handlers: map[string]func(s *state, cmd command) error{}}
	cmds.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("missing argument: command name")
	}
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	cmd := command{cmdName, cmdArgs}

	err = cmds.run(&programState, cmd)
	if err != nil {
		log.Fatalf("error running %s: %v", cmd.name, err)
	}
}
