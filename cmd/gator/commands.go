package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/FooWho/gator/internal/config"
	"github.com/FooWho/gator/internal/database"
	"github.com/google/uuid"
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

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("register requires username as argument")
	}
	user, err := s.db.CreateUser(context.Background(),
		database.CreateUserParams{ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      sql.NullString{String: cmd.args[0], Valid: true},
		})
	if err != nil {
		return fmt.Errorf("unable to create user: %s", cmd.args[0])
	}
	err = handlerLogin(s, command{name: "login", args: cmd.args})
	if err != nil {
		return err
	}
	fmt.Printf("created user: %v\n", user)
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
