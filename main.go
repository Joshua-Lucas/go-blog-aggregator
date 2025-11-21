package main

import (
	"fmt"
	"os"
	"github.com/Joshua-Lucas/blog-aggregator/internal/commands"
	"github.com/Joshua-Lucas/blog-aggregator/internal/config"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	state := commands.State{Config: &cfg}	

	cmd := commands.Commands{
		Handlers: make(map[string]func(*commands.State, commands.Command) error),
	}

	cmd.Register("login", commands.HandlerLogin)

	fmt.Println(cmd)

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "error: There are less than two arguments.\n")
		os.Exit(1)
	}

	userArg := commands.Command {
		Name: os.Args[1],
		Args: []string{},
	}

	if (len(os.Args) > 2) {
		for _, v := range os.Args[2:] { 
    	userArg.Args = append(userArg.Args, v)
		}	
	}

	err = cmd.Run(&state, userArg)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

}
