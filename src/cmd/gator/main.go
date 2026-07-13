package main

import (
	"fmt"

	"github.com/FooWho/gator/src/internal/config"
)

func main() {

	configuration := &config.Config{}
	configuration.Read()
	configuration.SetUser("jason")
	configuration.Read()
	fmt.Printf("%+v\n", *configuration)
}
