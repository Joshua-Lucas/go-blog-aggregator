package commands

import (
	"fmt"

	"github.com/Joshua-Lucas/blog-aggregator/internal/config"
)

type State struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*State, command) error
}

func (c *commands) run(s *State, cmd command) error {
	v, ok := c.handlers[cmd.name]

	if !ok {
		return fmt.Errorf("There is no handler for the provided command")
	}

	err := v(s, cmd)

	if err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(*State, command) error) {
	c.handlers[name] = f
}

func handlerLogin(s *State, cmd command) error {
	if len(cmd.args) <= 0 {
		return fmt.Errorf("Login in handler expects a single argument, the username")
	}

	userName := cmd.args[0]

	err := s.cfg.SetUser(userName)

	if err != nil {
		return err
	}

	println("User has been set")

	return nil
}
