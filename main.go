package main

import (
	"fmt"
	"os"

	"github.com/Sanghun1Adam1Park/blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("error reading config:", err)
		os.Exit(1)
	}

	s := &state{
		cfg: &cfg,
	}

	cmds := commands{
		commandMap: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		fmt.Print("error with number of argument")
		os.Exit(1)
	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}

	if err := cmds.run(s, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
