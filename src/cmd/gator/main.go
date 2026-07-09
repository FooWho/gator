package main

import (
	"fmt"

	"github.com/FooWho/gator/src/internal/config"
)

func main() {

	config := config.Read()
	fmt.Printf("db_url: %s\n", config.DBUrl)
}
