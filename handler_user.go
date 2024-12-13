package main

import (
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing argument: username")
	}
	username := cmd.args[0]
	err := s.config.SetUser(username)
	if err != nil {
		return fmt.Errorf("error setting user: %w", err)
	}
	fmt.Printf("username was set to %s\n", username)
	return nil
}
