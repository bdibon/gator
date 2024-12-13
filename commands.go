package main

import "fmt"

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(s *state, cmd command) error
}

func (c *commands) register(name string, f func(s *state, cmd command) error) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("no handler for command %s", cmd.name)
	}
	err := handler(s, cmd)
	if err != nil {
		return err
	}
	return nil
}
