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
	username := cmd.args[0]
	user, err := s.db.GetUser(context.Background(), sql.NullString{String: username, Valid: true})
	if err != nil {
		return fmt.Errorf("call to GetUser() failed for user: %s\n", username)
	}
	err = s.config.SetUser(user.Name.String)
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
		fmt.Printf("err: %v\n", err)
		return fmt.Errorf("unable to create user: %s", cmd.args[0])
	}
	err = handlerLogin(s, command{name: "login", args: cmd.args})
	if err != nil {
		return err
	}
	fmt.Printf("created user: %v\n", user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("reset does not take arguments")
	}
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return errors.New("could not TRUNCATE TABLE users")
	}
	fmt.Printf("TRUNCATED TABLE users done\n")
	return nil
}

func handlerListUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("users does not take arguments")
	}
	users, err := s.db.ListUsers(context.Background())
	if err != nil {
		return fmt.Errorf("users could not get users: %v\n", err)
	}
	for _, user := range users {
		fmt.Printf("* %s", user.Name.String)
		if user.Name.String == s.config.CurrentUserName {
			fmt.Printf(" (current)\n")
		} else {
			fmt.Print("\n")
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("agg does not take arguments")
	}
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("fetchFeed() failed")
	}
	fmt.Printf("Title: %s\n", feed.Channel.Title)
	fmt.Printf("Link: %s\n", feed.Channel.Link)
	fmt.Printf("Description: %s\n", feed.Channel.Description)
	fmt.Print("Items: \n")
	for _, item := range feed.Channel.Item {
		fmt.Printf("     Title: %s\n", item.Title)
		fmt.Printf("     Link: %s\n", item.Link)
		fmt.Printf("     Description: %s\n", item.Description)
		fmt.Printf("     PubDate: %s\n", item.PubDate)
	}

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
