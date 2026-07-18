package main

import (
	"fmt"
	"os"

	"github.com/FooWho/gator/src/internal/config"
)

func main() {

	gatorState := state{}
	gatorCommands := commands{}
	gatorState.config = &config.Config{}
	gatorState.config.Read()
	gatorCommands.cmdHandler = make(map[string]func(*state, command) error)
	gatorCommands.cmdHandler["login"] = handlerLogin

	args := os.Args
	if len(args) < 2 {
		fmt.Print("gator requires a command as an argument\n")
		os.Exit(1)
	} else {
		cmd := command{name: args[1], args: args[2:]}
		fmt.Printf("Got command:\n")
		fmt.Printf("   name: %s\n", cmd.name)
		fmt.Printf("   args[0]: %v\n", cmd.args[0])
		gatorCommands.cmdHandler[cmd.name](&gatorState, cmd)
	}

}
