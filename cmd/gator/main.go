package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/FooWho/gator/internal/config"
	"github.com/FooWho/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg := &config.Config{}
	if err := cfg.Read(); err != nil {
		fmt.Printf("error reading config: %s\n", err)
		os.Exit(1)
	}
	gatorState := state{config: cfg}
	gatorCommands := commands{}
	gatorCommands.register("login", handlerLogin)
	gatorCommands.register("register", handlerRegister)
	gatorCommands.register("reset", handlerReset)
	gatorCommands.register("users", handlerListUsers)
	gatorCommands.register("agg", handlerAgg)
	gatorCommands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	gatorCommands.register("feeds", handlerListFeeds)
	gatorCommands.register("follow", middlewareLoggedIn(handlerFollowFeed))
	gatorCommands.register("following", middlewareLoggedIn(handlerFollowing))
	gatorCommands.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		fmt.Print("could not open database")
		os.Exit(1)
	}
	dbQueries := database.New(db)
	gatorState.db = dbQueries

	args := os.Args
	if len(args) < 2 {
		fmt.Print("gator requires a command as an argument\n")
		os.Exit(1)
	}

	cmd := command{name: args[1], args: args[2:]}
	fmt.Printf("Got command:\n")
	fmt.Printf("   name: %s\n", cmd.name)
	if len(cmd.args) > 0 {
		fmt.Printf("   args[0]: %v\n", cmd.args[0])
	}
	if err := gatorCommands.run(&gatorState, cmd); err != nil {
		fmt.Printf("error running command: %s\n", err)
		os.Exit(1)
	}
}
