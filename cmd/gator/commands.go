package main

import (
	"errors"
	"fmt"

	"github.com/FooWho/gator/internal/config"
	"github.com/FooWho/gator/internal/database"
)

type command struct {
	name string
	args []string
}

type commands struct {
	cmdHandler map[string]func(*state, command) error
}

type state struct {
	config *config.Config
	db     *database.Queries
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("login requires username as argument")
	}
	err := s.config.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("unable to set user: %w", err)
	}
	fmt.Printf("Username set - %s\n", cmd.args[0])
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	if s.config == nil {
		return errors.New("state not initialized")
	}
	f, ok := c.cmdHandler[cmd.name]
	if !ok {
		return fmt.Errorf("no handler for command: %s", cmd.name)
	}
	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	if c.cmdHandler == nil {
		c.cmdHandler = make(map[string]func(*state, command) error)
	}
	c.cmdHandler[name] = f
}
